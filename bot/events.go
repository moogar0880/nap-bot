package bot

import (
	"fmt"

	"github.com/moogar0880/nap-bot/config"
	"github.com/nlopes/slack"
)

const (
	addEvent = "add"
)

func EmojiChangedEventHandler(e *slack.EmojiChangedEvent, client *slack.Client, m *ChannelIDManager, c config.Config) error {
	// we only care about new emojis, so fail gracefully if we get any other
	// kind of event
	if e.SubType != addEvent {
		return nil
	}

	for _, channel := range c.EmojiAddConfig.Channels {
		channelID, err := m.Get(channel)
		if err != nil {
			log.WithError(err).Error("unable to get channel id")
			continue
		}
		_, _, err = client.PostMessage(
			channelID,
			fmt.Sprintf("The following new emoji was just added :%s:", e.Name),
			slack.PostMessageParameters{},
		)
		if err != nil {
			log.WithError(err).Errorf("unable to send message to %s", channel)
		}
	}
	return nil
}
