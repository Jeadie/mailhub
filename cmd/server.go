package cmd

import (
	"fmt"
	air "github.com/Jeadie/mailhub/pkg/airtable"
	"github.com/Jeadie/mailhub/pkg/db"
	"github.com/Jeadie/mailhub/pkg/server"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"os"
)

var (
	// Variables from flags
	environment    string
	serverAddr     string
	sendToAirtable bool

	// Variable for airtable integration. (Used only if sendToAirtable)
	airtableApiKey  string
	airtableBaseId  string
	airtableTableId string

	rootCmd = &cobra.Command{
		Use:   "mailhub",
		Short: "Mailhub stores mail for people",
		Long: `Mailhub stores mail, in a variety of forms (e.g. SMS), for each person.
		That is, it persists and retrieves mail on a per-person (or unique identifier) basis.
	`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if sendToAirtable {
				airtableApiKey = os.Getenv("AIRTABLE_API_KEY")
				airtableBaseId = os.Getenv("AIRTABLE_BASE_ID")
				airtableTableId = os.Getenv("AIRTABLE_TABLE_ID")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			isTest := environment == "test"
			r := gin.Default()

			eventChan := make(chan db.Event)
			var eventHandler []db.EventHandler

			if sendToAirtable {
				a := air.CreateAirtable(airtableBaseId, airtableTableId, airtableApiKey)
				eventHandler = append(eventHandler, func(e db.Event) error { return air.HandleDbEvent(a, e) })
			}

			if isTest {
				eventHandler = append(eventHandler, func(e db.Event) error { fmt.Println(e); return nil })
			}

			// Listen to events from the DB, send to each handler
			go func(dbEvents chan db.Event, handlers []db.EventHandler) {
				for e := range dbEvents {
					for _, hndl := range eventHandler {
						hndl(e)
					}
				}
			}(eventChan, eventHandler)

			server.ConstructEndpoints(r, db.CreateDao(isTest, eventChan))
			r.Run(serverAddr)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&serverAddr, "server-addr", ":8080", "TCP network address for the server to start on")
	rootCmd.PersistentFlags().StringVar(&environment, "env", "prod", "Environment stage of the Mailhub.")
	rootCmd.PersistentFlags().BoolVar(
		&sendToAirtable,
		"send-to-airtable",
		false,
		"Whether to send results to an Airtable integration. Requires ENV variables: AIRTABLE_API_KEY, AIRTABLE_BASE_ID, AIRTABLE_TABLE_ID",
	)

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
