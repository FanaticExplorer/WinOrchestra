package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	flagTitle   string
	flagProcess string
	flagPID     int
	flagClass   string
	flagAll     bool
)

var rootCmd = &cobra.Command{
	Use:   "winorchestra",
	Short: "Control windows by title, process name, or PID",
	Long:  "Control windows by title, process name, or PID.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`Run "winorchestra --help" for usage.`)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flagTitle, "title", "t", "", "Filter by window title (partial, case-insensitive)")
	rootCmd.PersistentFlags().StringVarP(&flagProcess, "process", "p", "", "Filter by process .exe name (partial, case-insensitive)")
	rootCmd.PersistentFlags().IntVar(&flagPID, "pid", 0, "Filter by exact process ID")
	rootCmd.PersistentFlags().StringVar(&flagClass, "class", "", "Filter by window class (partial, case-insensitive)")
	rootCmd.PersistentFlags().BoolVar(&flagAll, "all", false, "Apply action to all matching windows, not just the first")
}
