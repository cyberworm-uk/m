/*
Copyright © 2024 cyberworm <contact@cyberworm.uk>
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/cyberworm-uk/m/pkg/m"
	"github.com/cyberworm-uk/m/pkg/png"
	"github.com/mazznoer/colorgrad"
	"github.com/spf13/cobra"
)

// pngCmd represents the png command
var pngCmd = &cobra.Command{
	Use:   "png",
	Short: "Generate a PNG format mandelbrot image",
	Long: `Generate a PNG format mandelbrot image.
	Generated from a provided range of the complex plane.
	Either generated from attributes provided or a precalculcated .json via --from-json`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			width, limit               int
			rStart, rEnd, iStart, iEnd float64
			gradient                   []string
			fromJson                   string
			g                          colorgrad.Gradient
			e                          error
		)
		if fromJson, e = cmd.Flags().GetString("from-json"); e != nil {
			fmt.Fprintln(os.Stderr, e)
			return
		}
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
		if gradient, e = cmd.Flags().GetStringSlice("gradient"); e != nil {
			fmt.Fprintln(os.Stderr, e)
			return
		}
		if g, e = colorgrad.NewGradient().HtmlColors(gradient...).Build(); e != nil {
			fmt.Fprintln(os.Stderr, e)
			return
		}
		var img []byte
		var raw *m.Raw
		if len(fromJson) > 0 {
			var in io.ReadCloser
			if fromJson == "-" {
				in = os.Stdin
			} else {
				if in, e = os.Open(fromJson); e != nil {
					fmt.Fprintln(os.Stderr, e)
					return
				}
				defer in.Close()
			}
			var buf = new(bytes.Buffer)
			if _, e = buf.ReadFrom(in); e != nil {
				fmt.Fprintln(os.Stderr, e)
				return
			}
			raw = m.RawFromJson(buf.Bytes())
		} else {
			raw = m.NewM(width, limit, complex(rStart, iStart), complex(rEnd, iEnd)).Calculate()
		}
		img = png.RawToPng(raw, g)
		var fname = fmt.Sprintf("m-%v-%v-%v-%v.png", complex(rStart, iStart), complex(rEnd, iEnd), width, limit)
		var out io.WriteCloser
		if out, e = os.Create(fname); e != nil {
			fmt.Fprintln(os.Stderr, e)
			return
		}
		defer out.Close()
		if _, e = out.Write(img); e != nil {
			fmt.Fprintln(os.Stderr, e)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(pngCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pngCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pngCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	pngCmd.Flags().StringSlice("gradient", []string{"rgb(0,0,0 / 0%)", "rgb(0,0,0)", "rgb(255,255,255)"}, "list of colours to gradiate through")
	pngCmd.Flags().String("from-json", "", "json file to read raw data from (- for stdin)")
}
