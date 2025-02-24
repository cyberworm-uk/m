/*
Copyright © 2024 cyberworm <contact@cyberworm.uk>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/cyberworm-uk/m/pkg/m"
	"github.com/spf13/cobra"
)

// rawCmd represents the raw command
var rawCmd = &cobra.Command{
	Use:   "raw",
	Short: "Generates a raw output for the mandelbrot set",
	Long:  `Generates a raw output for the mandelbrot set.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			width, limit               int
			rStart, rEnd, iStart, iEnd float64
			e                          error
		)
		if width, e = cmd.Flags().GetInt("width"); e != nil {
			fmt.Fprintln(os.Stderr, e)
			return
		}
		if limit, e = cmd.Flags().GetInt("limit"); e != nil {
			fmt.Fprintln(os.Stderr, e)
			return
		}
		if rStart, e = cmd.Flags().GetFloat64("re-start"); e != nil {
			fmt.Fprintln(os.Stderr, e)
			return
		}
		if rEnd, e = cmd.Flags().GetFloat64("re-end"); e != nil {
			fmt.Fprintln(os.Stderr, e)
			return
		}
		if iStart, e = cmd.Flags().GetFloat64("im-start"); e != nil {
			fmt.Fprintln(os.Stderr, e)
			return
		}
		if iEnd, e = cmd.Flags().GetFloat64("im-end"); e != nil {
			fmt.Fprintln(os.Stderr, e)
			return
		}
		mbrot := m.NewM(width, limit, complex(rStart, iStart), complex(rEnd, iEnd))
		fmt.Printf("%s\n", mbrot.Calculate().Json())
	},
}

func init() {
	rootCmd.AddCommand(rawCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rawCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rawCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
