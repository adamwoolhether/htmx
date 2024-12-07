package demo

import (
	"context"
	"net/http"

	"github.com/adamwoolhether/htmx/go/app/hypermedia/view/dog"
	dogbus "github.com/adamwoolhether/htmx/go/business/dog"
	"github.com/adamwoolhether/htmx/go/foundation/web"
)

type Group struct {
	dogStore *dogbus.Dogs
}

func NewGroup(dogStore *dogbus.Dogs) *Group {
	// Add default
	dogStore.Add("Comet", "Whippet")
	dogStore.Add("Oscar", "German Shorthaired Pointer")
	return &Group{dogStore: dogStore}
}

func (g *Group) DogRows(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	allDogs := g.dogStore.GetAll()

	return web.RenderHTML(ctx, w, dog.Rows(allDogs), http.StatusOK)
}

func (g *Group) CreateDog(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	name := r.PostFormValue("name")
	breed := r.PostFormValue("breed")

	newDog := g.dogStore.Add(name, breed)

	return web.RenderHTML(ctx, w, dog.Row(newDog), http.StatusOK)
}

func (g *Group) DeleteDog(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	g.dogStore.Delete(id)

	return web.RespondJSON(ctx, w, nil, http.StatusNoContent)
}
