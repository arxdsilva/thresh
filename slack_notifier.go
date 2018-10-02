package thresh

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nlopes/slack"
)

// SlackNotifier returns a NotifierFunc for slack using the provided token.
// An invalid token will result in an error during CheckStatus.
// Channel may be either a channel name or channel ID, like token if this
// is invalid it will result in an error during CheckStatus.
func SlackNotifier(token, channel string, opts ...slack.MsgOption) NotifierFunc {
	return func(addr string) error {
		client := slack.New(token)
		channels, err := client.GetChannels(true)
		if err != nil {
			return err
		}

		var c *slack.Channel
		for _, _c := range channels {
			if strings.EqualFold(_c.Name, channel) ||
				strings.Compare(_c.ID, channel) == 0 {
				c = &_c
				break
			}
		}

		if c == nil {
			return errors.New("invalid channel name or ID")
		}

		opts = append(opts, slack.MsgOptionText(
			fmt.Sprintf("Health check error on %s", addr),
			false))
		_, _, _, err = client.SendMessage(c.ID, opts...)

		return err
	}
}
