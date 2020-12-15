package member

import (
	"fmt"
	"strings"

	"github.com/Starshine113/proxy/db"
	"github.com/Starshine113/proxy/router"
	"github.com/bwmarrin/discordgo"
)

func avatar(ctx *router.Ctx) (err error) {
	if err = ctx.CheckSystem(); err != nil {
		return err
	}

	if err = ctx.CheckMinArgs(1); err != nil {
		return ctx.CommandError(err)
	}

	members, err := ctx.Database.GetAccountMembers(ctx.Author.ID)
	if err != nil {
		return ctx.CommandError(err)
	}

	member := ctx.Args[0]

	var m *db.Member

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

	if len(ctx.Args) == 1 {
		if len(ctx.Message.Attachments) == 0 {
			return viewAvatar(m, ctx)
		}
		return changeAvatar(m, ctx)
	}

	return changeAvatar(m, ctx)
}

func changeAvatar(m *db.Member, ctx *router.Ctx) (err error) {
	if len(ctx.Args) > 1 {
		member, err := ctx.ParseMember(strings.Join(ctx.Args[1:], " "))
		if err == nil {
			err = ctx.Database.SetAvatar(m.ID.String(), member.User.AvatarURL("256"))
			if err != nil {
				return ctx.CommandError(err)
			}
			_, err = ctx.Sendf("%v Member avatar changed to %v's avatar.\nNote that if %v changes their avatar, this member's avatar will need to be reset.", router.SuccessEmoji, member.Mention(), member.User.Username)
			return err
		}

		// it's not a member (or the fetch failed), so try if it's an image
		if router.HasAnySuffix(ctx.Args[1], ".jpg", ".jpeg", ".png", ".gif", ".webp") {
			err = ctx.Database.SetAvatar(m.ID.String(), ctx.Args[1])
			if err != nil {
				return ctx.CommandError(err)
			}
			_, err = ctx.Send(&discordgo.MessageSend{
				Content: fmt.Sprintf("%v Member avatar changed to the image at the given URL.", router.SuccessEmoji),
				Embed: &discordgo.MessageEmbed{
					Title: fmt.Sprintf("%v's avatar", m.Name),
					Image: &discordgo.MessageEmbedImage{
						URL: ctx.Args[1],
					},
				},
			})
			return err
		}
	}

	if len(ctx.Message.Attachments) == 0 {
		return nil
	}

	a := ctx.Message.Attachments[0]

	if router.HasAnySuffix(a.URL, ".jpg", ".jpeg", ".png", ".gif", ".webp") {
		err = ctx.Database.SetAvatar(m.ID.String(), a.URL)
		if err != nil {
			return ctx.CommandError(err)
		}
		_, err = ctx.Send(&discordgo.MessageSend{
			Content: fmt.Sprintf("%v Member avatar changed to the attached image.\nNote that if this message is deleted, the member's avatar will need to be reset.", router.SuccessEmoji),
			Embed: &discordgo.MessageEmbed{
				Title: fmt.Sprintf("%v's avatar", m.Name),
				Image: &discordgo.MessageEmbedImage{
					URL: a.URL,
				},
			},
		})
		return err
	}

	return
}

func viewAvatar(m *db.Member, ctx *router.Ctx) (err error) {
	if m.AvatarURL == "" {
		_, err = ctx.Sendf("%v has no avatar set.", m.Name)
		return err
	}

	_, err = ctx.Send(&discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%v's avatar", m.Name),
		Description: fmt.Sprintf("To clear, use `%vmember avatar %s clear`.", ctx.Bot.Prefix, m.ID),
		Image: &discordgo.MessageEmbedImage{
			URL: m.AvatarURL,
		},
	})
	return err
}
