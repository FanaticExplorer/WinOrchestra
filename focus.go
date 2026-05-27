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
		w, names, err := findFirstWindow()
		if err != nil {
			return err
		}

		focusWindow(w.handle)

		fmt.Printf("focused: %s (%s, %d)\n", w.title, names[w.pid], w.pid)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(focusCmd)
}
