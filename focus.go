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
		windows := enumerateWindows()
		names := processNames()

		matched := filterWindows(windows, flagTitle, flagProcess, flagPID, names)

		if len(matched) == 0 {
			return fmt.Errorf("no matching window found")
		}

		focusWindow(matched[0].handle)

		fmt.Println("Focused:", matched[0].title)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(focusCmd)
}
