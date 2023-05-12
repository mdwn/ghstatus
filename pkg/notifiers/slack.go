package notifiers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mdwn/ghstatus/pkg/ghstatus"
	"github.com/mdwn/ghstatus/pkg/notifier"
	"github.com/slack-go/slack"
	"github.com/spf13/pflag"
)

const (
	Slack = "slack"

	slackGoodEmoji = ":white_check_mark:"
	slackBadEmoji  = ":warning:"
	slackInfoEmoji = ":information_source:"
)

var (
	slackOauthToken  string
	slackChannel     string
	slackJoinChannel bool
)

func init() {
	registeredNotifiers[Slack] = NewSlackNotifier

	flags := pflag.NewFlagSet("slack", pflag.ContinueOnError)
	flags.StringVar(&slackOauthToken, "slack-oauth-token", "", "The Slack oauth token to use.")
	flags.StringVar(&slackChannel, "slack-channel", "", "The Slack channel to notify.")
	flags.BoolVar(&slackJoinChannel, "slack-join-channel", false, "Whether the bot should attempt to join the channel.")
	flags.FlagUsages()
	notifierFlags.AddFlagSet(flags)
}

// SlackNotifier writes the output to Slack.
type SlackNotifier struct {
	client    *slack.Client
	channelID string
}

// NewSlackNotifier will return a Slack notifier.
func NewSlackNotifier() (notifier.Notifier, error) {
	if slackOauthToken == "" {
		return nil, errors.New("OAuth token must be supplied for the Slack notifier")
	}

	if slackChannel == "" {
		return nil, errors.New("channel must be supplied for the Slack notifier")
	}

	client := slack.New(slackOauthToken)

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

	if slackJoinChannel {
		_, _, _, err := client.JoinConversation(channelID)
		if err != nil {
			return nil, fmt.Errorf("error joining channel: %w", err)
		}
	}

	return &SlackNotifier{
		client:    slack.New(slackOauthToken),
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

	resp, resp2, err := s.client.PostMessageContext(ctx, s.channelID, slack.MsgOptionBlocks(blocks.BlockSet...))
	if err != nil {
		return fmt.Errorf("error posting message: %s, %s, %w", resp, resp2, err)
	}

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
		slackMsgText = fmt.Sprintf("%s Github is reporting a **%s** outage", slackBadEmoji, status.Indicator)
	}

	text := slack.NewSectionBlock(slack.NewTextBlockObject(
		slack.MarkdownType, slackMsgText, false, false,
	), nil, nil, slack.SectionBlockOptionBlockID("status"))

	blocks.BlockSet = append(blocks.BlockSet,
		slack.NewHeaderBlock(slack.NewTextBlockObject(slack.PlainTextType, "Status", false, false)),
		text)
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
}
