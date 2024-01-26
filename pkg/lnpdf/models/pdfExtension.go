package model

import (
	"encoding/json"
	"os"
)

func (this *Pdf) FromFile(path *string) {
	file, err := os.ReadFile(*path)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(file, this); err != nil {
		panic(err)
	}
}
