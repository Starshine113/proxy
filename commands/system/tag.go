package system

import (
	"strings"

	"github.com/Starshine113/proxy/router"
)

func tag(ctx *router.Ctx) (err error) {
	tag := strings.Join(ctx.Args, " ")

	if err = ctx.CheckSystem(); err != nil {
		return err
	}

	s, err := ctx.Database.GetUserSystem(ctx.Author.ID)
	if err != nil {
		return ctx.CommandError(err)
	}

	if tag == "" {
		if s.Tag == "" {
			_, err = ctx.Sendf("%v You do not currently have a system tag set. To set one, use `%vsystem tag <new tag>`.", router.ErrorEmoji, ctx.Bot.Prefix)
			return
		}
		_, err = ctx.Sendf("Your current system tag is `%v`. To change it, use `%vsystem tag <new tag>`.", s.Tag, ctx.Bot.Prefix)
		return
	}

	if len(tag) > 100 {
		_, err = ctx.Sendf("%v Your new system tag would be too long (%v > 100 characters).", router.ErrorEmoji, len(tag))
		return
	}

	err = ctx.Database.SetTag(s.ID.String(), tag)
	if err != nil {
		return ctx.CommandError(err)
	}

	_, err = ctx.Sendf("%v Succesfully changed your system tag to `%v`.", router.SuccessEmoji, tag)
	return
}
