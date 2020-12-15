package server

import "github.com/Starshine113/proxy/router"

// Init ...
func Init(r *router.Router) {
	r.AddGroup(&router.Group{
		Name: "Log",

		Description: "Manage the proxy log",

		Command: &router.Command{
			Name: "Channel",

			Description: "Set the log channel",
			Usage:       "[channel]",

			Command: log,
		},
	})
}
