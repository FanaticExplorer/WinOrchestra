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
		windows := enumerateWindows()
		names := processNames()

		matched := filterWindows(windows, flagTitle, flagProcess, flagPID, names)

		if len(matched) == 0 {
			return fmt.Errorf("no matching window found")
		}

		closeWindow(matched[0].handle)

		fmt.Println("Closed:", matched[0].title)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(closeCmd)
}
