package sound

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func PlotOscillogram(name, outpath string, from, to int) error {
	snd, err := Load(name)
	if err != nil {
		return err
	}

	// make image
	height, width := 255, to - from
	red := color.RGBA{255, 0, 0, 255}
	image := image.NewGray(image.Rect(0, 0, width, height*2))

	// plot each point on image + a central red line at 0
	var amp float32
	for t := from; t < to; t++ {
		amp = float32(height) * (snd.GetChannel(t, 0) * 2)
		image.Set(t, int(amp), color.White)

		image.Set(t, height, red)
	}

	var of *os.File
	of, err = os.Create(outpath)
	if err != nil {
		return err
	}

	return png.Encode(of, image)
}
