package sound

import (
	"errors"
	"path"
)

type Reader interface {
	SampleRate() int
	TotalSamples() int
	Channels() int
	Get(int, int) float64
	GetSlice(int, int, int) []float64
}

type Writer interface {
	Reader
	Save() error
	Set(int, int, float64)
}

func Load(name string) (Reader, error) {
	if path.Ext(name) == ".wav" {
		return LoadWav(name)
	}
	return nil, errors.New("Unsupported sound format")
}

func Generate(sw Writer, gen func(int) float64) {
	// Use the generator function to fill in the sound data at each sample

	sw.Save()
}

// TODO: godoc compatible documentation
