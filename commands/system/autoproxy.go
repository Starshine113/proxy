package system

import (
	"github.com/Starshine113/proxy/db"
	"github.com/Starshine113/proxy/router"
)

func autoproxy(ctx *router.Ctx) (err error) {
	if err = ctx.CheckSystem(); err != nil {
		return err
	}

	if err = ctx.CheckArgRange(0, 1); err != nil {
		return ctx.CommandError(err)
	}

	gs, err := ctx.Database.GetGuildSystem(ctx.Author.ID, ctx.Message.GuildID)
	if err != nil {
		return ctx.CommandError(err)
	}

	if len(ctx.Args) == 0 {
		if gs.AutoproxyMode == db.AutoproxyModeOff {
			_, err = ctx.Sendf("Autoproxying in this server is currently **disabled** for your system.\nTo enable it, use `%vautoproxy <mode>` where `<mode>` is one of `latch`, `front`, or a specific member.", ctx.Bot.Prefix)
			return err
		}
		_, err = ctx.Sendf("Autoproxy for this server is set to **%v** mode.\nTo disable autoproxy, use `%vautoproxy off`.", gs.AutoproxyMode, ctx.Bot.Prefix)
		return
	}

	switch ctx.Args[0] {
	case "latch":
		err = ctx.Database.SetAutoproxyMode(ctx.Author.ID, ctx.Message.GuildID, db.AutoproxyModeLatch)
		if err != nil {
			return ctx.CommandError(err)
		}
		_, err = ctx.Sendf("%v Autoproxy set to **latch** mode in this server.", router.SuccessEmoji)
		return err
	case "front":
		err = ctx.Database.SetAutoproxyMode(ctx.Author.ID, ctx.Message.GuildID, db.AutoproxyModeFront)
		if err != nil {
			return ctx.CommandError(err)
		}
		_, err = ctx.Sendf("%v Autoproxy set to **front** mode in this server.", router.SuccessEmoji)
		return err
	case "off":
		err = ctx.Database.SetAutoproxyMode(ctx.Author.ID, ctx.Message.GuildID, db.AutoproxyModeOff)
		if err != nil {
			return ctx.CommandError(err)
		}
		_, err = ctx.Sendf("%v Disabled autoproxy in this server.", router.SuccessEmoji)
		return err
	}

	return
}
