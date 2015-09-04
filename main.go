package main

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/robfig/graphics-go/graphics"
	"github.com/robfig/graphics-go/graphics/interp"
	"github.com/soniakeys/quant/median"
)

var (
	delay   = flag.Int("d", 10, "delay")
	output  = flag.String("o", "animated.gif", "output filename")
	reverse = flag.Bool("r", false, "inverse rotation")
	speed   = flag.Float64("x", 1.0, "rotating speed")
	zoom    = flag.Bool("z", false, "zoom animation") // experimental
	bg      = flag.String("bg", "FFFFFF", "background color")
)

func main() {
	flag.Parse()

	c, err := colorful.Hex(*bg)
	if err != nil {
		log.Fatal(err)
	}

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
	q := median.Quantizer(256)
	p := q.Quantize(make(color.Palette, 0, 256), src)
	for i := 0; i < limit; i++ {
		dst := image.NewPaletted(src.Bounds(), p)
		draw.Draw(dst, src.Bounds(), &image.Uniform{c}, image.ZP, draw.Src)
		err = graphics.Rotate(dst, src, &graphics.RotateOptions{Angle: base * float64(i)})
		if err != nil {
			log.Fatal(err)
		}
		if *zoom {
			w, h := float64(src.Bounds().Dx()), float64(src.Bounds().Dy())
			tmp := image.NewPaletted(src.Bounds(), p)
			draw.Draw(tmp, src.Bounds(), &image.Uniform{c}, image.ZP, draw.Src)
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
