package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var minimizeCmd = &cobra.Command{
	Use:   "minimize",
	Short: "Minimize the first matching window",
	Long:  `Minimize the first window that matches the given filters.`,
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
		minimizeWindow(w.handle)

		fmt.Printf("minimized: %s (%s, %d)\n", w.title, names[w.pid], w.pid)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(minimizeCmd)
}
