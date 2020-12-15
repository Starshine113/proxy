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
	"strings"

	"github.com/Starshine113/proxy/bot"
	"github.com/Starshine113/proxy/db"
	"github.com/bwmarrin/discordgo"
)

const (
	// SuccessEmoji is the emoji used to designate a successful action
	SuccessEmoji = "✅"
	// ErrorEmoji is the emoji used for errors
	ErrorEmoji = "❌"
	// WarnEmoji is the emoji used to warn that a command failed
	WarnEmoji = "⚠️"
)

// Ctx is the context for a command
type Ctx struct {
	Command string
	Args    []string
	RawArgs string

	Bot      *bot.Bot
	BotUser  *discordgo.User
	Database *db.Db

	Message *discordgo.MessageCreate
	Channel *discordgo.Channel
	Author  *discordgo.User

	Cmd *Command

	AdditionalParams map[string]interface{}
}

// Errors when creating Context
var (
	ErrorNoBotUser = errors.New("bot user not found in state cache")
)

// Context creates a new Ctx
func Context(prefixes []string, messageContent string, m *discordgo.MessageCreate, b *bot.Bot) (ctx *Ctx, err error) {
	messageContent = TrimPrefixesSpace(messageContent, prefixes...)
	message := strings.Split(messageContent, " ")
	command := message[0]
	args := []string{}
	if len(message) > 1 {
		args = message[1:]
	}

	ctx = &Ctx{Command: command, Args: args, Message: m, Author: m.Author, Database: b.Db, Bot: b, RawArgs: strings.Join(args, " ")}

	channel, err := b.Session.Channel(m.ChannelID)
	if err != nil {
		return ctx, err
	}
	ctx.Channel = channel

	ctx.BotUser = b.Session.State.User
	if ctx.BotUser == nil {
		return ctx, ErrorNoBotUser
	}

	ctx.AdditionalParams = make(map[string]interface{})

	return ctx, nil
}
