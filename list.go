package main

import (
	"encoding/json"
	"fmt"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	listRaw bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List matching windows as JSON",
	Long:  `List all visible windows as JSON. Optionally filter by title, process name, or PID.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		windows := enumerateWindows()
		names := processNames()

		matched := filterWindows(windows, flagTitle, flagProcess, flagClass, flagPID, names)

		foregroundHwnd, _, _ := procGetForegroundWindow.Call()
		entries := make([]windowEntry, 0, len(matched))
		for _, w := range matched {
			entries = append(entries, toEntry(w, names, syscall.Handle(foregroundHwnd)))
		}

		var out []byte
		var err error
		if listRaw {
			out, err = json.Marshal(entries)
		} else {
			out, err = json.MarshalIndent(entries, "", "  ")
		}
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	},
}

func init() {
	listCmd.Flags().BoolVar(&listRaw, "raw", false, "Output JSON without indentation")

	rootCmd.AddCommand(listCmd)
}
