package sound

import (
	"path"
	"testing"
)

func TestWavBad(t *testing.T) {
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

func TestWavLoad(t *testing.T) {
	wav440 := path.Join("testdata", "440.wav")

	w, err := LoadWav(wav440)
	if err != nil {
		t.Error("Could not open file", err.Error())
		return
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

func TestWavCreate(t *testing.T) {
	w := CreateWav(100)

	if w.TotalSamples() != 100 {
		t.Error("Wrong number of samples", w.TotalSamples())
	}

	if w.Channels() != 2 {
		t.Error("Default number of channels should be 2")
	}
}

func TestWavSamples(t *testing.T) {
	wav440 := path.Join("testdata", "440.wav")

	w, err := LoadWav(wav440)
	if err != nil {
		t.Error("Could not open file", err.Error())
		return
	}

	if w.TotalSamples() != 441000 { // 10 seconds at 44100sps
		t.Error("Incorrect number of samples")
	}

	for i := 0; i < w.TotalSamples(); i += 10 {
		if w.Get(0, i) > 1.0 || w.Get(0, i) < 0.0 {
			t.Error("Sample out of bounds")
		}
	}
}

func TestDFT(t *testing.T) {
	wav440 := path.Join("testdata", "440.wav")

	w, err := LoadWav(wav440)
	if err != nil {
		t.Error("Could not open file", err.Error())
		return
	}

	d, err := dft(w, 0, w.TotalSamples()+1)
	if err == nil {
		t.Error("Should error when out of bounds")
	}

	d, err = dft(w, 0, w.TotalSamples())

	if len(d) != w.TotalSamples()/441 {
		t.Error("Time length of DFT incorrect:", len(d))
	}

	if len(d[0]) != 220 {
		t.Error("Frequency precision of DFT incorrect:", len(d[0]))
	}

	// not 440 because the window size chosen does not afford that much accuracy
	for i, _ := range d {
		if dominantFrequency(d[i]) != 400 {
			t.Error("Dominant frequency is incorrect:", dominantFrequency(d[i]))
			break
		}
	}
}
