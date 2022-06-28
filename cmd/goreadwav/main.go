package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
)

type Chunk struct {
	ID   string
	Size uint32
	Data []byte
}

type RIFF struct {
	ID     string
	Size   uint32
	Data   []byte
	Format string
}

type Wav struct {
	WavFormat
	WavData
}

type WavFormat struct {
	ID         string
	Size       uint32
	AudioFmt   uint16
	ChNum      uint16
	SmplRate   uint32
	ByteRate   uint32
	BlockSize  uint16
	BitPerSmpl uint16
	// TODO: 拡張パラメータ
	// ExParaSize uint16
}

type WavData struct {
	ID   string
	Size uint32
	Data []byte
}

func main() {
	// TODO: ハードコードあとでやめる
	wavfile := "./wav/sample.wav"

	bytes, err := ioutil.ReadFile(wavfile)
	if err != nil {
		log.Panicln(err)
		os.Exit(1)
	}

	// RIFFチャンク
	wav, err := readWav(bytes)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Printf("WavFormat: %+v\n", wav.WavFormat)
	log.Printf("WavData ID: %v, Size: %v, Data: %v\n",
		wav.WavData.ID, wav.WavData.Size, wav.WavData.Data[:10])
}

func readChunk(b []byte) (*Chunk, error) {
	if len(b) < 8 {
		return nil, fmt.Errorf("invalid chunk. too short to contain id and size")
	}
	id := string(b[:4])
	size := binary.LittleEndian.Uint32(b[4:8])
	// TODO: データフィールドがvalidな基準 (ifの条件部分) 考える
	// dataはsize以上 (paddingがなければ一致、なければdataの方が1byte長い
	// dataの長さがuint32の最大値以下
	//
	if len(b[8:]) > math.MaxUint32 && uint32(len(b[8:])) < size {
		return nil, fmt.Errorf("invalid chunk. not enough data")
	}
	data := b[8 : 8+size]
	return &Chunk{
		ID:   id,
		Size: size,
		Data: data,
	}, nil
}

func readRIFF(b []byte) (*RIFF, error) {
	chunk, err := readChunk(b)
	if err != nil {
		return nil, err
	}
	if chunk.ID != "RIFF" {
		return nil, fmt.Errorf("not \"RIFF\", found %s", chunk.ID)
	}
	if len(chunk.Data) < 4 {
		return nil, fmt.Errorf("this RIFF has no format")
	}
	format := string(chunk.Data[:4])
	return &RIFF{
		ID:     chunk.ID,
		Size:   chunk.Size,
		Format: format,
		Data:   chunk.Data,
	}, nil
}

func readWav(b []byte) (*Wav, error) {
	riff, err := readRIFF(b)
	if err != nil {
		return nil, err
	}
	if riff.Format != "WAVE" {
		return nil, fmt.Errorf("not \"WAVE\", found %s", riff.Format)
	}
	// フォーマットチャンク
	wfmt, err := readWavFormat(riff.Data[4:])
	if err != nil {
		return nil, err
	}
	// データチャンク
	wdata, err := readWavData(riff.Data, 4+8+wfmt.Size)
	if err != nil {
		return nil, err
	}
	return &Wav{
		WavFormat: *wfmt,
		WavData:   *wdata,
	}, nil
}

func readWavFormat(b []byte) (*WavFormat, error) {
	chunk, err := readChunk(b)
	if err != nil {
		return nil, err
	}
	// TODO: バリデーション

	if chunk.ID != "fmt " {
		return nil, fmt.Errorf("not \"fmt \", found %s", chunk.ID)
	}

	id := chunk.ID
	size := chunk.Size

	audioFmt := binary.LittleEndian.Uint16(chunk.Data[:2])
	chNum := binary.LittleEndian.Uint16(chunk.Data[2:4])
	smplRate := binary.LittleEndian.Uint32(chunk.Data[4:8])
	byteRate := binary.LittleEndian.Uint32(chunk.Data[8:12])
	blockSize := binary.LittleEndian.Uint16(chunk.Data[12:14])
	bitPerSmpl := binary.LittleEndian.Uint16(chunk.Data[14:16])
	// TODO: 拡張パラメータ
	//if size > 16 {
	//exParaSize := binary.LittleEndian.Uint16(b[20:22])
	//}
	return &WavFormat{
		ID:         id,
		Size:       size,
		AudioFmt:   audioFmt,
		ChNum:      chNum,
		SmplRate:   smplRate,
		ByteRate:   byteRate,
		BlockSize:  blockSize,
		BitPerSmpl: bitPerSmpl,
	}, nil
}

func readWavData(b []byte, offset uint32) (*WavData, error) {
	chunk, err := readChunk(b[offset:])
	if err != nil {
		return nil, err
	}
	if chunk.ID != "data" {
		return nil, fmt.Errorf("not \"data\", found %s", chunk.ID)
	}
	return &WavData{
		ID:   chunk.ID,
		Size: chunk.Size,
		Data: chunk.Data,
	}, nil
}
