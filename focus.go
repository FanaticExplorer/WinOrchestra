package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var focusCmd = &cobra.Command{
	Use:   "focus",
	Short: "Focus the first matching window",
	Long:  `Restore (if minimized) and bring the first matching window to the foreground.`,
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
		focusWindow(w.handle)

		fmt.Printf("focused: %s (%s, %d)\n", w.title, names[w.pid], w.pid)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(focusCmd)
}
