// bifextract - CLI utility for extracting images from a BIF file.
// https://sdkdocs.roku.com/display/sdkdoc/Trick+Mode+Support
package main

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	usage = `
Usage:
	bifextract <file-path|url> <output-dir>
`
	version = "0.1.0"
)

var (
	file      string
	outputDir string
)

// BIF represents BIF file data.
type BIF struct {
	File *os.File
}

// Frame represents each frame in BIF.
type Frame struct {
	Timestamp uint32
	Offset    uint32
}

// Initialize and parse flags/arguments
func init() {
	if len(os.Args[1:]) < 1 {
		fmt.Println("Please provide a BIF path or URL.")
		fmt.Println(usage)
		os.Exit(1)
	}

	// TODO: Support URL paths.
	// isHTTP := strings.HasPrefix(os.Args[len(os.Args)-1], "http")
}

func main() {
	args := os.Args[1:]

	option := args[0]
	if option == "version" {
		fmt.Printf("Version: %s\n", version)
	} else {
		file = option
		outputDir = args[1]
		extractBIF()
	}
}

func extractBIF() {

	// Open file.
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	// Create output dir.
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, os.ModePerm)
	}

	// Load BIF instance.
	bif, _ := NewBIF(f)

	// Version
	v := bif.getVersion()
	fmt.Printf("BIF Version: %d\n", v)

	// Frame Count
	fc := bif.getFramesCount()
	fmt.Printf("Number of frames: %d\n", fc)

	// Framewise Separation
	fs := bif.getFramewiseSeparation()
	fmt.Printf("Framewise Separation: %d ms\n", fs)

	// Get frames
	var byteIndex int64 = 64 // BIF index starts at byte 64.
	var frames []Frame

	for i := 0; i < fc; i++ {
		ts, offset := bif.readFrame(byteIndex)
		frame := Frame{
			Timestamp: ts * fs,
			Offset:    offset,
		}
		frames = append(frames, frame)
		byteIndex += 8 // Next frame every 8 bytes.
	}

	// Generate frame images to output dir.
	fmt.Printf("Generating %d frames...\n", len(frames))
	for k, v := range frames {
		var nextOffset uint32
		if k == len(frames)-1 {
			// nextOffset = frames[len(frames)-1].Offset
			fmt.Println("Finished.")
		} else {
			nextOffset = frames[k+1].Offset

			// Calculate frame image size from next offset.
			frameLen := nextOffset - v.Offset

			// Create image.
			bif.createFrameImage(k, int64(v.Offset), int(frameLen))
		}
	}
	f.Close()
}

// NewBIF Creates a BIF instance.
func NewBIF(f *os.File) (*BIF, error) {
	// Validate filetype
	err := checkBIF(f)
	if err != nil {
		return &BIF{
			File: f,
		}, err
	}

	return &BIF{
		File: f,
	}, nil
}

func checkBIF(f *os.File) error {
	b := make([]byte, 8)
	_, err := f.ReadAt(b, 0)
	if err != nil {
		panic(err)
	}
	magic := string(b)
	isBIF := strings.Contains(magic, "BIF")
	if !isBIF {
		return errors.New("invalid BIF file")
	}
	return nil
}

func (b *BIF) getVersion() uint32 {
	f := b.File
	f.Seek(8, 0)
	buf := make([]byte, 4)
	_, err := f.Read(buf)
	if err != nil {
		panic(err)
	}
	version := binary.LittleEndian.Uint32(buf)
	return version
}

func (b *BIF) getFramesCount() int {
	f := b.File
	f.Seek(12, 0)
	buf := make([]byte, 4)
	_, err := f.Read(buf)
	if err != nil {
		panic(err)
	}
	numFrames := binary.LittleEndian.Uint32(buf)
	return int(numFrames)
}

func (b *BIF) getFramewiseSeparation() uint32 {
	f := b.File
	f.Seek(16, 0)
	buf := make([]byte, 4)
	_, err := f.Read(buf)
	if err != nil {
		panic(err)
	}
	framewiseSeparation := binary.LittleEndian.Uint32(buf)
	return framewiseSeparation
}

func (b *BIF) readFrame(offset int64) (uint32, uint32) {
	f := b.File
	f.Seek(offset, 0)
	buf := make([]byte, 4)
	f.Read(buf)
	frameTimestamp := binary.LittleEndian.Uint32(buf)

	f.Seek(offset+4, 0)
	b2 := make([]byte, 4)
	f.Read(b2)
	frameOffset := binary.LittleEndian.Uint32(b2)
	return frameTimestamp, frameOffset
}

func (b *BIF) getFrameImage(offset int64, len int) string {
	f := b.File
	f.Seek(offset, 0)
	buf := make([]byte, len)
	f.Read(buf)
	enc := base64.StdEncoding.EncodeToString(buf)
	return enc
}

func (b *BIF) createFrameImage(i int, offset int64, len int) {
	f := b.File
	f.Seek(offset, 0)
	buf := make([]byte, len)
	f.Read(buf)

	filename := fmt.Sprintf("%s/frame_%s.jpg", outputDir, strconv.Itoa(i))
	err := ioutil.WriteFile(filename, buf, 0644)
	if err != nil {
		panic(err)
	}
}
