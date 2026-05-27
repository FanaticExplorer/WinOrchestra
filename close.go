package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var closeCmd = &cobra.Command{
	Use:   "close",
	Short: "Close the first matching window",
	Long:  `Close the first window that matches the given filters.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if flagTitle == "" && flagProcess == "" && flagPID == 0 {
			return fmt.Errorf("at least one filter is required: -t, -p, or --pid")
		}

		windows := enumerateWindows()
		names := processNames()

		matched := filterWindows(windows, flagTitle, flagProcess, flagPID, names)

		if len(matched) == 0 {
			return fmt.Errorf("no matching window found")
		}

		w := matched[0]
		closeWindow(w.handle)

		fmt.Printf("closed: %s (%s, %d)\n", w.title, names[w.pid], w.pid)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(closeCmd)
}
