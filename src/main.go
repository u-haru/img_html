package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/soniakeys/quant/median"
	"golang.org/x/image/draw"
)

func main() {
	var (
		Colors int
		Width  int
	)
	flag.IntVar(&Colors, "c", 0, "Color num")
	flag.IntVar(&Width, "w", 256, "Width(0 to raw)")
	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Printf("usage:\n %s [Option] from.png to.html\n Options\n", os.Args[0])
		flag.PrintDefaults()
		return
	}
	file, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if Width != 0 {
		img = resize(img, Width)
	}

	if Colors != 0 {
		q := median.Quantizer(Colors)
		img = q.Paletted(img).SubImage(img.Bounds())
	}

	if err != nil {
		fmt.Println("Error: Image could not be decoded")
		os.Exit(1)
	}
	img_to_html(img, flag.Arg(1))
}

func img_to_html(img image.Image, file string) {
	f, err := os.Create(file)
	if err != nil {
		fmt.Println("create:", err)
		return
	}
	defer f.Close()

	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	f.Write([]byte(`<html><body bgcolor="#ffffff"><pre><font face="monospace" style="height: 10px;line-height: 10px;font-size: 10px;">`))
	var precolor Pixel = Pixel{0, 0, 0, 1}
	var color Pixel
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			color = rgbaToPixel(r, g, b, 0)
			dot := fmt.Sprintf(`<font color="#%2x%2x%2x">`, r>>8, g>>8, b>>8)

			if precolor != color && precolor.A == 1 {
				f.Write([]byte(dot))
			} else if precolor != color {
				f.Write([]byte(`</font>` + dot))
			}
			f.Write([]byte("█"))
			precolor = color
		}
		f.Write([]byte("<br>"))
	}
	f.Write([]byte(`</font></font></pre></body></html>`))
}

func resize(img image.Image, w int) image.Image {
	rct := img.Bounds()
	m := float64(w) / float64(rct.Dx())
	h := int(float64(rct.Dy()) * m)
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, rct, draw.Over, nil)
	return dst
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{uint(r >> 8), uint(g >> 8), uint(b >> 8), 0}
}

type Pixel struct {
	R uint
	G uint
	B uint
	A uint
}
