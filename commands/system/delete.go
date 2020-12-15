package system

import (
	"strings"

	"github.com/Starshine113/proxy/router"
	"github.com/bwmarrin/discordgo"
)

func delete(ctx *router.Ctx) (err error) {
	if err = ctx.CheckSystem(); err != nil {
		return err
	}

	s, err := ctx.Database.GetUserSystem(ctx.Author.ID)
	if err != nil {
		return ctx.CommandError(err)
	}

	_, err = ctx.Sendf("%v Are you sure you want to delete your system? If so, reply with your system ID (`%s`).\n*This action is irreversible!*", router.WarnEmoji, s.ID)

	ctx.AddMessageHandler(ctx.Message.ID, func(ctx *router.Ctx, m *discordgo.MessageCreate) {
		c := strings.TrimSpace(m.Content)

		if c != s.ID.String() {
			_, err = ctx.Sendf("%v System deletion cancelled. Note that you must reply with your system ID (`%s`) *verbatim*.", router.ErrorEmoji, s.ID)
			return
		}
		err = ctx.Database.DeleteSystem(c)
		if err != nil {
			ctx.CommandError(err)
			return
		}
		_, err = ctx.Sendf("%v System deleted.", router.SuccessEmoji)
		return
	})

	return
}
