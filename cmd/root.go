package cmd

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/oxodao/photobooth/config"
	"github.com/oxodao/photobooth/logs"
	"github.com/oxodao/photobooth/models"
	"github.com/oxodao/photobooth/orm"
	"github.com/oxodao/photobooth/routes"
	"github.com/oxodao/photobooth/services"
	"github.com/oxodao/photobooth/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "photobooth",
	Short: "The photobooth main app",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := orm.GET.AppState.GetState()
		if err != nil {
			if err != sql.ErrNoRows {
				logs.Error("Failed to load appstate: ", err)
				os.Exit(1)
			}

			as := models.AppState{
				HardwareID: uuid.New().String(),
				ApiToken:   nil, // @TODO: The token should be retreived from the API server while setting the photobooth up
			}
			err := orm.GET.AppState.CreateState(as)
			if err != nil {
				logs.Error("Failed to save the state: ", err)
				os.Exit(1)
			}

			logs.Info("Initializing the photobooth with id ", as.HardwareID)
		}

		err = orm.GET.Events.ClearExporting()
		if err != nil {
			logs.Error("Failed to clear event exporting.")
			logs.Error("Some event might be in a wrong state...")
			logs.Error(err)
		}

		r := mux.NewRouter()

		r.PathPrefix("/media/photobooth").Handler(http.StripPrefix("/media/photobooth", http.FileServer(http.Dir(utils.GetPath("images")))))

		routes.Register(r.PathPrefix("/api").Subrouter())
		if services.GET.AdminappFS != nil {
			r.PathPrefix("/admin").Handler(http.StripPrefix("/admin", http.FileServer(http.FS(*services.GET.AdminappFS))))
		} else {
			logs.Error("Failed to embed admin: not loaded")
		}

		if services.GET.WebappFS != nil {
			r.PathPrefix("/").Handler(http.FileServer(http.FS(*services.GET.WebappFS)))
		} else {
			logs.Error("Failed to embed webapp: not loaded")
		}

		logs.Infof("Photobooth app is listening on %v\n", config.GET.Web.ListeningAddr)
		err = http.ListenAndServe(config.GET.Web.ListeningAddr, r)
		if err != nil {
			logs.Error("Failed to listen on the given address/port", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
