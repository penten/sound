package sound

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"errors"
	"math"
	"math/cmplx"
)

func PlotOscillogram(snd Sounder, outpath string, from, to int) error {
	// make image
	height, width := 255, to-from
	gray := color.RGBA{255, 0, 0, 255}
	image := image.NewGray(image.Rect(0, 0, width, height*2))

	// plot each point on image + a central gray line at 0
	var amp float64
	for t := from; t < to; t++ {
		amp = float64(height) * (snd.Get(0, t) * 2)
		image.Set(t, int(amp), color.White)

		image.Set(t, height, gray)
	}

	of, err := os.Create(outpath)
	if err != nil {
		return err
	}

	return png.Encode(of, image)
}

func PlotSpectogram(snd Sounder, outpath string, from, to int) error {
	d, err := dft(snd, from, to)
	if err != nil {
		return err
	}

	// make image
	height, width := len(d[0]), len(d)
	image := image.NewGray(image.Rect(0, 0, width, height))

	// plot each point on image
	for i, col := range d {
		for j, amp := range col {
			// TODO: is there a color.Gray?
			// TODO: currently we can only see a few frequencies with the highest
			// amplitudes. Using color will help. (use log amplitudes otherwise)
			brightness := uint8(float64(255) * amp)
			image.Set(i, j, color.RGBA{brightness, brightness, brightness, 255})
		}
	}

	of, err := os.Create(outpath)
	if err != nil {
		return err
	}

	return png.Encode(of, image)
}

func dft(snd Sounder, from, to int) ([][]float64, error) {
	if from < 0 || to > snd.TotalSamples() {
		return nil, errors.New("Out of bounds")
	}

	// at a 40ms window length. At 44100sps, this gives us 1764 samples per window
	// frequency range is from 2/N to 1 cycles per sample == 22050Hz to 25Hz
	// currently using 10ms window
	window := int(float64(snd.SampleRate()) * 0.01)
	length := (to-from) / window
	var dft [][]float64

	// create dft slice (time x frequency)
	dft = make([][]float64, length)
	for t := 0; t < length; t++ {
		dft[t] = dftWindow(snd.GetSlice(0, t*window, (t+1)*window))
	}

	// tmp: find maximum and reduce all to 0.0-1.0 range
	// awful code to be moved to dftWindow
	max := 0.0
	for _, col := range dft {
		for _, amp := range col {
			if amp > max {
				max = amp
			}
		}
	}

	for t, col := range dft {
		for freq, amp := range col {
			dft[t][freq] = amp/max
		}
	}

	return dft, nil
}

func dftWindow(xj []float64) []float64 {
	// X_k = \sum_{n=0}^{N-1}x_n e^{-i 2 \pi kn/N}
	// Up to N/2 since second half is mirrored
	N := len(xj)
	Xj := make([]float64, N/2)

	// start at one to ignore the first "flat line" frequency
	for k := 1; k < N/2; k++ {
		Xk := complex(0, 0)
		for n := 0; n < N; n++ {
			Xk += complex(xj[n], 0) * cmplx.Exp(complex(0, float64(-2.0 * math.Pi * float64(k) * float64(n) / float64(N))))
		}
		Xj[k] = cmplx.Abs(Xk)
	}

	return Xj
}