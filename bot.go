package dbot

import (
	"context"

	"log/slog"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/paginator"
)

func New(logger *slog.Logger, version string, config Config) *Bot {
    return &Bot{
        Logger:    logger,
        Config:    config,
        Paginator: paginator.New(),
        Version:   version,
    }
}

type Bot struct {
    Logger    *slog.Logger
    Client    bot.Client
    Paginator *paginator.Manager
    Config    Config
    Version   string
}

func (b *Bot) SetupBot(listeners ...bot.EventListener) {
	var err error
	b.Client, err = disgo.New(b.Config.Token,
		bot.WithLogger(b.Logger),
		bot.WithGatewayConfigOpts(gateway.WithIntents(
			gateway.IntentGuilds,
			gateway.IntentGuildMessages,
			gateway.IntentMessageContent,
			gateway.IntentGuildVoiceStates)),
		bot.WithCacheConfigOpts(cache.WithCaches(
			cache.FlagGuilds,
			cache.FlagVoiceStates,
			cache.FlagChannels)),
		bot.WithEventListeners(b.Paginator),
		bot.WithEventListeners(listeners...),
	)
	if err != nil {
		b.Logger.Error("Failed to setup b: ", err)
	}
}

func (b *Bot) OnReady(_ *events.Ready) {
	b.Logger.Info("Butler ready")
	if err := b.Client.SetPresence(context.TODO(), gateway.WithListeningActivity("you"), gateway.WithOnlineStatus(discord.OnlineStatusOnline)); err != nil {
		b.Logger.Error("Failed to set presence: %s", err)
	}
}
