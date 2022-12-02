package routes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/oxodao/photomaton/config"
	"github.com/oxodao/photomaton/orm"
	"github.com/oxodao/photomaton/services"
	"golang.org/x/exp/slices"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Register(r *mux.Router) {
	r.HandleFunc("/socket/{type}", socket)
	r.HandleFunc("/settings", settings)
	r.HandleFunc("/picture", picture).Methods(http.MethodPost)
}

func socket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	socketType := strings.ToUpper(vars["type"])
	if !slices.Contains(services.SOCKET_TYPES, socketType) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Println("Failed to parse hostport: ", err)
		fmt.Println("Got hostport: ", r.RemoteAddr)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if slices.Contains([]string{"[::1]", "127.0.0.1"}, host) {
		if !config.GET.DebugMode {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		fmt.Println("[Debug mode] Letting a remote connection from ", host)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade connection: ", err)
		return
	}

	services.GET.Join(socketType, conn)
}

func settings(w http.ResponseWriter, r *http.Request) {
	settings := services.GET.GetFrontendSettings()

	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(settings)
	w.Write(data)
}

func picture(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(256 * 1024) // Max picture size = 256mo, we should be good.
	if err != nil {
		fmt.Println("Unable to save picture: Parse form error => ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	event := r.FormValue("event")
	unattended := r.FormValue("unattended")
	image := r.FormValue("image")

	if len(event) == 0 || len(unattended) == 0 || len(image) == 0 {
		fmt.Println("Failed to save picture: bad request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	eventId, err := strconv.ParseInt(event, 10, 64)
	if err != nil {
		fmt.Println("Failed to get event id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isUnattended, err := strconv.ParseBool(unattended)
	if err != nil {
		fmt.Println("Failed to parse unattended var: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	evt, err := orm.GET.Events.GetEvent(eventId)
	if err != nil {
		fmt.Println("No event for the given id")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	imageName := ""

	img, err := orm.GET.Events.InsertImage(eventId, isUnattended)
	if err != nil {
		fmt.Println("Failed to insert image: ", err)
		fmt.Println("Defaulting name to current timestamp in the root folder for the event")
		imageName = fmt.Sprintf("%v.jpg", time.Now().Format("20060102-150405"))
	} else {
		imageName = fmt.Sprintf("%v.jpg", img.Id)
	}

	path, err := config.GET.GetImageFolder(evt.Id, isUnattended)
	if err != nil {
		fmt.Println("Failed to create path: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f, err := os.Create(filepath.Join(path, imageName))
	if err != nil {
		fmt.Println("Failed to create image file...")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	image = image[len("data:image/jpeg;base64,"):]
	data, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		fmt.Println("Failed to decode image, writing it to file as-is")
		_, err = f.Write([]byte(image))
		if err != nil {
			fmt.Println("Even failed to write the b64... sad")
		}
	} else {
		_, err = f.Write(data)
		if err != nil {
			fmt.Println("Failed to write the image to disk")
		}
	}

	if err = f.Sync(); err != nil {
		fmt.Println("Failed to sync the data ! be careful")
	}
}
