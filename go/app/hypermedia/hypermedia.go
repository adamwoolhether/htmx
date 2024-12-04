package hypermedia

import (
	"github.com/adamwoolhether/htmx/go/foundation/logger"
	"github.com/adamwoolhether/htmx/go/foundation/web"
)

type Config struct {
	Log *logger.Logger
}

func Routes(app *web.App, cfg Config) {

}

func webRoutes(app *web.App, cfg Config) {
	const root = ""

}
