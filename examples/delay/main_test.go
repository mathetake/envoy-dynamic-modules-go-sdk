package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestFoo(t *testing.T) {
	m := &moduleContext{}
	fmt.Println(m)
	fmt.Println("come on")
}
