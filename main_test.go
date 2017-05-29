package main

import (
	"os"
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
