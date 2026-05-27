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
		windows := enumerateWindows()
		names := processNames()

		matched := filterWindows(windows, flagTitle, flagProcess, flagPID, names)

		if len(matched) == 0 {
			return fmt.Errorf("no matching window found")
		}

		minimizeWindow(matched[0].handle)

		fmt.Println("Minimized:", matched[0].title)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(minimizeCmd)
}
