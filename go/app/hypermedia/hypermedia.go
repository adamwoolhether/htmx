package hypermedia

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/adamwoolhether/htmx/go/app/hypermedia/handlers/demo"
	"github.com/adamwoolhether/htmx/go/business/dog"
	"github.com/adamwoolhether/htmx/go/foundation/web"
)

//go:embed static
var publicFS embed.FS

func StaticFS() http.Handler {
	// Create a sub-filesystem rooted at "static"
	fs, err := fs.Sub(publicFS, "static")
	if err != nil {
		log.Fatalf("failed to create sub FS: %v", err)
	}
	return http.FileServer(http.FS(fs)) // Serve the files directly
}

func Routes(app *web.App) {
	demoGrp := demo.NewGroup(dog.NewStore())
	app.Get("/rows", demoGrp.DogRows)
	app.Post("/dog", demoGrp.CreateDog)
	app.Delete("/dog/{id}", demoGrp.DeleteDog)
}
