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
		windows, names, err := findWindows()
		if err != nil {
			return err
		}

		for _, w := range windows {
			minimizeWindow(w.handle)
			fmt.Printf("minimized: %s (%s, %d)\n", w.title, names[w.pid], w.pid)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(minimizeCmd)
}
