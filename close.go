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
		w, names, err := findFirstWindow()
		if err != nil {
			return err
		}

		closeWindow(w.handle)

		fmt.Printf("closed: %s (%s, %d)\n", w.title, names[w.pid], w.pid)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(closeCmd)
}
