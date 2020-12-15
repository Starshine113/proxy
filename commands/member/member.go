package member

import "github.com/Starshine113/proxy/router"

// Init ...
func Init(r *router.Router) {
	g := r.AddGroup(&router.Group{
		Name:    "Member",
		Aliases: []string{"M"},

		Description: "Commands to view and manage member",

		Command: &router.Command{
			Name:    "View",
			Aliases: []string{"Show"},

			Description: "View a member",
			Usage:       "<member name/ID>",

			Command: view,
		},
	})

	g.AddCommand(&router.Command{
		Name: "New",

		Description: "Create a new member with the given name",
		Usage:       "<name>",

		Command: new,
	})

	g.AddCommand(&router.Command{
		Name: "Proxy",

		Description: "Set a member's proxy",
		Usage:       "<name> <proxy>",

		Command: proxy,
	})

	g.AddCommand(&router.Command{
		Name:    "Avatar",
		Aliases: []string{"Av", "Pfp"},

		Description: "Set or view a member's avatar",
		Usage:       "<name> [avatar]",

		Command: avatar,
	})

	g.AddCommand(&router.Command{
		Name:    "DisplayName",
		Aliases: []string{"Nickname", "Nick", "DN"},

		Description: "Set a member's display name",
		Usage:       "<name> <display name>",

		Command: displayName,
	})
}
