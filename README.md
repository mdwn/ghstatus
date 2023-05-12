# ghstatus

[![](https://pkg.go.dev/badge/github.com/mdwn/ghstatus?utm_source=godoc)](https://pkg.go.dev/github.com/mdwn/ghstatus)

ghstatus is a lightweight client and monitor for the Github Status API as documented [here](https://www.githubstatus.com/api/).
The intention of this library is to provide programmatic golang access to the API and provide monitor
methods for notifications when Github is down.

## The client

The client can be found in [`pkg/ghstatus`](https://github.com/mdwn/ghstatus/tree/main/pkg/ghstatus). Using the client is pretty
straightforward, as the Github Status API is read only, and no auth tokens or API keys are currently required to use it:

```go
ctx := context.Background()
client := ghstatus.NewClient()

summaryResponse, err := client.Summary(ctx)
if err != nil {
  return fmt.Errorf("error getting client summary: %w", err)
}

...

allIncidents, err := client.AllIncidents(ctx)

...

```

Retry/backoff is supplied by using using Hashicorp's [retryablehttp module](https://github.com/hashicorp/go-retryablehttp).
However, there are no documented rate limits or recommended backoff timings, so this may be overkill.

## CLI

The CLI provides methods for querying the current Github Status and to output it in various formats.

### Table output

Table output provides a human readable way of reading the various Github Status endpoints.

```
$ ghstatus summary -f table
# Status

+-----------+-------------------------+
| INDICATOR |       DESCRIPTION       |
+-----------+-------------------------+
| none      | All Systems Operational |
+-----------+-------------------------+


# Components
...
```

### YAML output

Outputs the status in a YAML format which mirrors the native JSON format.

```
$ ghstatus status -f yaml
page:
    id: kctbh9vrtdwd
    name: GitHub
    url: https://www.githubstatus.com
    updated_at: 2023-05-15T07:51:14.591Z
status:
    description: All Systems Operational
    indicator: none
```

### JSON output

Outputs the status in the native JSON format.

```
$ ghstatus status -f json | jq .
{
  "page": {
    "id": "kctbh9vrtdwd",
    "name": "GitHub",
    "url": "https://www.githubstatus.com",
    "updated_at": "2023-05-15T07:51:14.591Z"
  },
  "status": {
    "description": "All Systems Operational",
    "indicator": "none"
  }
}
```

### All API query commands

All of the currently available API endpoints are queryable from the CLI utility.

```
$ ghstatus summary
$ ghstatus status
$ ghstatus components
$ ghstatus incidents unresolved
$ ghstatus incidents all
$ ghstatus scheduled-maintenances upcoming
$ ghstatus scheduled-maintenances active
$ ghstatus scheduled-maintenances all
```

## Monitor

The CLI additionally supports the `monitor` command which can use various notifiers. The current notifiers are:

### stdout

This notifier writes changes to stdout.

### file

This notifier writes changes to a configured file. Requires the following flags:

- `--fn-file-path` to the expected output file.

### slack

This notifier writes changes to a Slack channel. Requires the following flags:

- `--slack-oauth-token` to set the Slack oauth token.
- `--slack-channel` the Slack channel to post updates to. Can be either of the form `#channel-name` or the actual channel ID.
- `--slack-join-channel` whether the bot should attempt to join the channel.

The oauth token requires the following Slack scopes:

- `channels:join` to join the target channel. This is only needed if attempting to use `--slack-join-channel`. If you would rather not
  use this, you can elect to invite the bot explicitly.
- `channels:read` to find the target channel by its friendly name rather than the channel ID. If using a channel ID, this is not needed.
- `chat:write` to write status messages to the channel.