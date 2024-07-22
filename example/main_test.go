package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	defer func() {
		if stop != nil {
			stop()
		}
	}()
	os.Exit(m.Run())
}
