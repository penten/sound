package sound

import (
	"errors"
	"path"
)

type Sounder interface {
	SampleRate() int
	TotalSamples() int
	Channels() int
	Get(int, int) float64
	GetSlice(int, int, int) []float64
}

func Load(name string) (Sounder, error) {
	if path.Ext(name) == ".wav" {
		return LoadWav(name)
	}
	return nil, errors.New("Unsuported sound format")
}

// TODO: godoc compatible documentation
