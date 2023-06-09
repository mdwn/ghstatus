package notifiers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/mdwn/ghstatus/pkg/ghstatus"
	"github.com/mdwn/ghstatus/pkg/logging"
	"github.com/mdwn/ghstatus/pkg/notifier"
	"github.com/ory/viper"
	"github.com/slack-go/slack"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

const (
	Slack = "slack"

	slackGoodEmoji = ":white_check_mark:"
	slackBadEmoji  = ":warning:"
	slackInfoEmoji = ":information_source:"

	slackOAuthTokenCfg  = "slack.oauth.token"
	slackOAuthTokenFlag = "slack-oauth-token"
	slackOAuthTokenEnv  = "SLACK_OAUTH_TOKEN"

	slackChannelCfg  = "slack.channel"
	slackChannelFlag = "slack-channel"
	slackChannelEnv  = "SLACK_CHANNEL"

	slackJoinChannelCfg  = "slack.join.channel"
	slackJoinChannelFlag = "slack-join-channel"
	slackJoinChannelEnv  = "SLACK_JOIN_CHANNEL"
)

func init() {
	if err := RegisterNotifier(Slack, NewSlackNotifier); err != nil {
		panic(err.Error())
	}

	flags := pflag.NewFlagSet("slack", pflag.ContinueOnError)
	flags.String(slackOAuthTokenFlag, "", "The Slack oauth token to use.")
	flags.String(slackChannelFlag, "", "The Slack channel to notify.")
	flags.Bool(slackJoinChannelFlag, false, "Whether the bot should attempt to join the channel.")

	notifierFlags.AddFlagSet(flags)

	err := multierror.Append(nil,
		viper.BindPFlag(slackOAuthTokenCfg, flags.Lookup(slackOAuthTokenFlag)),
		viper.BindEnv(slackOAuthTokenCfg, slackOAuthTokenEnv),

		viper.BindPFlag(slackChannelCfg, flags.Lookup(slackChannelFlag)),
		viper.BindEnv(slackChannelCfg, slackChannelEnv),

		viper.BindPFlag(slackJoinChannelCfg, flags.Lookup(slackJoinChannelFlag)),
		viper.BindEnv(slackJoinChannelCfg, slackJoinChannelEnv),
	)

	if err.ErrorOrNil() != nil {
		panic(fmt.Sprintf("error binding Slack configs: %v", err))
	}

}

// SlackNotifier writes the output to Slack.
type SlackNotifier struct {
	log       *zap.Logger
	client    *slack.Client
	channelID string
}

// NewSlackNotifier will return a Slack notifier.
func NewSlackNotifier(log *zap.Logger) (notifier.Notifier, error) {
	slackOAuthToken := viper.GetString(slackOAuthTokenCfg)
	slackChannel := viper.GetString(slackChannelCfg)

	if slackOAuthToken == "" {
		return nil, errors.New("OAuth token must be supplied for the Slack notifier")
	}

	if slackChannel == "" {
		return nil, errors.New("channel must be supplied for the Slack notifier")
	}

	client := slack.New(slackOAuthToken)

	var channelID string
	if strings.HasPrefix(slackChannel, "#") {
		channelName := strings.TrimPrefix(slackChannel, "#")

		conversations, next, err := client.GetConversations(&slack.GetConversationsParameters{})
		for {
			if err != nil {
				return nil, fmt.Errorf("error finding channel: %w", err)
			}

			for _, conversation := range conversations {
				if conversation.Name == channelName {
					channelID = conversation.ID
					break
				}
			}

			if channelID != "" || next == "" {
				break
			}

			conversations, next, err = client.GetConversations(&slack.GetConversationsParameters{Cursor: next})
		}
	} else {
		channelID = slackChannel
	}

	if channelID == "" {
		return nil, fmt.Errorf("unable to find channel %s", slackChannel)
	}

	if viper.GetBool(slackJoinChannelCfg) {
		_, _, _, err := client.JoinConversation(channelID)
		if err != nil {
			return nil, fmt.Errorf("error joining channel: %w", err)
		}
	}

	return &SlackNotifier{
		log:       logging.WithComponent(log, "slack"),
		client:    slack.New(slackOAuthToken),
		channelID: channelID,
	}, nil
}

// Name is the name of the notifier.
func (*SlackNotifier) Name() string {
	return Slack
}

// Notify will notify an underlying system with the given message.
func (s *SlackNotifier) Notify(ctx context.Context, msg notifier.Message) error {
	blocks := &slack.Blocks{}

	s.changedStatus(msg, blocks)
	s.changedComponents(msg, blocks)
	s.changedIncidents(msg, blocks)
	s.changedScheduledMaintenances(msg, blocks)

	if len(blocks.BlockSet) == 0 {
		s.log.Debug("Slack notifier found no changes.")
		return nil
	}

	_, _, err := s.client.PostMessageContext(ctx, s.channelID, slack.MsgOptionBlocks(blocks.BlockSet...))
	if err != nil {
		return fmt.Errorf("error posting message: %w", err)
	}

	s.log.Debug("Slack notified of changes.")

	return nil
}

// Cleanup performs any cleanup steps.
func (s *SlackNotifier) Cleanup() error {
	return nil
}

// changedStatus updates the message to contain any information about the changed status.
func (s *SlackNotifier) changedStatus(msg notifier.Message, blocks *slack.Blocks) {
	status := msg.ChangedStatus
	if status == nil {
		return
	}

	var slackMsgText string
	switch status.Indicator {
	case ghstatus.None:
		slackMsgText = fmt.Sprintf("%s Github reports no outages", slackGoodEmoji)
	default:
		slackMsgText = fmt.Sprintf("%s Github is reporting a *%s* outage", slackBadEmoji, status.Indicator)
	}

	text := slack.NewSectionBlock(slack.NewTextBlockObject(
		slack.MarkdownType, slackMsgText, false, false,
	), nil, nil, slack.SectionBlockOptionBlockID("status"))

	blocks.BlockSet = append(blocks.BlockSet,
		slack.NewHeaderBlock(slack.NewTextBlockObject(slack.PlainTextType, "Status", false, false)),
		text)

	s.log.Debug("Status change being sent to Slack")
}

// changedComponents updates the message to contain any information about the changed components.
func (s *SlackNotifier) changedComponents(msg notifier.Message, blocks *slack.Blocks) {
	if len(msg.ChangedComponents) == 0 {
		return
	}

	blocks.BlockSet = append(blocks.BlockSet,
		slack.NewHeaderBlock(slack.NewTextBlockObject(slack.PlainTextType, "Components", false, false)))

	for _, component := range msg.ChangedComponents {
		var slackMsgText string

		switch component.Status {
		case ghstatus.Operational:
			slackMsgText = fmt.Sprintf("%s %s is operational", slackGoodEmoji, component.Name)
		default:
			slackMsgText = fmt.Sprintf("%s %s is reporting %s", slackBadEmoji, component.Name, component.Status)
		}

		text := slack.NewSectionBlock(slack.NewTextBlockObject(
			slack.MarkdownType, slackMsgText, false, false,
		), nil, nil, slack.SectionBlockOptionBlockID(fmt.Sprintf("component-%s", component.Name)))

		blocks.BlockSet = append(blocks.BlockSet, text)
	}

	s.log.Debug("Components change being sent to Slack")
}

// changedIncidents updates the message to contain any information about the changed incidents.
func (s *SlackNotifier) changedIncidents(msg notifier.Message, blocks *slack.Blocks) {
	if len(msg.ChangedIncidents) == 0 {
		return
	}

	blocks.BlockSet = append(blocks.BlockSet,
		slack.NewHeaderBlock(slack.NewTextBlockObject(slack.PlainTextType, "Incidents", false, false)))

	for _, incident := range msg.ChangedIncidents {
		var slackMsgText string

		switch incident.Status {
		case ghstatus.Investigating:
			slackMsgText = fmt.Sprintf("%s %q is being investigated", slackBadEmoji, incident.Name)
		case ghstatus.Identified:
			slackMsgText = fmt.Sprintf("%s The cause of %q has been identified", slackInfoEmoji, incident.Name)
		case ghstatus.Monitoring:
			slackMsgText = fmt.Sprintf("%s %q is being monitored", slackInfoEmoji, incident.Name)
		case ghstatus.Resolved:
			slackMsgText = fmt.Sprintf("%s %q has been resolved", slackGoodEmoji, incident.Name)
		case ghstatus.Postmorten:
			slackMsgText = fmt.Sprintf("%s %q has a postmortem", slackGoodEmoji, incident.Name)
		default:
			slackMsgText = fmt.Sprintf("%s %q has status %s", slackInfoEmoji, incident.Name, incident.Status)
		}

		slackMsgText += fmt.Sprintf(" (impact %s)", incident.Impact)

		if len(incident.IncidentUpdates) > 0 {
			slackMsgText += fmt.Sprintf(": %s", incident.IncidentUpdates[0].Body)
		}

		text := slack.NewSectionBlock(slack.NewTextBlockObject(
			slack.MarkdownType, slackMsgText, false, false,
		), nil, nil, slack.SectionBlockOptionBlockID(fmt.Sprintf("incident-%s", incident.ID)))

		blocks.BlockSet = append(blocks.BlockSet, text)
	}

	s.log.Debug("Incidents change being sent to Slack")
}

// changedScheduledMaintenances updates the message to contain any information about the changed scheduled maintenances.
func (s *SlackNotifier) changedScheduledMaintenances(msg notifier.Message, blocks *slack.Blocks) {
	if len(msg.ChangedScheduledMaintenances) == 0 {
		return
	}

	blocks.BlockSet = append(blocks.BlockSet,
		slack.NewHeaderBlock(slack.NewTextBlockObject(slack.PlainTextType, "Scheduled Maintenances", false, false)))

	for _, scheduledMaintenance := range msg.ChangedScheduledMaintenances {
		var slackMsgText string

		switch scheduledMaintenance.Status {
		case ghstatus.Scheduled:
			slackMsgText = fmt.Sprintf("%s %q is scheduled", slackInfoEmoji, scheduledMaintenance.Name)
		case ghstatus.InProgress:
			slackMsgText = fmt.Sprintf("%s %q is in progress", slackInfoEmoji, scheduledMaintenance.Name)
		case ghstatus.Verifying:
			slackMsgText = fmt.Sprintf("%s %q is being verified", slackInfoEmoji, scheduledMaintenance.Name)
		case ghstatus.Completed:
			slackMsgText = fmt.Sprintf("%s %q is completed", slackInfoEmoji, scheduledMaintenance.Name)
		default:
			slackMsgText = fmt.Sprintf("%s %q has status %s", slackInfoEmoji, scheduledMaintenance.Name, scheduledMaintenance.Status)
		}

		slackMsgText += fmt.Sprintf(" (expected impact %s)", scheduledMaintenance.Impact)

		text := slack.NewSectionBlock(slack.NewTextBlockObject(
			slack.PlainTextType, slackMsgText, false, false,
		), nil, nil, slack.SectionBlockOptionBlockID(fmt.Sprintf("scheduled-maintenance-%s", scheduledMaintenance.ID)))

		blocks.BlockSet = append(blocks.BlockSet, text)
	}

	s.log.Debug("Scheduled maintenances change being sent to Slack")
}
