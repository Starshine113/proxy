package commands

import (
	"github.com/Starshine113/proxy/commands/member"
	"github.com/Starshine113/proxy/commands/static"
	"github.com/Starshine113/proxy/commands/system"
	"github.com/Starshine113/proxy/router"
)

// Init ...
func Init(r *router.Router) {
	static.Init(r)
	system.Init(r)
	member.Init(r)
}
