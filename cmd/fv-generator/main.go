package main

import (
	"fv-generator/fv"
	"log"
)

func main() {
	fv, err := fv.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = fv.SavePDF()
	if err != nil {
		log.Fatal(err)
	}
}
