// https://sdkdocs.roku.com/display/sdkdoc/Trick+Mode+Support
package main

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const file = "gladiator.bif"
const version = "0.1.0"
const outputDir = "./test"
const output = "./test/test"

type BIF struct {
	FileType            string
	Version             int
	FrameCount          int
	FramewiseSeparation int
	Frames              []Frame
}

type Frame struct {
	Timestamp uint32
	Offset    uint32
}

func main() {

	fmt.Printf("bifparser - %s\n", version)

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, os.ModePerm)
	}

	// Type
	bif := isBIF(f)
	fmt.Printf("is BIF?: %t\n", bif)

	// Version
	v := getVersion(f)
	fmt.Printf("Version: %d\n", v)

	// Frame Count
	fc := getFramesCount(f)
	fmt.Printf("Number of frames: %d\n", fc)

	// Framewise Separation
	fs := getFramewiseSeparation(f)
	fmt.Printf("Framewise Separation: %d ms\n", fs)

	// Get frames
	var byteIndex int64 = 64 // BIF index starts at byte 64.
	var frames []Frame

	for i := 0; i < fc; i++ {
		ts, offset := readFrame(f, byteIndex)
		frame := Frame{
			Timestamp: ts * fs,
			Offset:    offset,
		}
		frames = append(frames, frame)
		byteIndex += 8 // Next frame every 8 bytes.
	}

	for k, v := range frames {
		var nextOffset uint32
		if k == len(frames)-1 {
			// nextOffset = frames[len(frames)-1].Offset
			fmt.Println("end")
		} else {
			nextOffset = frames[k+1].Offset

			// Calculate frame image size from next offset.
			frameLen := nextOffset - v.Offset

			// Create image.
			createFrameImage(f, k, int64(v.Offset), int(frameLen))
		}
	}
	f.Close()
}

func isBIF(f *os.File) bool {
	b := make([]byte, 8)
	_, err := f.Read(b)
	if err != nil {
		panic(err)
	}
	magic := string(b)
	return strings.Contains(magic, "BIF")
}

func getVersion(f *os.File) uint32 {
	f.Seek(8, 0)
	b := make([]byte, 4)
	_, err := f.Read(b)
	if err != nil {
		panic(err)
	}
	version := binary.LittleEndian.Uint32(b)
	return version
}

func getFramesCount(f *os.File) int {
	f.Seek(12, 0)
	b := make([]byte, 4)
	_, err := f.Read(b)
	if err != nil {
		panic(err)
	}
	numFrames := binary.LittleEndian.Uint32(b)
	return int(numFrames)
}

func getFramewiseSeparation(f *os.File) uint32 {
	f.Seek(16, 0)
	b := make([]byte, 4)
	_, err := f.Read(b)
	if err != nil {
		panic(err)
	}
	framewiseSeparation := binary.LittleEndian.Uint32(b)
	return framewiseSeparation
}

func readFrame(f *os.File, offset int64) (uint32, uint32) {
	f.Seek(offset, 0)
	b := make([]byte, 4)
	f.Read(b)
	frameTimestamp := binary.LittleEndian.Uint32(b)

	f.Seek(offset+4, 0)
	b2 := make([]byte, 4)
	f.Read(b2)
	frameOffset := binary.LittleEndian.Uint32(b2)
	return frameTimestamp, frameOffset
}

func getFrameImage(f *os.File, offset int64, len int) string {
	f.Seek(offset, 0)
	b := make([]byte, len)
	f.Read(b)
	enc := base64.StdEncoding.EncodeToString(b)
	return enc
}

func createFrameImage(f *os.File, i int, offset int64, len int) {
	f.Seek(offset, 0)
	b := make([]byte, len)
	f.Read(b)

	filename := fmt.Sprintf("%s_%s.jpg", output, strconv.Itoa(i))
	err := ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		panic(err)
	}
}
