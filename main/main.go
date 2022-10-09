package main

import (
	"os"
	"path/filepath"

	ws "github.com/GoombaG/Marking/webserver"
)

func main() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	err = os.Chdir(exPath)
	if err != nil {
		panic(err)
	}

	ws.Start()
}
