/*
Copyright © 2024 cyberworm <contact@cyberworm.uk>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/cyberworm-uk/m/pkg/anim"
	"github.com/cyberworm-uk/m/pkg/m"
	"github.com/cyberworm-uk/m/pkg/mp4"
	"github.com/mazznoer/colorgrad"
	"github.com/spf13/cobra"
)

// mp4Cmd represents the mp4 command
var mp4Cmd = &cobra.Command{
	Use:   "mp4",
	Short: "Generate a H264 format mandelbrot animation",
	Long: `Generate a H264 format mandelbrot animation.
	Generated from a provided range of the complex plane.
	Either generated from attributes provided or a precalculcated .json via --from-json`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			width, limit, frames                            int
			rStart, rEnd, iStart, iEnd, rP, iP, rRad, scale float64
			gradient                                        []string
			grad                                            colorgrad.Gradient
			animType                                        string
			e                                               error
		)
		if frames, e = cmd.Flags().GetInt("frames"); e != nil {
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
		if rP, e = cmd.Flags().GetFloat64("re-point"); e != nil {
			fmt.Fprintln(os.Stderr, e)
			return
		}
		if iP, e = cmd.Flags().GetFloat64("im-point"); e != nil {
			fmt.Fprintln(os.Stderr, e)
			return
		}
		if scale, e = cmd.Flags().GetFloat64("scale"); e != nil {
			fmt.Fprintln(os.Stderr, e)
			return
		}
		if rRad, e = cmd.Flags().GetFloat64("radius"); e != nil {
			fmt.Fprintln(os.Stderr, e)
			return
		}
		if gradient, e = cmd.Flags().GetStringSlice("gradient"); e != nil {
			fmt.Fprintln(os.Stderr, e)
			return
		}
		if grad, e = colorgrad.NewGradient().HtmlColors(gradient...).Build(); e != nil {
			fmt.Fprintln(os.Stderr, e)
			return
		}
		if animType, e = cmd.Flags().GetString("type"); e != nil {
			fmt.Fprintln(os.Stderr, e)
			return
		}
		var raws chan *m.Raw
		switch animType {
		case "julia":
			raws = anim.GenerateJuliaFrames(width, limit, frames, rStart, rEnd, iStart, iEnd, rRad)
		case "zoom":
			raws = anim.GenerateZoomFrames(width, limit, frames, rStart, rEnd, iStart, iEnd, rP, iP, scale)
		}
		var fname = fmt.Sprintf("m-%v-%v-%v-%v.264", complex(rStart, iStart), complex(rEnd, iEnd), width, limit)
		var out io.WriteCloser
		if out, e = os.Create(fname); e != nil {
			fmt.Fprintln(os.Stderr, e)
			return
		}
		defer out.Close()
		fmt.Fprintf(os.Stderr, "Saving to %q\n", fname)
		mp4.RawsToMp4(raws, grad, out)
	},
}

func init() {
	rootCmd.AddCommand(mp4Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mp4Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mp4Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	mp4Cmd.Flags().Float64("re-point", -0.13856524454488, "real point to zoom to")
	mp4Cmd.Flags().Float64("im-point", -0.64935990748190, "imaginary point to zoom to")
	mp4Cmd.Flags().Float64("scale", 0.01, "amount to zoom per frame")
	mp4Cmd.Flags().StringSlice("gradient", []string{"rgb(0,0,0 / 0%)", "rgb(0,0,0)", "rgb(255,255,255)"}, "list of colours to gradiate through")
	mp4Cmd.Flags().String("from-json", "", "json file to read raw data from (- for stdin)")
	mp4Cmd.Flags().String("type", "zoom", "animation type to generate")
	mp4Cmd.Flags().Float64("radius", 0.7636753236814714, "radius to rotate around the origin of the complex plane for julia")
	mp4Cmd.Flags().Int("frames", 100, "number of frames in gif")
}
