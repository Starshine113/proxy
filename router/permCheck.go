package router

// Discord proxy bot
// Copyright (C) 2020  Starshine System

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// Errors relating to missing permissions
var (
	ErrorMissingManagerPerms = errors.New("you are missing the `Manage Server` permission")
	ErrorMissingBotOwner     = errors.New("you are not a bot owner")
)

// Check checks if the user has permissions to run a command
func (ctx *Ctx) Check() (err error) {
	if ctx.Cmd.GuildOnly && ctx.Message.GuildID == "" {
		return ErrorNoDMs
	}
	if ctx.Cmd.Permissions == PermLevelManager {
		return ctx.checkManager(ctx.Author.ID)
	} else if ctx.Cmd.Permissions == PermLevelOwner {
		return checkOwner(ctx.Author.ID, ctx.Bot.Config.Bot.BotOwners)
	}
	return nil
}

func checkOwner(userID string, owners []string) (err error) {
	for _, u := range owners {
		if userID == u {
			return nil
		}
	}
	return ErrorMissingBotOwner
}

func (ctx *Ctx) checkManager(userID string) (err error) {
	// check if in DMs
	if ctx.Message.GuildID == "" {
		return ErrorNoDMs
	}

	// get the guild
	guild, err := ctx.Bot.Session.State.Guild(ctx.Message.GuildID)
	if err == discordgo.ErrStateNotFound {
		guild, err = ctx.Bot.Session.Guild(ctx.Message.GuildID)
	}
	if err != nil && err != discordgo.ErrStateNotFound {
		return err
	}

	// get the member
	member, err := ctx.Bot.Session.State.Member(ctx.Message.GuildID, ctx.Author.ID)
	if err == discordgo.ErrStateNotFound {
		member, err = ctx.Bot.Session.GuildMember(ctx.Message.GuildID, ctx.Author.ID)
	}
	if err != nil && err != discordgo.ErrStateNotFound {
		return err
	}

	// if the user is the guild owner, they have permission to use the command
	if member.User.ID == guild.OwnerID {
		return nil
	}

	// iterate through all guild roles
	for _, r := range guild.Roles {
		// iterate through member roles
		for _, u := range member.Roles {
			// if they have the role...
			if u == r.ID {
				// ...and the role has admin perms, return
				if checkPerms(r.Permissions, discordgo.PermissionAdministrator, discordgo.PermissionManageServer) {
					return nil
				}
			}
		}
	}

	return ErrorMissingManagerPerms
}

func checkPerms(p int, c ...int) bool {
	for _, perm := range c {
		if p&perm == perm {
			return true
		}
	}
	return false
}
