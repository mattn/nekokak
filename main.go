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
)

var (
	delay  = flag.Int("d", 10, "delay")
	output = flag.String("o", "animated.gif", "output filename")
)

func main() {
	flag.Parse()

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	src, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	g := &gif.GIF{
		Image:     []*image.Paletted{},
		Delay:     []int{},
		LoopCount: 0,
	}
	for i := 0; i < 360/10; i++ {
		dst := image.NewPaletted(src.Bounds(), palette.WebSafe)
		draw.Draw(dst, src.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
		err = graphics.Rotate(dst, src, &graphics.RotateOptions{Angle: math.Pi * 2 * float64(10*i) / 360})
		if err != nil {
			log.Fatal(err)
		}
		g.Image = append(g.Image, dst)
		g.Delay = append(g.Delay, *delay)
	}
	out, err := os.Create("animated.gif")
	if err != nil {
		log.Fatal(err)
	}
	err = gif.EncodeAll(out, g)
	if err != nil {
		log.Fatal(err)
	}
}
