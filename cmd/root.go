package cmd

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/oxodao/photomaton/config"
	"github.com/oxodao/photomaton/models"
	"github.com/oxodao/photomaton/orm"
	"github.com/oxodao/photomaton/routes"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "photomaton",
	Short: "The photomaton main app",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := orm.GET.AppState.GetState()
		if err != nil {
			if err != sql.ErrNoRows {
				fmt.Println("Failed to load appstate: ", err)
				os.Exit(1)
			}

			as := models.AppState{
				HardwareID: uuid.New().String(),
				ApiToken:   nil, // @TODO: The token should be retreived from the API server while setting the photomaton up
			}
			err := orm.GET.AppState.CreateState(as)
			if err != nil {
				fmt.Println("Failed to save the state: ", err)
				os.Exit(1)
			}

			fmt.Println("Initializing the photobooth with id ", as.HardwareID)
		}

		r := mux.NewRouter()

		routes.Register(r.PathPrefix("/api").Subrouter())

		fmt.Printf("Photobooth app is listening on %v\n", config.GET.Web.ListeningAddr)
		err = http.ListenAndServe(config.GET.Web.ListeningAddr, r)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {

}
