package demo

import (
	"context"
	"errors"
	"net/http"

	"github.com/adamwoolhether/htmx/go/app/hypermedia/view/dog"
	dogbus "github.com/adamwoolhether/htmx/go/business/dog"
	"github.com/adamwoolhether/htmx/go/foundation/web"
)

var selectedID string // terrible example...

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

func (g *Group) Form(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if selectedID != "" {
		found, ok := g.dogStore.Get(selectedID)
		if ok {
			return web.RenderHTML(ctx, w, dog.Form(&found), http.StatusOK)
		}
	}

	return web.RenderHTML(ctx, w, dog.Form(nil), http.StatusOK)
}

func (g *Group) CreateDog(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	name := r.PostFormValue("name")
	breed := r.PostFormValue("breed")

	newDog := g.dogStore.Add(name, breed)

	return web.RenderHTML(ctx, w, dog.Row(newDog, false), http.StatusOK)
}

func (g *Group) SelectDog(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	selectedID = r.PathValue("id")
	w.Header().Set("HX-Trigger", "selection-change")

	return web.RespondJSON(ctx, w, nil, http.StatusNoContent)
}

func (g *Group) DeselectDog(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	selectedID = ""
	w.Header().Set("HX-Trigger", "selection-change")

	return web.RespondJSON(ctx, w, nil, http.StatusNoContent)
}

func (g *Group) UpdateDog(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	name := r.PostFormValue("name")
	breed := r.PostFormValue("breed")

	selectedID = ""
	w.Header().Set("HX-Trigger", "selection-change")

	ok := g.dogStore.Update(id, name, breed)
	if !ok {
		return errors.New("failed to update dog")
	}

	return web.RenderHTML(ctx, w, dog.Rows(g.dogStore.GetAll()), http.StatusOK)
}

func (g *Group) DeleteDog(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	g.dogStore.Delete(id)

	return web.RenderDelete(ctx, w, http.StatusOK)
}
