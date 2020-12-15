package router

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Execute actually executes the router
func (r *Router) Execute(ctx *Ctx) (err error) {
	help := r.GetCommand("commands")
	if ctx.Match(append([]string{help.Name}, help.Aliases...)...) {
		err = r.Help(ctx)
		return
	}
	for _, g := range r.Groups {
		if ctx.Match(append([]string{g.Name}, g.Aliases...)...) || ctx.MatchRegexp(g.Regex) {
			if len(ctx.Args) == 0 {
				ctx.Command = ""
			} else {
				ctx.Command = ctx.Args[0]
			}
			if len(ctx.Args) > 1 {
				ctx.Args = ctx.Args[1:]
			} else {
				ctx.Args = []string{}
			}
			err = g.Execute(ctx)
			return
		}
	}
	for _, cmd := range r.Commands {
		if ctx.Match(append([]string{cmd.Name}, cmd.Aliases...)...) || ctx.MatchRegexp(cmd.Regex) {
			if len(ctx.Args) > 0 {
				if ctx.Args[0] == "help" || ctx.Args[0] == "usage" {
					ctx.Args[0] = ctx.Command
					err = r.Help(ctx)
					return
				}
			}
			ctx.Cmd = cmd
			if perms := ctx.Check(); perms != nil {
				return ctx.CommandError(perms)
			}
			if cmd.Cooldown != time.Duration(0) {
				if _, e := r.Cooldowns.Get(fmt.Sprintf("%v-%v-%v", ctx.Channel.ID, ctx.Author.ID, cmd.Name)); e == nil {
					_, err = ctx.Sendf("The command `%v` can only be run once every **%v**.", cmd.Name, cmd.Cooldown.Round(time.Millisecond).String())
					return err
				}
				err = r.Cooldowns.SetWithTTL(fmt.Sprintf("%v-%v-%v", ctx.Channel.ID, ctx.Author.ID, cmd.Name), true, cmd.Cooldown)
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

	_, err = ctx.Send(&discordgo.MessageSend{
		Content: fmt.Sprintf("%v Unknown command `%v`. For a list of commands, try `%vcommands`.", ErrorEmoji, ctx.Command, ctx.Bot.Prefix),
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{
				discordgo.AllowedMentionTypeUsers,
			},
		},
	})
	return
}
