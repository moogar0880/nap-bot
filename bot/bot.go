package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/moogar0880/nap-bot/config"
	"github.com/nlopes/slack"
)

const excludeArchived = true

// ChannelIDManager is responsible for mapping channel names to their IDs
type ChannelIDManager struct {
	channelIDsByName map[string]string
	client           *slack.Client
}

// NewChannelIDManager returns a newly initialized ChannelIDManager
func NewChannelIDManager() *ChannelIDManager {
	return &ChannelIDManager{
		channelIDsByName: make(map[string]string),
	}
}

// Init initializes the ChannelIDManager with the provided slack.Client
func (m *ChannelIDManager) Init(c *slack.Client) error {
	m.client = c
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	channels, err := m.client.GetChannelsContext(ctx, excludeArchived)
	if err != nil {
		return err
	}
	for _, channel := range channels {
		m.channelIDsByName[channel.Name] = channel.ID
	}
	return nil
}

// Get returns the id of a channel by name, if one exists in the map
func (m *ChannelIDManager) Get(name string) (string, error) {
	if channelID, ok := m.channelIDsByName[name]; ok {
		return channelID, nil
	}
	// TODO: update channel mapping if we can't find the requested channel
	return "", fmt.Errorf("bot: no such channel: %s", name)
}

// A Bot encapsulates all components of the slack bot
type Bot struct {
	Config         config.Config
	client         *slack.Client
	auth           *slack.AuthTestResponse
	channelManager *ChannelIDManager
	Done           chan bool
}

// New creates a newly initialized Bot instance based on the provided config
func New(c config.Config) *Bot {
	return &Bot{
		Config:         c,
		client:         slack.New(c.AuthToken),
		channelManager: NewChannelIDManager(),
		Done:           make(chan bool),
	}
}

// Run is the primary service to generate and kick off the slackbot listener
// This portion receives all incoming Real Time Messages (RTM) notices from the
// workspace as registered by the API token
func (b *Bot) Run(ctx context.Context) (err error) {
	InitializeLogger(b.Config)

	b.auth, err = b.client.AuthTest()
	if err != nil {
		return fmt.Errorf("bot: authentication failed: %s", err.Error())
	}

	// initialize our mapping of client names to ids
	if err = b.channelManager.Init(b.client); err != nil {
		return err
	}

	log.WithFields(logrus.Fields{
		"user":    b.auth.User,
		"user_id": b.auth.UserID,
	}).Debugf("successfully authenticated")

	go b.run(ctx)
	return nil
}

func (b *Bot) run(ctx context.Context) {
	rtm := b.client.NewRTM()
	go rtm.ManageConnection()

	log.Debug("listening for incoming events")
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.EmojiChangedEvent:
			EmojiChangedEventHandler(ev, b.client, b.channelManager, b.Config)
		case *slack.RTMError:
			log.WithField("error", ev.Error()).Error("an error occurred")
		}
	}
	b.Done <- true
}
