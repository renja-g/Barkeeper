package commands

import "github.com/disgoorg/disgo/discord"

var Commands = []discord.ApplicationCommandCreate{
	test,
	version,
	rate,
	info,
	teams,
	leaderboard,
	history,
}
