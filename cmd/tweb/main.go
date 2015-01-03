package main

import (
	"log"

	"github.com/spf13/cobra"

	"try/thunder"
)

func main() {
	var (
		DEBUG  bool
		local  bool
		listen string
	)

	cmd := &cobra.Command{
		Use:   "tweb",
		Short: "Starts the HTTP interface for a connected Thunder Launcher",
		Long:  "Starts the HTTP interface for a connected Thunder Laucher",
		Run: func(cmd *cobra.Command, args []string) {
			var launcher thunder.Launcher

			if DEBUG {
				launcher = thunder.NewMockLauncher()
			} else {
				tl, err := thunder.GetConnectedThunderLaunchers()
				if err != nil {
					log.Fatal("There was an error looking for connected launchers: ", err)
				}
				if len(tl) < 1 {
					log.Fatal("No connected launchers found. Connect one and try again.")
				}
				// for now just use the first launcher
				launcher = tl[0]
				for i := 1; i < len(tl); i++ {
					tl[i].Close()
				}
			}

			log.Fatal(Listen(listen, launcher, local))
		},
	}

	cmd.Flags().BoolVarP(&DEBUG, "debug", "d", false, "Debug mode doesn't require a connected Thunder Launcher")
	cmd.Flags().BoolVarP(&local, "local", "l", false, "Use local static assets instead of those packaged with the binary")
	cmd.Flags().StringVarP(&listen, "http", "b", ":8080", "The address and port to bind the HTTP interface to")

	cmd.Execute()
}
