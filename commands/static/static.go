package static

import "github.com/Starshine113/proxy/router"

// Init ...
func Init(r *router.Router) {
	r.AddCommand(&router.Command{
		Name:        "ping",
		Description: "Show the bot's latency",

		Permissions: router.PermLevelNone,

		Command: ping,
	})
}
