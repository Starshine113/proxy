package router

import "errors"

// ErrorNoSystem is returned when the user has no system
var ErrorNoSystem = errors.New("user doesn't have a system")

// CheckSystem checks if the user has a system
func (ctx *Ctx) CheckSystem() (err error) {
	if !ctx.Database.HasSystem(ctx.Author.ID) {
		ctx.Sendf("%v You do not have a system registered with %v. To create one, use `%vsystem new`.", ErrorEmoji, ctx.BotUser.Username, ctx.Bot.Prefix)
		return ErrorNoSystem
	}

	return nil
}
