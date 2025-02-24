package m

import (
	"bytes"
	"encoding/json"
	"log"
	"math"
	"math/cmplx"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type Mandelbrot struct {
	width, height, limit int
	start, end           complex128
	raw                  [][]int
	f                    func(c complex128, l int) int
	logger               *log.Logger
}

type Raw struct {
	Width  int     `json:"width"`
	Height int     `json:"height"`
	Limit  int     `json:"limit"`
	Data   [][]int `json:"data"`
}

func RawFromJson(jr []byte) *Raw {
	var r = &Raw{}
	var buf = bytes.NewBuffer(jr)
	var j = json.NewDecoder(buf)
	j.Decode(r)
	return r
}

func (r *Raw) Json() []byte {
	var buf = new(bytes.Buffer)
	var j = json.NewEncoder(buf)
	j.Encode(r)
	return buf.Bytes()
}

func f(c complex128, l int) int {
	var i int
	var z complex128 = 0 + 0i
	for i = 0; i < l && cmplx.Abs(z) < 2; i++ {
		z = z*z + c
	}
	return i
}

func NewM(width, limit int, a, b complex128) *Mandelbrot {
	var m = &Mandelbrot{
		width: width,
		limit: limit,
	}
	m.SetRange(a, b)
	m.F(f)
	m.logger = log.New(log.Writer(), "m: ", log.Ldate|log.Ltime|log.Lshortfile)
	return m
}

// F takes in a function f.
// the function f should take in a complex number and a limit l
// the function should return an integer in the range of (0,l]
func (m *Mandelbrot) F(f func(c complex128, l int) int) {
	m.f = f
}

func (m *Mandelbrot) Ranges() (complex128, complex128) {
	return m.start, m.end
}

func (m *Mandelbrot) SetRange(a, b complex128) {
	ra, rb := real(a), real(b)
	ia, ib := imag(a), imag(b)
	if ra > rb {
		rb, ra = ra, rb
	}
	if ia > ib {
		ib, ia = ia, ib
	}
	m.start = complex(ra, ia)
	m.end = complex(rb, ib)
	m.setHeight()
}

func (m *Mandelbrot) setHeight() {
	if real(m.start) != real(m.end) && imag(m.start) != imag(m.end) {
		m.height = int(math.Round(float64(m.width) / ((real(m.start) - real(m.end)) / (imag(m.start) - imag(m.end)))))
	} else {
		m.height = 0
	}
}

func (m *Mandelbrot) complex(x, y int) complex128 {
	xf := float64(x)
	wf := float64(m.width)
	yf := float64(y)
	hf := float64(m.height)
	return complex(
		real(m.start)+xf/wf*(real(m.end)-real(m.start)),
		imag(m.start)+yf/hf*(imag(m.end)-imag(m.start)),
	)
}

type job struct {
	x, y int
}

func (m *Mandelbrot) worker(in <-chan *job, wg *sync.WaitGroup, counter *atomic.Int64) {
	for j := range in {
		m.raw[j.x][j.y] = m.f(m.complex(j.x, j.y), m.limit)
		counter.Add(1)
	}
	wg.Done()
}

func (m *Mandelbrot) Calculate() *Raw {
	m.raw = make([][]int, m.width)
	for i := 0; i < m.width; i++ {
		m.raw[i] = make([]int, m.height)
	}
	var counter = &atomic.Int64{}
	var wg = &sync.WaitGroup{}
	var out = make(chan *job, runtime.NumCPU())
	var die = make(chan bool)
	defer close(die)
	t := time.NewTicker(5 * time.Second)
	go func() {
		total := m.width * m.height
		m.logger.Printf("INFO: Populating...(0/%d)\n", total)
		for {
			select {
			case <-t.C:
				curr := counter.Load()
				pct := float64(curr) * 100 / float64(total)
				m.logger.Printf("INFO: %0.2f%% (%d/%d)\n", pct, curr, total)
			case <-die:
				curr := counter.Load()
				m.logger.Printf("INFO: Done...(%d/%d)\n", curr, total)
				return
			}
		}
	}()
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go m.worker(out, wg, counter)
	}
	for x := 0; x < m.width; x++ {
		for y := 0; y < m.height; y++ {
			out <- &job{x: x, y: y}
		}
	}
	close(out)
	wg.Wait()
	return &Raw{
		Width:  m.width,
		Height: m.height,
		Limit:  m.limit,
		Data:   m.raw,
	}
}
