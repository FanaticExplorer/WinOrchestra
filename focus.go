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
		windows, names, err := findWindows()
		if err != nil {
			return err
		}

		for _, w := range windows {
			focusWindow(w.handle)
			fmt.Printf("focused: %s (%s, %d)\n", w.title, names[w.pid], w.pid)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(focusCmd)
}
