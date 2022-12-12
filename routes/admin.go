package routes

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/oxodao/photobooth/models"
	"github.com/oxodao/photobooth/services"
	"github.com/oxodao/photobooth/utils"
)

func registerAdminRoutes(r *mux.Router) {
	r.Use(AuthenticatedMiddleware)

	r.HandleFunc("/login", isPasswordValid)
	r.HandleFunc("/set_mode/{mode}", setMode).Methods(http.MethodPost)

	eventRouter := r.PathPrefix("/event/{event}").Subrouter()
	eventRouter.Use(WithEventMiddleware)

	imageRouter := eventRouter.PathPrefix("/image/{image}").Subrouter()
	imageRouter.Use(WithImageMiddleware)

	imageRouter.HandleFunc("", serveImage).Methods(http.MethodGet)
}

func isPasswordValid(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("yes"))
}

func setMode(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	(*services.GET.MqttClient).Publish("photobooth/admin/set_mode", 2, false, values["mode"])
}

func serveImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	event := ctx.Value(ContextEventKey).(*models.Event)
	image := ctx.Value(ContextImageKey).(*models.Image)

	subpath := "pictures"
	if image.Unattended {
		subpath = "unattended"
	}

	path := utils.GetPath("image", fmt.Sprintf("%v", event.Id), "images", subpath, fmt.Sprintf("%v", image.Id)) + ".jpg"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		WriteError(w, errors.New("image not found"), http.StatusNotFound, "Can't find file on filesystem")
		return
	}

	http.ServeFile(w, r, path)
}
