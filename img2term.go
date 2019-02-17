package main

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
)

func main() {
	flagIRC := flag.Bool("irc", false, "Output IRC color codes")
	flag256 := flag.Bool("256", false, "Use 256 colors")
	flag24bit := flag.Bool("24bit", false, "Use 24-bit colors")

	// flagAnimated := flag.Bool("animated", false, "Animated GIF playback")
	flagSpaces := flag.Bool("spaces", false, "Use 2 spaces per pixel instead of fitting two pixels in ▀")
	flagGrayscale := flag.Bool("gray", false, "Make the image grayscale")
	flagAutocrop := flag.Bool("crop", false, "Automatically crop out same-color or transparent borders")
	flagAutoresize := flag.Bool("autoresize", false, "Automatically downscale image so it fits your terminal")
	flagResizeW := flag.Int("width", 0, "Downscale image if greater than width")
	flagResizeH := flag.Int("height", 0, "Downscale image if greater than height")

	flag.Parse()

	mode := term16
	setMode := func(m RenderMode) {
		if mode != term16 {
			fmt.Print("Only one of -irc, -256 or -24bit must be given")
			os.Exit(1)
		}
		mode = m
	}
	if *flag256 {
		setMode(term256)
	}
	if *flag24bit {
		setMode(term24bit)
	}
	if *flagIRC {
		setMode(irc)
	}
	w, h := *flagResizeW, *flagResizeH
	if *flagAutoresize {
		var err error
		w, h, err = terminal.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			log.Fatal(err)
		}
		h -= 3 // Some vertical padding for shell prompts
	}
	if *flagSpaces {
		w /= 2
	} else {
		h *= 2
	}
	for _, file := range flag.Args() {
		img := DecodeImage(file)
		res := RenderToText(img, *flagGrayscale, *flagAutocrop, *flagSpaces, w, h, mode)
		fmt.Println(res)
	}
}
