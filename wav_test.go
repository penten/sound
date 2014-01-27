package sound

import (
	"path"
	"testing"
)

func TestBad(t *testing.T) {
	bad := path.Join("testdata", "bad.wav")
	notfound := path.Join("testdata", "notfound.wav")

	_, err := LoadWav(notfound)
	if err == nil {
		t.Error("File not found: Should return error")
	}

	_, err = LoadWav(bad)
	if err == nil {
		t.Error("Invalid wav file: Should return error")
	}
}

func TestLoad(t *testing.T) {
	wav440 := path.Join("testdata", "440.wav")

	w, err := LoadWav(wav440)
	if err != nil {
		t.Error("Could not open file")
	}

	if string(w.chunkID[:]) != "RIFF" {
		t.Error("Wrong Chunk ID")
	}

	if string(w.format[:]) != "WAVE" {
		t.Error("Wrong format")
	}

	if w.audioFormat != 1 || w.Channels() != 2 || w.SampleRate() != 44100 {
		t.Error("Incorrect metadata")
	}
}

func TestSamples(t *testing.T) {
	wav440 := path.Join("testdata", "440.wav")

	_, err := LoadWav(wav440)
	if err != nil {
		t.Error("Could not open file")
	}
}
