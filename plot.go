package sound

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"
)

const windowSize = 0.01

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

	// find maximum value in DFT
	max := maxAmp(d)
	max = math.Log(1 + max)

	// plot each point on image
	for i, col := range d {
		for j, amp := range col {
			// Log(Xk+1) will give us a positive value
			// Using log here allows low amplitudes to be more visible
			bright := uint8(float64(255) * math.Log(1+amp) / max)
			image.Set(i, height-j, color.RGBA{bright, bright, bright, 255})
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

	// at a 40ms window length. At 44100sps, giving us 1764 samples per window
	// frequency range is from 2/N to 1 cycles per sample == 22050Hz to 25Hz
	// currently using 10ms window
	window := int(float64(snd.SampleRate()) * windowSize)
	length := (to - from) / window
	var dft [][]float64

	// create dft slice (time x frequency)
	dft = make([][]float64, length)
	for t := 0; t < length; t++ {
		dft[t] = dftWindow(snd.GetSlice(0, t*window, (t+1)*window))
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
			Xk += complex(xj[n], 0) * cmplx.Exp(complex(0, float64(-2.0*math.Pi*float64(k)*float64(n)/float64(N))))
		}
		Xj[k] = cmplx.Abs(Xk)
	}

	return Xj
}

func maxAmp(d [][]float64) float64 {
	max := 0.0
	for _, col := range d {
		for _, amp := range col {
			if amp > max {
				max = amp
			}
		}
	}
	return max
}

func dominantFrequency(Xj []float64) int {
	max := 0.0
	maxi := 0
	for k, amp := range Xj {
		if amp > max {
			max = amp
			maxi = k
		}
	}

	// X[k] corresponds to the amplitude of e^{i2\pi kn/N}
	// we have one oscillation every time the exponent goes through i2pi
	// so the frequency is k per window
	return maxi
}
