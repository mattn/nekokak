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
	reverse = flag.Bool("r", false, "inverse rotation")
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

	var base float64 = math.Pi * 2
	if *reverse {
		base *= -1
	}

	for i := 0; i < 360/10; i++ {
		dst := image.NewPaletted(src.Bounds(), palette.WebSafe)
		draw.Draw(dst, src.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
		err = graphics.Rotate(dst, src, &graphics.RotateOptions{Angle: base * float64(10*i) / 360})
		if err != nil {
			log.Fatal(err)
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
