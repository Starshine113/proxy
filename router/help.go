package router

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type cmdList []*Command

func (c cmdList) Len() int      { return len(c) }
func (c cmdList) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c cmdList) Less(i, j int) bool {
	return sort.StringsAreSorted([]string{c[i].Name, c[j].Name})
}

// Invite returns an invite link for the bot
func (ctx *Ctx) Invite() string {
	// perms is the list of permissions the bot will be granted by default
	var perms = discordgo.PermissionSendMessages +
		discordgo.PermissionManageMessages +
		discordgo.PermissionEmbedLinks +
		discordgo.PermissionAttachFiles +
		discordgo.PermissionReadMessageHistory +
		discordgo.PermissionUseExternalEmojis +
		discordgo.PermissionManageWebhooks +
		discordgo.PermissionAddReactions +
		discordgo.PermissionViewChannel

	return fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%v&permissions=%v&scope=bot", ctx.Bot.Session.State.User.ID, perms)
}

// Help is the help command
func (r *Router) Help(ctx *Ctx) (err error) {
	err = ctx.TriggerTyping()
	if err != nil {
		return
	}

	if len(ctx.Args) == 0 {
		permLevel := PermLevelNone

		if err = checkOwner(ctx.Author.ID, r.BotOwners); err == nil {
			permLevel = PermLevelOwner
		} else if err = ctx.checkManager(ctx.Author.ID); err == nil {
			permLevel = PermLevelManager
		}

		return r.details(ctx, permLevel)
	}

	var cmd *Command
	g := r.GetGroup(ctx.Args[0])
	if g != nil {
		if len(ctx.Args) == 1 {
			_, err = ctx.Send(ctx.GroupEmbed(g))
			return
		}
		cmd = g.GetCommand(ctx.Args[1])
		if cmd != nil {
			_, err = ctx.Send(ctx.GroupCmdEmbed(g, cmd))
			return
		}
	}
	cmd = r.GetCommand(ctx.Args[0])
	if cmd != nil {
		_, err = ctx.Send(ctx.CmdEmbed(cmd))
		return
	}

	_, err = ctx.Send(fmt.Sprintf("%v Invalid command or group provided:\n> `%v` is not a known command, group or alias.", ErrorEmoji, ctx.Args[0]))

	return
}

// GroupEmbed ...
func (ctx *Ctx) GroupEmbed(g *Group) *discordgo.MessageEmbed {
	var aliases string
	if g.Aliases == nil {
		aliases = "N/A"
	} else {
		aliases = strings.Join(g.Aliases, ", ")
	}

	var subCmds []string
	for _, cmd := range g.Subcommands {
		subCmds = append(subCmds, fmt.Sprintf("[%d] %s", cmd.Permissions, cmd.Name))
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("```%v```", strings.ToUpper(g.Name)),
		Description: g.Description,
		Color:       0x21a1a8,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Subcommands",
				Value:  fmt.Sprintf("```%v```", strings.Join(subCmds, "\n")),
				Inline: false,
			},
			{
				Name:   "Aliases",
				Value:  fmt.Sprintf("```%v```\n** **", aliases),
				Inline: false,
			},
			{
				Name:   "Default command",
				Value:  g.Command.Description,
				Inline: false,
			},
			{
				Name:   "Usage",
				Value:  fmt.Sprintf("```%v%v %v```", ctx.Bot.Config.Bot.Prefixes[0], strings.ToLower(g.Name), g.Command.Usage),
				Inline: false,
			},
			{
				Name:   "Permission level",
				Value:  "```" + g.Command.Permissions.String() + "```",
				Inline: false,
			},
		},
	}

	return embed
}

// GroupCmdEmbed ...
func (ctx *Ctx) GroupCmdEmbed(g *Group, cmd *Command) *discordgo.MessageEmbed {
	var aliases string

	if cmd.Aliases == nil {
		aliases = "N/A"
	} else {
		aliases = strings.Join(cmd.Aliases, ", ")
	}

	fields := make([]*discordgo.MessageEmbedField, 0)

	if cmd.LongDescription != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Description",
			Value:  cmd.LongDescription,
			Inline: false,
		})
	}

	fields = append(fields, []*discordgo.MessageEmbedField{
		{
			Name:   "Usage",
			Value:  fmt.Sprintf("```%v%v %v %v```", ctx.Bot.Config.Bot.Prefixes[0], strings.ToLower(g.Name), strings.ToLower(cmd.Name), cmd.Usage),
			Inline: false,
		},
		{
			Name:   "Aliases",
			Value:  fmt.Sprintf("```%v```", aliases),
			Inline: false,
		},
		{
			Name:   "Permission level",
			Value:  "```" + cmd.Permissions.String() + "```",
			Inline: false,
		},
	}...)

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("```%v %v```", strings.ToUpper(g.Name), strings.ToUpper(cmd.Name)),
		Description: cmd.Description,
		Color:       0x21a1a8,
		Fields:      fields,
	}

	return embed
}

// CmdEmbed ...
func (ctx *Ctx) CmdEmbed(cmd *Command) *discordgo.MessageEmbed {
	var aliases string

	if cmd.Aliases == nil {
		aliases = "N/A"
	} else {
		aliases = strings.Join(cmd.Aliases, ", ")
	}

	fields := make([]*discordgo.MessageEmbedField, 0)

	if cmd.LongDescription != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Description",
			Value:  cmd.LongDescription,
			Inline: false,
		})
	}

	fields = append(fields, []*discordgo.MessageEmbedField{
		{
			Name:   "Usage",
			Value:  fmt.Sprintf("```%v%v %v```", ctx.Bot.Config.Bot.Prefixes[0], strings.ToLower(cmd.Name), cmd.Usage),
			Inline: false,
		},
		{
			Name:   "Aliases",
			Value:  fmt.Sprintf("```%v```", aliases),
			Inline: false,
		},
		{
			Name:   "Permission level",
			Value:  "```" + cmd.Permissions.String() + "```",
			Inline: false,
		},
	}...)

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("```%v```", strings.ToUpper(cmd.Name)),
		Description: cmd.Description,
		Color:       0x21a1a8,
		Fields:      fields,
	}

	return embed
}

func (r *Router) details(ctx *Ctx, p PermLevel) (err error) {
	if err = ctx.TriggerTyping(); err != nil {
		return err
	}

	var cmds cmdList
	for _, c := range r.Commands {
		if c.Permissions <= p {
			cmds = append(cmds, c)
		}
	}

	for _, g := range r.Groups {
		if g.Command.Permissions <= p {
			cmds = append(cmds, &Command{
				Name:        g.Name,
				Permissions: g.Command.Permissions,
				Description: g.Command.Description,
			})
		}
	}

	sort.Sort(cmds)
	cmdSlices := make([][]*Command, 0)

	for i := 0; i < len(cmds); i += 10 {
		end := i + 10

		if end > len(cmds) {
			end = len(cmds)
		}

		cmdSlices = append(cmdSlices, cmds[i:end])
	}

	embeds := make([]*discordgo.MessageEmbed, 0)

	for i, c := range cmdSlices {
		x := make([]string, 0)
		for _, cmd := range c {
			x = append(x, fmt.Sprintf("`[%d] %v`: %v", cmd.Permissions, cmd.Name, cmd.Description))
		}
		embeds = append(embeds, &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name:    ctx.BotUser.Username + " help",
				IconURL: ctx.BotUser.AvatarURL("128"),
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Page %v/%v", i+1, len(cmdSlices)),
			},
			Timestamp:   time.Now().Format(time.RFC3339),
			Description: strings.Join(x, "\n"),
			Fields: []*discordgo.MessageEmbedField{{
				Name:   "Usage",
				Value:  "Use ⬅️ ➡️ to navigate between pages, and use ❌ to delete this message.",
				Inline: false,
			}},
			Color: 0x21a1a8,
		})
	}

	_, err = ctx.PagedEmbed(embeds)
	return
}

// PrettyDurationString ...
func PrettyDurationString(duration time.Duration) (out string) {
	var days, hours, hoursFrac, minutes float64

	hours = duration.Hours()
	hours, hoursFrac = math.Modf(hours)
	minutes = hoursFrac * 60

	hoursFrac = math.Mod(hours, 24)
	days = (hours - hoursFrac) / 24
	hours = hours - (days * 24)
	minutes = minutes - math.Mod(minutes, 1)

	if days != 0 {
		out += fmt.Sprintf("%v days, ", days)
	}
	if hours != 0 {
		out += fmt.Sprintf("%v hours, ", hours)
	}
	out += fmt.Sprintf("%v minutes", minutes)

	return
}
