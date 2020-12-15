package static

import (
	"fmt"
	"time"

	"github.com/Starshine113/proxy/router"
)

func ping(ctx *router.Ctx) (err error) {
	heartbeat := ctx.Bot.Session.HeartbeatLatency().Round(time.Millisecond).String()

	// get current time
	cmdStart := time.Now()

	// send initial message
	msg, err := ctx.Embedf("Pong!", "Heartbeat: %v", heartbeat)
	if err != nil {
		return fmt.Errorf("Ping: %w", err)
	}

	// get time difference, edit message
	diff := time.Now().Sub(cmdStart).Round(time.Millisecond).String()
	_, err = ctx.EditEmbedf(msg, "Pong!", "Heartbeat: %v\nMessage latency: %v", heartbeat, diff)
	return err
}
