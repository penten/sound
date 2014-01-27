package sound

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"os"
)

type Wav struct {
	chunkID   [4]byte
	chunkSize uint32
	format    [4]byte

	subchunk1ID   [4]byte
	subchunk1Size uint32
	audioFormat   uint16
	numChannels   uint16
	sampleRate    uint32
	byteRate      uint32
	blockAlign    uint16
	bitsPerSample uint16

	subchunk2ID [4]byte
	dataSize    uint32
	Data        []byte

	Samples [][]float64
}

func LoadWav(name string) (*Wav, error) {
	var w Wav
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	// Check file length
	var stat os.FileInfo
	stat, err = file.Stat()
	if err != nil || stat.Size() < 44 {
		return nil, errors.New("Not a valid wav file")
	}

	r := bufio.NewReader(file)

	// Read the header and verify the file is a WAV
	binary.Read(r, binary.BigEndian, &w.chunkID)
	binary.Read(r, binary.LittleEndian, &w.chunkSize)
	binary.Read(r, binary.BigEndian, &w.format)

	if string(w.format[:]) != "WAVE" {
		return nil, errors.New("Not a valid wav file")
	}

	// Read the fmt subchunk
	binary.Read(r, binary.BigEndian, &w.subchunk1ID)
	binary.Read(r, binary.LittleEndian, &w.subchunk1Size)
	binary.Read(r, binary.LittleEndian, &w.audioFormat)
	binary.Read(r, binary.LittleEndian, &w.numChannels)
	binary.Read(r, binary.LittleEndian, &w.sampleRate)
	binary.Read(r, binary.LittleEndian, &w.byteRate)
	binary.Read(r, binary.LittleEndian, &w.blockAlign)
	binary.Read(r, binary.LittleEndian, &w.bitsPerSample)

	if w.bitsPerSample > 16 || w.audioFormat != 1 {
		return nil, errors.New("Unsupported format")
	}

	// Read the data subchunk
	binary.Read(r, binary.BigEndian, &w.subchunk2ID)
	binary.Read(r, binary.LittleEndian, &w.dataSize)

	// Read the data
	w.Data = make([]byte, w.dataSize)
	binary.Read(r, binary.LittleEndian, &w.Data)

	file.Close()

	w.loadSamples()

	return &w, nil
}

func (w *Wav) SampleRate() int {
	return int(w.sampleRate)
}

func (w *Wav) TotalSamples() int {
	return int(w.dataSize) / (int(w.numChannels) * (int(w.bitsPerSample) / 8))
}

func (w *Wav) Channels() int {
	return int(w.numChannels)
}

func (w *Wav) loadSamples() {
	width := int(w.bitsPerSample) / 8
	c, s := w.Channels(), w.TotalSamples()

	// Create slices to hold the samples (channel x time)
	w.Samples = make([][]float64, c)
	for i := 0; i < c; i++ {
		w.Samples[i] = make([]float64, s)
	}

	// PCM wav can either be an 8 bit unsigned or 16 bit signed integer
	var min, max int
	if width == 1 {
		min, max = 0, 255
	} else {
		min, max = -32768, 32767
	}

	// Convert from raw bytes to amplitudes, between 0 and 1
	var amp, offset int 
	var b []byte
	for i := 0; i < s; i ++ {
		for j := 0; j < c; j++ {
			b = w.Data[offset : offset + width]
			if width == 1 {
				amp = int(b[0])
			} else {
				var s int16
				binary.Read(bytes.NewReader(b), binary.LittleEndian, &s)
				amp = int(s)
			}
			w.Samples[j][i] = float64(amp-min) / float64(max-min)
			offset += width
		}
	}
}

func (w *Wav) Get(c, t int) float64 {
	return w.Samples[c][t]
}

func (w *Wav) GetSlice(c, from, to int) []float64 {
	return w.Samples[c][from:to]	
}