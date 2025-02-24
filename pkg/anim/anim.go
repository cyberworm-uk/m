package anim

import (
	"log"
	"math"
	"math/cmplx"
	"runtime"

	"github.com/cyberworm-uk/m/pkg/m"
)

func ZoomTransform(mbrot *m.Mandelbrot, rP, iP, scale float64) {
	start, end := mbrot.Ranges()
	cscale := complex(scale, 0)
	p := complex(rP, iP)
	mbrot.SetRange(
		start+((p-start)*cscale),
		end-((end-p)*cscale),
	)
}

func JuliaTransform(mbrot *m.Mandelbrot, cur, max int, r float64) {
	c := cmplx.Rect(r, float64(cur*2)*math.Pi/float64(max))
	mbrot.F(
		func(z complex128, l int) int {
			var i int
			for i = 0; i <= l && cmplx.Abs(z) < 2; i++ {
				z = z*z + c
			}
			return i
		},
	)
}

func GenerateJuliaFrames(width, limit, fcount int, rStart, rEnd, iStart, iEnd, rad float64) chan *m.Raw {
	var rs = make(chan *m.Raw, runtime.NumCPU())
	go func() {
		var mbrot = m.NewM(width, limit, complex(rStart, iStart), complex(rEnd, iEnd))
		var logger = log.New(log.Writer(), "anim: ", log.Ldate|log.Ltime|log.Lshortfile)
		for i := 0; i < fcount; i++ {
			JuliaTransform(mbrot, i, fcount, rad)
			rs <- mbrot.Calculate()
			logger.Printf("INFO: %0.2f%% (%v/%v)\n", float64(i+1)/float64(fcount)*100, i+1, fcount)
		}
		close(rs)
	}()
	return rs
}

func GenerateZoomFrames(width, limit, fcount int, rStart, rEnd, iStart, iEnd, rP, iP, scale float64) chan *m.Raw {
	var rs = make(chan *m.Raw, runtime.NumCPU())
	go func() {
		var mbrot = m.NewM(width, limit, complex(rStart, iStart), complex(rEnd, iEnd))
		var logger = log.New(log.Writer(), "anim: ", log.Ldate|log.Ltime|log.Lshortfile)
		for i := 0; i < fcount; i++ {
			rs <- mbrot.Calculate()
			ZoomTransform(mbrot, rP, iP, scale)
			logger.Printf("INFO: %0.2f%% (%v/%v)\n", float64(i+1)/float64(fcount)*100, i+1, fcount)
		}
		close(rs)
	}()
	return rs
}
