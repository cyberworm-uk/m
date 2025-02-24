# m
mandlebrot image generator in go

```
A set of commands to generate mandelbrot fractal graphics

Usage:
  m [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  gif         A brief description of your command
  help        Help about any command
  julia       Generates a raw output for the quadratic julia set
  mp4         Generate a H264 format mandelbrot animation
  png         Generate a PNG format mandelbrot image
  raw         Generates a raw output for the mandelbrot set

Flags:
      --config string    config file (default is $HOME/.m.yaml)
  -h, --help             help for m
      --im-end float     imaginary range end (default 1.2)
      --im-start float   imaginary range start (default -1.2)
      --limit int        iteration test limit (default 200)
      --re-end float     real range end (default 1)
      --re-start float   real range start (default -2)
      --width int        output width (default 1024)

Use "m [command] --help" for more information about a command.
```

# png
```
$ go install github.com/cyberworm-uk/m@latest
$ ~/go/bin/m png --help
Generate a PNG format mandelbrot image.
        Generated from a provided range of the complex plane.
        Either generated from attributes provided or a precalculcated .json via --from-json

Usage:
  m png [flags]

Flags:
      --from-json string   json file to read raw data from (- for stdin)
      --gradient strings   list of colours to gradiate through (default ["rgb(0,0,0 / 0%)","rgb(0,0,0)","rgb(255,255,255)"])
  -h, --help               help for png

Global Flags:
      --config string    config file (default is $HOME/.m.yaml)
      --im-end float     imaginary range end (default 1.2)
      --im-start float   imaginary range start (default -1.2)
      --limit int        iteration test limit (default 200)
      --re-end float     real range end (default 1)
      --re-start float   real range start (default -2)
      --width int        output width (default 1024)
```
![mandlebrot png](m-(-2-1.2i)-(1+1.2i)-1024-200.png)

# gif
```
$ go install github.com/cyberworm-uk/m/cmd/gif@latest
Usage of /home/user/go/bin/gif:
  -anim string
        type of animation: "resolve" or "zoom" (default "resolve")
  -frames int
        number of frames of animation (default 100)
  -imag-end float
        end of imaginary range (default 1.2)
  -imag-start float
        start of the imaginary range (default -1.2)
  -limit int
        iteration limit for bounding check (precision) (default 200)
  -real-end float
        end of the real range (default 1)
  -real-start float
        start of the real range (default -2)
  -width int
        image output width (default 1024)
```
![mandlebrot gif](m-(-2-1.2i)-(1+1.2i)-1024-1000.gif)

# mp4

![julia mp4](m-(-2-1.2i)-(2+1.2i)-800-200.mp4)