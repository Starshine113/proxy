package system

import (
	"strings"

	"github.com/Starshine113/proxy/router"
)

// Init ...
func Init(r *router.Router) {
	g := r.AddGroup(&router.Group{
		Name:    "System",
		Aliases: []string{"S"},

		Description: "Commands to view and manage systems",

		Command: &router.Command{
			Name:    "View",
			Aliases: []string{"Show"},

			Description: "View a system",
			Usage:       "[user/system ID]",

			Command: view,
		},
	})

	g.AddCommand(&router.Command{
		Name: "New",

		Description: "Create a new system with an optional name",
		Usage:       "[name]",

		Command: new,
	})

	g.AddCommand(&router.Command{
		Name: "Tag",

		Description: "View your system's current tag, or change it",
		Usage:       "[new tag]",

		Command: tag,
	})

	g.AddCommand(&router.Command{
		Name:    "Delete",
		Aliases: []string{"Yeet"},

		Description: "Delete your system",

		Command: delete,
	})

	g.AddCommand(&router.Command{
		Name: "Rename",

		Description: "Rename your system",

		Command: rename,
	})

	g.AddCommand(&router.Command{
		Name:    "List",
		Aliases: []string{"L"},

		Description: "List members in your system",

		Command: list,
	})
}

func new(ctx *router.Ctx) (err error) {
	if ctx.Database.HasSystem(ctx.Author.ID) {
		_, err = ctx.Sendf("%v You already have a system registered with %v. To view it, use `%vsystem`.", router.ErrorEmoji, ctx.BotUser.Username, ctx.Bot.Prefix)
		return err
	}

	name := strings.Join(ctx.Args, " ")

	if len(name) > 100 {
		_, err = ctx.Sendf("%v The name you gave is too long (%v > 100 characters).", router.ErrorEmoji, len(name))
		return
	}

	s, err := ctx.Database.CreateSystem(ctx.Author.ID, name)
	if err != nil {
		return ctx.CommandError(err)
	}

	_, err = ctx.Sendf("%v Your system has been created. Type `%vsystem %v` to view it.", router.SuccessEmoji, ctx.Bot.Prefix, s.ID.String())
	return
}
