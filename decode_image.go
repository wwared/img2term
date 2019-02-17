package main

import (
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)


func DecodeImage(path string) image.Image {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

func DecodeGIF(path string) *gif.GIF {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	img, err := gif.DecodeAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return img
}
