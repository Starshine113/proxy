package system

import "github.com/Starshine113/proxy/router"

func autoproxy(ctx *router.Ctx) (err error) {
	if err = ctx.CheckSystem(); err != nil {
		return err
	}

	if err = ctx.CheckArgRange(0, 1); err != nil {
		return ctx.CommandError(err)
	}
	return
}
