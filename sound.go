package sound

import (
	"errors"
	"path"
)

type Sounder interface {
	SampleRate() int
	TotalSamples() int
	Channels() int
	Get(int) float32
	GetChannel(int, int) float32
}

func Load(name string) (Sounder, error) {
	if path.Ext(name) == ".wav" {	
		return LoadWav(name)
	}
	return nil, errors.New("Unsuported sound format")
}