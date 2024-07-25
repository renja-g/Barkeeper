package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/handler"

	dbot "github.com/renja-g/Barkeeper"
	"github.com/renja-g/Barkeeper/commands"
	"github.com/renja-g/Barkeeper/components"
)

var (
	shouldSyncCommands *bool
	version            = "dev"
)

func init() {
	shouldSyncCommands = flag.Bool("sync-commands", false, "Whether to sync commands to discord")
	flag.Parse()
}

func main() {
	cfg, err := dbot.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	logLevel := slog.LevelInfo
	if cfg.LogLevel != 0 {
		logLevel = slog.Level(cfg.LogLevel)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	logger.Info("Starting bot", "version", version)
	logger.Info("Syncing commands?", "value", *shouldSyncCommands)

	b := dbot.New(logger, version, *cfg)

	h := handler.New()
	h.Command("/rate", commands.RateHandler)
	h.Autocomplete("/rate", commands.RateAutocompleteHandler)
	h.Command("/info", commands.InfoHandler)
	h.Command("/teams", func(e *handler.CommandEvent) error {
		return commands.TeamsHandler(e, b)
	})
	h.Command("/leaderboard", commands.LeaderboardHandler)
	h.Command("/history", commands.HistoryHandler)
	h.Command("/list", func(e *handler.CommandEvent) error {
		return commands.ListHandler(e, b)
	})

	h.Component("/reshuffle_button", components.ReshuffleComponent)
	h.Component("/start_match_button", components.StartMatchComponent)
	h.Component("/team1_wins_button", components.SetWinnerComponent)
	h.Component("/team2_wins_button", components.SetWinnerComponent)
	h.Component("/cancel_match_button", components.CancelMatchComponent)

	b.SetupBot(h, bot.NewListenerFunc(b.OnReady))

	if *shouldSyncCommands {
		if cfg.DevMode {
			logger.Info("Syncing commands in dev mode")
			_, err = b.Client.Rest().SetGuildCommands(b.Client.ApplicationID(), cfg.DevGuildID, commands.Commands)
		} else {
			logger.Info("Syncing commands")
			_, err = b.Client.Rest().SetGlobalCommands(b.Client.ApplicationID(), commands.Commands)
		}
		if err != nil {
			logger.Error("Failed to sync commands", "error", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = b.Client.OpenGateway(ctx); err != nil {
		logger.Error("Failed to connect to gateway", "error", err)
	}
	defer b.Client.Close(context.TODO())

	logger.Info("Bot is running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
	logger.Info("Shutting down...")
}