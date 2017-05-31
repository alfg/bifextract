package main

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"testing"
)

var testfile = "tears-of-steel.bif"
var f *os.File

func TestMain(m *testing.M) {
	setup()

	retCode := m.Run()

	teardown()
	os.Exit(retCode)
}

func setup() {
	buf, err := os.Open(testfile)
	if err != nil {
		panic(err)
	}
	f = buf
}

func teardown() {
	f.Close()
}

func Test_checkBIF(t *testing.T) {
	err := checkBIF(f)
	if err != nil {
		t.Error()
	}
}

func Test_NewBIF(t *testing.T) {
	_, err := NewBIF(f)
	if err != nil {
		t.Error()
	}
}

func Test_BIF_getVersion(t *testing.T) {
	bif, err := NewBIF(f)
	if err != nil {
		t.Error()
	}
	version := bif.getVersion()
	if version != 0 {
		t.Error()
	}
}

func Test_BIF_getFramesCount(t *testing.T) {
	bif, err := NewBIF(f)
	if err != nil {
		t.Error()
	}
	framesCount := bif.getFramesCount()
	if framesCount != 185 {
		t.Error()
	}
}

func Test_BIF_getFramewiseSeparation(t *testing.T) {
	bif, err := NewBIF(f)
	if err != nil {
		t.Error()
	}
	fs := bif.getFramewiseSeparation()
	if fs != 1000 {
		t.Error()
	}
}

func Test_BIF_readFrame(t *testing.T) {
	bif, err := NewBIF(f)
	if err != nil {
		t.Error()
	}
	timestamp1, offset1 := bif.readFrame(64)
	if timestamp1 != 1 || offset1 != 1552 {
		t.Error()
	}

	timestamp2, offset2 := bif.readFrame(72)
	if timestamp2 != 2 || offset2 != 2861 {
		t.Error()
	}
}

func Test_BIF_getFrameImage(t *testing.T) {
	bif, err := NewBIF(f)
	if err != nil {
		t.Error()
	}

	// Get first frame.
	frameSize := 2861 - 1552 // frameSize = nextFrameOffet - currentOffset
	img := bif.getFrameImage(1552, frameSize)

	// Test if valid base64 string.
	match, _ := regexp.MatchString("^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$", img)
	if !match {
		t.Error()
	}

	// Test if valid jpeg image.
	dec, _ := base64.StdEncoding.DecodeString(img)
	contentType := http.DetectContentType(dec)
	if contentType != "image/jpeg" {
		t.Error()
	}
}

func Test_BIF_createFrameImage(t *testing.T) {
	bif, err := NewBIF(f)
	if err != nil {
		t.Error()
	}

	// Get first frame.
	frameSize := 2861 - 1552 // frameSize = nextFrameOffet - currentOffset
	out := "test"

	// Create image.
	bif.createFrameImage(0, 2861, frameSize, out)

	// Test image.
	dat, err := ioutil.ReadFile("test/frame_0.jpg")
	if err != nil {
		t.Error()
	}

	contentType := http.DetectContentType(dat)
	if contentType != "image/jpeg" {
		t.Error()
	}
}
