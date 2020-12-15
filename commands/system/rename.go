package system

import (
	"strings"

	"github.com/Starshine113/proxy/router"
)

func rename(ctx *router.Ctx) (err error) {
	if err = ctx.CheckSystem(); err != nil {
		return err
	}

	name := strings.Join(ctx.Args, " ")
	if len(name) == 0 {
		_, err = ctx.Sendf("%v You must specify a new name.", router.ErrorEmoji)
		return
	} else if len(name) > 100 {
		_, err = ctx.Sendf("%v Your new system name would be too long (%v > 100 characters).", router.ErrorEmoji, len(name))
		return
	}

	s, err := ctx.Database.GetUserSystem(ctx.Author.ID)
	if err != nil {
		return ctx.CommandError(err)
	}

	err = ctx.Database.SetName(s.ID.String(), name)
	if err != nil {
		return ctx.CommandError(err)
	}

	_, err = ctx.Sendf("%v Successfully changed your system name to `%v`.", router.SuccessEmoji, name)
	return
}
