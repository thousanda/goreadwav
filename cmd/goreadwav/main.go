package main

import (
	"encoding/binary"
	"io/ioutil"
	"log"
	"os"
)

type Chunk struct {
	ID   string
	Size uint
	Data []byte
}

func main() {
	// TODO: ハードコードあとでやめる
	wav := "./wav/sample2.wav"

	bytes, err := ioutil.ReadFile(wav)
	if err != nil {
		log.Panicln(err)
		os.Exit(1)
	}

	// RIFFチャンク
	id := string(bytes[:4])
	size := binary.LittleEndian.Uint32(bytes[4:8])
	data := string(bytes[8:12])
	log.Println(id, size, data)

	// fmtチャンク
	id = string(bytes[12:16])
	size = binary.LittleEndian.Uint32(bytes[16:20])
	audioFmt := binary.LittleEndian.Uint16(bytes[20:22])
	chNum := binary.LittleEndian.Uint16(bytes[22:24])
	byteRate := binary.LittleEndian.Uint32(bytes[24:28])
	blockSize := binary.LittleEndian.Uint16(bytes[28:30])
	bitPerSmpl := binary.LittleEndian.Uint16(bytes[32:34])
	exParaSize := binary.LittleEndian.Uint16(bytes[34:36])
	log.Printf("{id: %v, size: %v, audioFmt: %v, chNum: %v, byteRate: %v, blockSize: %v, bitPerSmpl: %v, exParaSize: %v}\n",
		id, size, audioFmt, chNum, byteRate, blockSize, bitPerSmpl, exParaSize)
}
