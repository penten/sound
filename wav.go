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

func (w *Wav) Get(t int) float32 {
	// Returns an amplitude between 0 and 1

	// TODO: error checking
	var total float32

	total = 0
	for c := 0; c < w.Channels(); c++ {
		total += w.GetChannel(t, c)
	}
	return total / float32(w.Channels())
}

func (w *Wav) GetChannel(t int, c int) float32 {
	// Returns an amplitude between 0 and 1

	// TODO: error checking
	// TODO: Move the amplitude decoding to Load?
	var amp, min, max int

	width := int(w.bitsPerSample) / 8
	offset := t * width * (w.Channels() + c)
	b := w.Data[offset : offset+width]

	// PCM wav can either be an 8 bit unsigned or 16 bit signed integer
	if width == 1 {
		amp = int(b[0])
		min, max = 0, 255
	} else {
		var s int16
		binary.Read(bytes.NewReader(b), binary.LittleEndian, &s)
		amp = int(s)
		min, max = -32768, 32767
	}

	return float32(amp-min) / float32(max-min)
}

