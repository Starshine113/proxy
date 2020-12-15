package member

import (
	"strings"

	"github.com/Starshine113/proxy/router"
)

func new(ctx *router.Ctx) (err error) {
	if err = ctx.CheckSystem(); err != nil {
		return err
	}

	name := strings.Join(ctx.Args, " ")
	if name == "" {
		_, err = ctx.Sendf("%v You must pass a name.", router.ErrorEmoji)
		return
	}

	if len(name) > 100 {
		_, err = ctx.Sendf("%v The name you gave is too long (%v > 100 characters).", router.ErrorEmoji, len(name))
		return
	}

	s, err := ctx.Database.GetUserSystem(ctx.Author.ID)
	if err != nil {
		return ctx.CommandError(err)
	}

	m, err := ctx.Database.NewMember(s.ID.String(), name)
	if err != nil {
		return ctx.CommandError(err)
	}

	_, err = ctx.Sendf("%v Member \"%v\" (`%v`) registered!\nTo get started proxying this member, use the command `%vmember proxy %v`.", router.SuccessEmoji, m.Name, m.ID, ctx.Bot.Prefix, m.ID)
	return
}
