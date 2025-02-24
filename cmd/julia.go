/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"math/cmplx"
	"os"

	"github.com/cyberworm-uk/m/pkg/m"
	"github.com/spf13/cobra"
)

// juliaCmd represents the julia command
var juliaCmd = &cobra.Command{
	Use:   "julia",
	Short: "Generates a raw output for the quadratic julia set",
	Long: `Generates a raw output for the quadratic julia set.
	Some interesting values for c are:
	--re-c=-.79 --im-c=-.15
	--re-c=-.162 --im-c=1.04
	--re-c=.3 --im-c=-0.1
	--re-c=-1.476 --im-c=0
	--re-c=-.12 --im-c=-.77
	--re-c=.28 --im-c=0.008`,
	Run: func(cmd *cobra.Command, args []string) {
		var width, limit int
		var rStart, rEnd, iStart, iEnd, rc, ic float64
		var e error
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
		if rc, e = cmd.Flags().GetFloat64("re-c"); e != nil {
			fmt.Fprintln(os.Stderr, e)
			return
		}
		if ic, e = cmd.Flags().GetFloat64("im-c"); e != nil {
			fmt.Fprintln(os.Stderr, e)
			return
		}
		var c = complex(rc, ic)
		mbrot := m.NewM(width, limit, complex(rStart, iStart), complex(rEnd, iEnd))
		mbrot.F(
			func(z complex128, l int) int {
				var i int
				for i = 0; i <= l && cmplx.Abs(z) < 2; i++ {
					z = z*z + c
				}
				return i
			},
		)
		fmt.Printf("%s\n", mbrot.Calculate().Json())
	},
}

func init() {
	rootCmd.AddCommand(juliaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// juliaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// juliaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	juliaCmd.Flags().Float64("re-start", -1.5, "real range start")
	juliaCmd.Flags().Float64("re-end", 1.5, "real range start")
	juliaCmd.Flags().Float64("re-c", -.8, "real part of c in z(n+1)=z(n)^2+c")
	juliaCmd.Flags().Float64("im-c", .156, "imaginary part of c in z(n+1)=z(n)^2+c")
}
