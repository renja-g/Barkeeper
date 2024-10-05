package utils

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

func HasAdminRole(member *discord.Member, adminRoleID snowflake.ID) bool {
	if member == nil {
		return false
	}
	for _, roleID := range member.RoleIDs {
		if roleID == adminRoleID {
			return true
		}
	}
	return false
}
