package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"golang.org/x/crypto/ssh/terminal"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flagIRC := flag.Bool("irc", false, "Output IRC color codes")
	flagIRC16 := flag.Bool("irc16", false, "Output IRC colors codes (compatibility mode)")
	flag256 := flag.Bool("256", false, "Use 256 colors")
	flag24bit := flag.Bool("24bit", false, "Use 24-bit colors")
	flagBraille := flag.Bool("braille", false, "Use braille characters") // TODO add color support

	// flagAnimated := flag.Bool("animated", false, "Animated GIF playback")
	flagSpaces := flag.Bool("spaces", false, "Use 2 spaces per pixel instead of fitting two pixels in ▀")
	flagGrayscale := flag.Bool("gray", false, "Make the image grayscale")
	flagInvert := flag.Bool("invert", false, "Invert the image colors (useful with -braille)")
	flagAutocrop := flag.Bool("crop", false, "Automatically crop out same-color or transparent borders")
	flagAutoresize := flag.Bool("autoresize", false, "Automatically downscale image so it fits your terminal")
	flagResizeW := flag.Int("width", 0, "Downscale image if greater than width")
	flagResizeH := flag.Int("height", 0, "Downscale image if greater than height")

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	mode := term16
	setMode := func(m RenderMode) {
		if mode != term16 {
			fmt.Print("Only one of -irc, -irc16, -256 or -24bit must be given")
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
	if *flagIRC16 {
		setMode(irc16)
	}
	if *flagBraille {
		setMode(braille)
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
		res := RenderToText(img, *flagGrayscale, *flagInvert, *flagAutocrop, *flagSpaces, w, h, mode)
		fmt.Print(res)
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
