package commands

import "github.com/disgoorg/disgo/discord"

var Commands = []discord.ApplicationCommandCreate{
	rate,
	info,
	teams,
	leaderboard,
	history,
}
