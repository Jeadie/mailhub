package cmd

import (
	"fmt"
	"github.com/Jeadie/mailhub/pkg/db"
	"github.com/Jeadie/mailhub/pkg/server"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"os"
)

var (
	// Variables from flags
	environment string
	serverAddr  string

	rootCmd = &cobra.Command{
		Use:   "mailhub",
		Short: "Mailhub stores mail for people",
		Long: `Mailhub stores mail, in a variety of forms (e.g. SMS), for each person.
		That is, it persists and retrieves mail on a per-person (or unique identifier) basis.
	`,
		Run: func(cmd *cobra.Command, args []string) {
			isTest := environment == "test"
			r := gin.Default()

			eventChan := make(chan db.Event)
			if isTest {
				go func(dbEvents chan db.Event) {
					for e := range dbEvents {
						fmt.Println(e)
					}
				}(eventChan)
			}

			server.ConstructEndpoints(r, db.CreateDao(isTest, eventChan))
			r.Run(serverAddr)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&serverAddr, "server-addr", ":8080", "TCP network address for the server to start on")
	rootCmd.PersistentFlags().StringVar(&environment, "env", "prod", "Environment stage of the Mailhub.")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}