package member

// Discord proxy bot
// Copyright (C) 2020  Starshine System

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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

	g.AddCommand(&router.Command{
		Name:    "Delete",
		Aliases: []string{"Yeet"},

		Description: "Delete a member",
		Usage:       "<name/ID>",

		Command: memberDelete,
	})
}
