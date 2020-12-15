package member

import (
	"strings"

	"github.com/Starshine113/proxy/db"
	"github.com/Starshine113/proxy/router"
)

func displayName(ctx *router.Ctx) (err error) {
	if err = ctx.CheckSystem(); err != nil {
		return err
	}

	if ctx.CheckMinArgs(2); err != nil {
		return ctx.CommandError(err)
	}

	members, err := ctx.Database.GetAccountMembers(ctx.Author.ID)
	if err != nil {
		return ctx.CommandError(err)
	}

	member := ctx.Args[0]
	name := strings.Join(ctx.Args[1:], " ")

	if len(name) == 0 {
		_, err = ctx.Sendf("%v You must specify a new display name.", router.ErrorEmoji)
		return err
	} else if len(name) > 80 {
		_, err = ctx.Sendf("%v The display name you gave was too long (%v > 80 characters).", router.ErrorEmoji, len(name))
		return err
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

	err = ctx.Database.SetDisplayName(m.ID.String(), name)
	if err != nil {
		return ctx.CommandError(err)
	}

	_, err = ctx.Sendf("%v Member display name updated. This member will now be proxied with the name `%v.`", router.SuccessEmoji, name)
	return
}
