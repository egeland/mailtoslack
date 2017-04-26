# Mail to Slack

This is a simple SMTP to Slack server.

It listens for mail on a `PORT` (defaults to `2525`), and forwards it into
a `SLACK_CHANNEL` as specified by the `SLACK_TOKEN`'s permissions'.

## Configuration

Everything is done with environment variables:

* `PORT` - the port the server should listen on. Defaults to `2525`.
* `SLACK_TOKEN` - A token in the format
`xoxp-xxxxxxxxxxx-yyyyyyyyyyyy-zzzzzzzzzzzz-ssssssssssssssssssssssssssssssss`. Required.
* `SLACK_CHANNEL` - The channel to send messages to, e.g. `#mail`. Required.
* `DOMAIN_LIST` - A comma-separated list of domains to accept messages for. Optional. If unset, all messages are accepted.

## License
MIT - https://opensource.org/licenses/MIT
