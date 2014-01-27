package sound

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func PlotOscillogram(snd Sounder, outpath string, from, to int) error {
	// make image
	height, width := 255, to-from
	red := color.RGBA{255, 0, 0, 255}
	image := image.NewGray(image.Rect(0, 0, width, height*2))

	// plot each point on image + a central red line at 0
	var amp float64
	for t := from; t < to; t++ {
		amp = float64(height) * (snd.Get(0, t) * 2)
		image.Set(t, int(amp), color.White)

		image.Set(t, height, red)
	}

	of, err := os.Create(outpath)
	if err != nil {
		return err
	}

	return png.Encode(of, image)
}
