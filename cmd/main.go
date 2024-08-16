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
		logLevel = cfg.LogLevel
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	logger.Info("Starting bot", "version", version)
	logger.Info("Syncing commands?", "value", *shouldSyncCommands)

	b := dbot.New(logger, version, *cfg)

	h := handler.New()
	h.SlashCommand("/rate", commands.RateHandler())
	h.Autocomplete("/rate", commands.RateAutocompleteHandler)
	h.SlashCommand("/info", commands.InfoHandler())
	h.SlashCommand("/teams", commands.TeamsHandler(b))
	h.SlashCommand("/leaderboard", commands.LeaderboardHandler())
	h.SlashCommand("/history", commands.HistoryHandler(b))
	h.SlashCommand("/list", commands.ListHandler(b))
	h.SlashCommand("/invite", commands.InviteHandler(b))
	h.SlashCommand("/help", commands.HelpHandler())
	h.SlashCommand("/link_account", commands.LinkAccountHandler(cfg))

	h.ButtonComponent("/reshuffle_button", components.ReshuffleComponent())
	h.ButtonComponent("/start_match_button", components.StartMatchComponent(cfg))
	h.ButtonComponent("/team1_wins_button", components.SetWinnerComponent(cfg))
	h.ButtonComponent("/team2_wins_button", components.SetWinnerComponent(cfg))
	h.ButtonComponent("/cancel_match_button", components.CancelMatchComponent(cfg))
	h.ButtonComponent("/accept_the_invite_button", components.AcceptTheInviteComponent())
	h.ButtonComponent("/verify_acc/{data}", components.VerifyAccountLinkComponent(cfg))

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
