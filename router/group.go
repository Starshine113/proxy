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
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
)

// Group is a group of subcommands
type Group struct {
	Name        string
	Aliases     []string
	Regex       *regexp.Regexp
	Description string
	Command     *Command
	Subcommands []*Command
	Router      *Router
	Cooldowns   *ttlcache.Cache
}

// AddGroup adds a group to the router
func (r *Router) AddGroup(group *Group) *Group {
	cache := ttlcache.NewCache()
	cache.SkipTTLExtensionOnHit(true)

	group.Router = r
	group.Cooldowns = cache
	r.Groups = append(r.Groups, group)
	return r.GetGroup(group.Name)
}

// AddCommand adds a command to a group
func (g *Group) AddCommand(cmd *Command) {
	cmd.Router = g.Router
	if cmd.Cooldown == 0 {
		cmd.Cooldown = 500 * time.Millisecond
	}
	g.Subcommands = append(g.Subcommands, cmd)
}

// GetGroup returns a group by name
func (r *Router) GetGroup(name string) (group *Group) {
	name = strings.ToLower(name)
	for _, group := range r.Groups {
		if strings.ToLower(group.Name) == name {
			return group
		}
		for _, a := range group.Aliases {
			if strings.ToLower(a) == name {
				return group
			}
		}
	}
	return nil
}

// GetCommand gets a command by name
func (g *Group) GetCommand(name string) (c *Command) {
	for _, cmd := range g.Subcommands {
		if strings.ToLower(cmd.Name) == strings.ToLower(name) {
			return cmd
		}
		for _, a := range cmd.Aliases {
			if strings.ToLower(a) == strings.ToLower(name) {
				return cmd
			}
		}
	}
	if strings.ToLower(g.Command.Name) == strings.ToLower(name) {
		return g.Command
	}
	for _, a := range g.Command.Aliases {
		if strings.ToLower(a) == strings.ToLower(name) {
			return g.Command
		}
	}
	return nil
}

// Execute executes any command given
func (g *Group) Execute(ctx *Ctx) (err error) {
	g.Subcommands = append(g.Subcommands, g.Command)
	if ctx.Command == "help" || ctx.Command == "usage" {
		if len(ctx.Args) > 0 {
			ctx.Args[0] = g.Name
		} else {
			ctx.Args = []string{g.Name}
		}
		err = g.Router.Help(ctx)
		return
	}
	for _, cmd := range g.Subcommands {
		if ctx.Match(append([]string{cmd.Name}, cmd.Aliases...)...) || ctx.MatchRegexp(cmd.Regex) {
			if len(ctx.Args) > 0 {
				if ctx.Args[0] == "help" || ctx.Args[0] == "usage" {
					ctx.Args[0] = g.Name
					if len(ctx.Args) > 1 {
						ctx.Args[1] = ctx.Command
					} else {
						ctx.Args = append(ctx.Args, ctx.Command)
					}
					err = g.Router.Help(ctx)
					return
				}
			}

			ctx.Cmd = cmd
			if perms := ctx.Check(); perms != nil {
				ctx.CommandError(perms)
				return nil
			}
			if cmd.Cooldown != time.Duration(0) {
				if _, e := g.Cooldowns.Get(fmt.Sprintf("%v-%v-%v", ctx.Channel.ID, ctx.Author.ID, cmd.Name)); e == nil {
					_, err = ctx.Sendf("The command `%v` can only be run once every **%v**.", cmd.Name, cmd.Cooldown.Round(time.Millisecond).String())
					return err
				}
				err = g.Cooldowns.SetWithTTL(fmt.Sprintf("%v-%v-%v", ctx.Channel.ID, ctx.Author.ID, cmd.Name), true, cmd.Cooldown)
				if err != nil {
					return err
				}
			}
			err = cmd.Command(ctx)
			if err != ErrorNoSystem {
				return err
			}
			return nil
		}
	}
	ctx.Cmd = g.Command
	ctx.Args = append([]string{ctx.Command}, ctx.Args...)
	ctx.Command = g.Name
	if perms := ctx.Check(); perms != nil {
		ctx.CommandError(perms)
		return nil
	}
	err = g.Command.Command(ctx)
	if err != ErrorNoSystem {
		return err
	}
	return nil
}
