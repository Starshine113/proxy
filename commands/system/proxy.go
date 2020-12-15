package system

import (
	"fmt"
	"github.com/Starshine113/proxy/router"
)

func proxy(ctx *router.Ctx) (err error) {
	if err = ctx.CheckSystem(); err != nil {
		return err
	}

	s, err := ctx.Database.GetUserSystem(ctx.Author.ID)
	if err != nil {
		return ctx.CommandError(err)
	}

	gs, err := ctx.Database.GetGuildSystem(ctx.Author.ID, ctx.Message.GuildID)
	if err != nil {
		return ctx.CommandError(err)
	}

	if len(ctx.Args) == 0 {
		var desc, title string
		guild, err := ctx.Bot.Session.State.Guild(ctx.Message.GuildID)
		if err == nil {
			title = "Proxy status for " + guild.Name
			desc = fmt.Sprintf("Proxying in this server (%v) is currently", guild.Name)
		} else {
			title = "Proxy status"
			desc = "Proxying in this server is currently"
		}

		if gs.ProxyEnabled {
			desc += fmt.Sprintf(" **enabled** for your system. To disable it, use `%vsystem proxy off`.", ctx.Bot.Prefix)
		} else {
			desc += fmt.Sprintf(" **disabled** for your system. To enable it, use `%vsystem proxy on`.", ctx.Bot.Prefix)
		}

		_, err = ctx.Embed(title, desc, 0)
		return err
	}

	if len(ctx.Args) > 0 {
		switch ctx.Args[0] {
		case "off", "disable":
			if !gs.ProxyEnabled {
				_, err = ctx.Sendf("%v Proxying in this server is already disabled for your system. To enable it, use `%vsystem proxy on`.", router.ErrorEmoji, ctx.Bot.Prefix)
				return err
			}
			err = ctx.Database.SetGuildProxy(s.ID.String(), ctx.Message.GuildID, false)
			if err != nil {
				return ctx.CommandError(err)
			}
			_, err = ctx.Sendf("%v Disabled proxying in this server for your system", router.SuccessEmoji)
		case "on", "enable":
			if gs.ProxyEnabled {
				_, err = ctx.Sendf("%v Proxying in this server is already enabled for your system. To disable it, use `%vsystem proxy off`.", router.ErrorEmoji, ctx.Bot.Prefix)
				return err
			}
			err = ctx.Database.SetGuildProxy(s.ID.String(), ctx.Message.GuildID, true)
			if err != nil {
				return ctx.CommandError(err)
			}
			_, err = ctx.Sendf("%v Enabled proxying in this server for your system", router.SuccessEmoji)
		}
	}

	return
}