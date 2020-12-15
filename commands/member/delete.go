package member

import (
	"github.com/Starshine113/proxy/db"
	"github.com/Starshine113/proxy/router"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func memberDelete(ctx *router.Ctx) (err error) {
	if err = ctx.CheckSystem(); err != nil {
		return err
	}

	members, err := ctx.Database.GetAccountMembers(ctx.Author.ID)
	if err != nil {
		return ctx.CommandError(err)
	}

	member := strings.Join(ctx.Args, " ")
	if member == "" {
		_, err = ctx.Sendf("%v You must supply a member to delete.", router.ErrorEmoji)
		return
	}

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

	if m == nil {
		_, err = ctx.Sendf("%v Member `%v` not found.", router.ErrorEmoji, member)
		return
	}

	_, err = ctx.Sendf("%v Are you sure you want to delete this member? If so, reply with this member's ID (`%s`).\n*This action is irreversible!*", router.WarnEmoji, m.ID)

	ctx.AddMessageHandler(ctx.Message.ID, func(ctx *router.Ctx, msg *discordgo.MessageCreate) {
		c := strings.TrimSpace(msg.Content)

		if c != m.ID.String() {
			_, err = ctx.Sendf("%v Member deletion cancelled. Note that you must reply with this member's ID (`%s`) *verbatim*.", router.ErrorEmoji, m.ID)
			return
		}
		err = ctx.Database.DeleteMember(c)
		if err != nil {
			ctx.CommandError(err)
			return
		}
		_, err = ctx.Sendf("%v Member deleted.", router.SuccessEmoji)
		return
	})

	return
}