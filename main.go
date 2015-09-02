package main

import (
	"flag"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"

	"github.com/robfig/graphics-go/graphics"
	"github.com/robfig/graphics-go/graphics/interp"
)

var (
	delay   = flag.Int("d", 10, "delay")
	output  = flag.String("o", "animated.gif", "output filename")
	reverse = flag.Bool("r", false, "inverse rotation")
	speed   = flag.Float64("x", 1.0, "speed")
	zoom    = flag.Bool("z", false, "zoom") // experimental
)

func main() {
	flag.Parse()

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	src, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	g := &gif.GIF{
		Image:     []*image.Paletted{},
		Delay:     []int{},
		LoopCount: 0,
	}

	var base float64 = (math.Pi * 2) * 10 * (*speed) / 360
	if *reverse {
		base *= -1
	}

	limit := int(360 / 10 / (*speed))
	for i := 0; i < limit; i++ {
		dst := image.NewPaletted(src.Bounds(), palette.WebSafe)
		draw.Draw(dst, src.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
		err = graphics.Rotate(dst, src, &graphics.RotateOptions{Angle: base * float64(i)})
		if err != nil {
			log.Fatal(err)
		}
		if *zoom {
			w, h := float64(src.Bounds().Dx()), float64(src.Bounds().Dy())
			tmp := image.NewPaletted(src.Bounds(), palette.WebSafe)
			draw.Draw(tmp, src.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
			z := float64(0.5 + float64(i)/30.0)
			graphics.I.
				Scale(z, z).
				Translate((w-w*z)/2, (h-h*z)/2).
				Transform(tmp, dst, interp.Bilinear)
			dst = tmp
		}
		g.Image = append(g.Image, dst)
		g.Delay = append(g.Delay, *delay)
	}
	out, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	err = gif.EncodeAll(out, g)
	if err != nil {
		log.Fatal(err)
	}
}
