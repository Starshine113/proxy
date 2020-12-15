package member

import (
	"fmt"
	"strings"
	"time"

	"github.com/Starshine113/proxy/db"
	"github.com/Starshine113/proxy/router"
	"github.com/bwmarrin/discordgo"
)

func view(ctx *router.Ctx) (err error) {
	member := strings.Join(ctx.Args, " ")
	var m *db.Member
	if len(member) == 0 {
		_, err = ctx.Sendf("%v You must specify a member (name or ID).", router.ErrorEmoji)
		return
	}
	s, err := ctx.Database.GetUserSystem(ctx.Author.ID)
	if err == nil {
		members, err := ctx.Database.GetSystemMembers(s.ID.String())
		if err != nil {
			return ctx.CommandError(err)
		}
		for _, mem := range members {
			if member == mem.ID.String() {
				m = mem
				break
			}
			if strings.ToLower(member) == strings.ToLower(mem.Name) {
				m = mem
				break
			}
		}
	}

	if m != nil {
		_, err = ctx.Send(memberCard(s, m))
		return err
	}

	return
}

func memberCard(s *db.System, m *db.Member) *discordgo.MessageEmbed {
	title := m.Name
	if s.Name != "" {
		title += " (" + s.Name + ")"
	}

	fields := make([]*discordgo.MessageEmbedField, 0)

	if m.DisplayName != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Display name",
			Value:  m.DisplayName,
			Inline: true,
		})
	}

	if m.Prefix != "" || m.Suffix != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Proxy",
			Value:  fmt.Sprintf("`%vtext%v`", m.Prefix, m.Suffix),
			Inline: true,
		})
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "** **",
		Value:  fmt.Sprintf("Member ID: `%s`\nSystem ID: `%s`", m.ID, s.ID),
		Inline: false,
	})

	return &discordgo.MessageEmbed{
		Title: title,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: m.AvatarURL,
		},
		Fields: fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Created at",
		},
		Timestamp: m.Created.UTC().Format(time.RFC3339),
	}
}
