package e2e

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := os.Chdir(".."); err != nil {
		log.Fatal(err)
	}

	// Check if `envoy` is installed.
	if _, err := os.Stat("envoy"); os.IsNotExist(err) {
		log.Fatal("envoy binary not found. Please install it from https://github.com/envoyproxyx/envoyx/pkgs/container/envoy")
	}

	os.Exit(m.Run())
}
