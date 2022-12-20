package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/oxodao/photobooth/models"
	"github.com/oxodao/photobooth/orm"
	"github.com/oxodao/photobooth/services"
	"github.com/oxodao/photobooth/utils"
)

func registerAdminRoutes(r *mux.Router) {
	r.Use(AuthenticatedMiddleware)

	r.HandleFunc("/login", isPasswordValid)
	r.HandleFunc("/set_mode/{mode}", setMode).Methods(http.MethodPost)

	r.HandleFunc("/exports/{event_id}", getLastExports).Methods(http.MethodGet)
	r.HandleFunc("/exports/{export_id}/download", downloadExport).Methods(http.MethodGet)

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

func getLastExports(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	event_id := values["event_id"]

	eventId, err := strconv.ParseInt(event_id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event, err := orm.GET.Events.GetEvent(eventId)
	if err != nil {
		// @Todo: handle better
		w.WriteHeader(http.StatusNotFound)
		return
	}

	exportedEvents, err := orm.GET.Events.GetExportedEvents(event, 5)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	data, _ := json.MarshalIndent(exportedEvents, "", "  ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func downloadExport(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	exportedEventStr := values["export_id"]

	exportedEventId, err := strconv.ParseInt(exportedEventStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exportedEvent, err := orm.GET.Events.GetExportedEvent(exportedEventId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)

		return
	}

	path := utils.GetPath("images", fmt.Sprintf("%v", exportedEvent.EventId), "exports", exportedEvent.Filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		WriteError(w, errors.New("export not found"), http.StatusNotFound, "Can't find file on filesystem")
		return
	}

	http.ServeFile(w, r, path)
}
