package main

import (
	"os"
	"testing"
)

func TestIntegration(t *testing.T) {
	tests := map[string]func(*testing.T){
		"should finish succesfully": func(t *testing.T) {
			// given
			os.Args = []string{"sherlock", "--file", "../data/demo.json"}

			// then
			main()
		},
		"should panic when file is not found": func(t *testing.T) {
			// given
			os.Args = []string{"sherlock", "--file", "file/not/found"}

			// then
			assertPanic(t, main)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("a panic was expected")
		}
	}()
	f()
}
