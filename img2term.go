package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/crypto/ssh/terminal"
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strconv"
	"time"
)

// Render modes
const (
	term16    = iota
	term256   = iota
	term24bit = iota
	irc       = iota
)

var mode int

// Colors
var colors = [][]colorful.Color{
	// term16 (default xterm palette)
	{
		{R: 0.0, G: 0.0, B: 0.0},
		{R: 205.0 / 255.0, G: 0.0, B: 0.0},
		{R: 0.0, G: 205.0 / 255.0, B: 0.0},
		{R: 205.0 / 255.0, G: 205.0 / 255.0, B: 0.0},
		{R: 0.0, G: 0.0, B: 238.0 / 255.0},
		{R: 205.0 / 255.0, G: 0.0, B: 205.0 / 255.0},
		{R: 0.0, G: 205.0 / 255.0, B: 205.0 / 255.0},
		{R: 229.0 / 255.0, G: 229.0 / 255.0, B: 229.0 / 255.0},
		{R: 127.0 / 255.0, G: 127.0 / 255.0, B: 127.0 / 255.0},
		{R: 1.0, G: 0.0, B: 0.0},
		{R: 0.0, G: 1.0, B: 0.0},
		{R: 1.0, G: 1.0, B: 0.0},
		{R: 92.0 / 255.0, G: 92.0 / 255.0, B: 1.0},
		{R: 1.0, G: 0.0, B: 1.0},
		{R: 0.0, G: 1.0, B: 1.0},
		{R: 1.0, G: 1.0, B: 1.0},
	},
	// term256 (from https://jonasjacek.github.io/colors/ not sure how accurate)
	{
		{R: 0.0, G: 0.0, B: 0.0},
		{R: 128.0 / 255.0, G: 0.0, B: 0.0},
		{R: 0.0, G: 128.0 / 255.0, B: 0.0},
		{R: 128.0 / 255.0, G: 128.0 / 255.0, B: 0.0},
		{R: 0.0, G: 0.0, B: 128.0 / 255.0},
		{R: 128.0 / 255.0, G: 0.0, B: 128.0 / 255.0},
		{R: 0.0, G: 128.0 / 255.0, B: 128.0 / 255.0},
		{R: 192.0 / 255.0, G: 192.0 / 255.0, B: 192.0 / 255.0},
		{R: 128.0 / 255.0, G: 128.0 / 255.0, B: 128.0 / 255.0},
		{R: 1.0, G: 0.0, B: 0.0},
		{R: 0.0, G: 1.0, B: 0.0},
		{R: 1.0, G: 1.0, B: 0.0},
		{R: 0.0, G: 0.0, B: 1.0},
		{R: 1.0, G: 0.0, B: 1.0},
		{R: 0.0, G: 1.0, B: 1.0},
		{R: 1.0, G: 1.0, B: 1.0},
		{R: 0.0, G: 0.0, B: 0.0},
		{R: 0.0, G: 0.0, B: 95.0 / 255.0},
		{R: 0.0, G: 0.0, B: 135.0 / 255.0},
		{R: 0.0, G: 0.0, B: 175.0 / 255.0},
		{R: 0.0, G: 0.0, B: 215.0 / 255.0},
		{R: 0.0, G: 0.0, B: 1.0},
		{R: 0.0, G: 95.0 / 255.0, B: 0.0},
		{R: 0.0, G: 95.0 / 255.0, B: 95.0 / 255.0},
		{R: 0.0, G: 95.0 / 255.0, B: 135.0 / 255.0},
		{R: 0.0, G: 95.0 / 255.0, B: 175.0 / 255.0},
		{R: 0.0, G: 95.0 / 255.0, B: 215.0 / 255.0},
		{R: 0.0, G: 95.0 / 255.0, B: 1.0},
		{R: 0.0, G: 135.0 / 255.0, B: 0.0},
		{R: 0.0, G: 135.0 / 255.0, B: 95.0 / 255.0},
		{R: 0.0, G: 135.0 / 255.0, B: 135.0 / 255.0},
		{R: 0.0, G: 135.0 / 255.0, B: 175.0 / 255.0},
		{R: 0.0, G: 135.0 / 255.0, B: 215.0 / 255.0},
		{R: 0.0, G: 135.0 / 255.0, B: 1.0},
		{R: 0.0, G: 175.0 / 255.0, B: 0.0},
		{R: 0.0, G: 175.0 / 255.0, B: 95.0 / 255.0},
		{R: 0.0, G: 175.0 / 255.0, B: 135.0 / 255.0},
		{R: 0.0, G: 175.0 / 255.0, B: 175.0 / 255.0},
		{R: 0.0, G: 175.0 / 255.0, B: 215.0 / 255.0},
		{R: 0.0, G: 175.0 / 255.0, B: 1.0},
		{R: 0.0, G: 215.0 / 255.0, B: 0.0},
		{R: 0.0, G: 215.0 / 255.0, B: 95.0 / 255.0},
		{R: 0.0, G: 215.0 / 255.0, B: 135.0 / 255.0},
		{R: 0.0, G: 215.0 / 255.0, B: 175.0 / 255.0},
		{R: 0.0, G: 215.0 / 255.0, B: 215.0 / 255.0},
		{R: 0.0, G: 215.0 / 255.0, B: 1.0},
		{R: 0.0, G: 1.0, B: 0.0},
		{R: 0.0, G: 1.0, B: 95.0 / 255.0},
		{R: 0.0, G: 1.0, B: 135.0 / 255.0},
		{R: 0.0, G: 1.0, B: 175.0 / 255.0},
		{R: 0.0, G: 1.0, B: 215.0 / 255.0},
		{R: 0.0, G: 1.0, B: 1.0},
		{R: 95.0 / 255.0, G: 0.0, B: 0.0},
		{R: 95.0 / 255.0, G: 0.0, B: 95.0 / 255.0},
		{R: 95.0 / 255.0, G: 0.0, B: 135.0 / 255.0},
		{R: 95.0 / 255.0, G: 0.0, B: 175.0 / 255.0},
		{R: 95.0 / 255.0, G: 0.0, B: 215.0 / 255.0},
		{R: 95.0 / 255.0, G: 0.0, B: 1.0},
		{R: 95.0 / 255.0, G: 95.0 / 255.0, B: 0.0},
		{R: 95.0 / 255.0, G: 95.0 / 255.0, B: 95.0 / 255.0},
		{R: 95.0 / 255.0, G: 95.0 / 255.0, B: 135.0 / 255.0},
		{R: 95.0 / 255.0, G: 95.0 / 255.0, B: 175.0 / 255.0},
		{R: 95.0 / 255.0, G: 95.0 / 255.0, B: 215.0 / 255.0},
		{R: 95.0 / 255.0, G: 95.0 / 255.0, B: 1.0},
		{R: 95.0 / 255.0, G: 135.0 / 255.0, B: 0.0},
		{R: 95.0 / 255.0, G: 135.0 / 255.0, B: 95.0 / 255.0},
		{R: 95.0 / 255.0, G: 135.0 / 255.0, B: 135.0 / 255.0},
		{R: 95.0 / 255.0, G: 135.0 / 255.0, B: 175.0 / 255.0},
		{R: 95.0 / 255.0, G: 135.0 / 255.0, B: 215.0 / 255.0},
		{R: 95.0 / 255.0, G: 135.0 / 255.0, B: 1.0},
		{R: 95.0 / 255.0, G: 175.0 / 255.0, B: 0.0},
		{R: 95.0 / 255.0, G: 175.0 / 255.0, B: 95.0 / 255.0},
		{R: 95.0 / 255.0, G: 175.0 / 255.0, B: 135.0 / 255.0},
		{R: 95.0 / 255.0, G: 175.0 / 255.0, B: 175.0 / 255.0},
		{R: 95.0 / 255.0, G: 175.0 / 255.0, B: 215.0 / 255.0},
		{R: 95.0 / 255.0, G: 175.0 / 255.0, B: 1.0},
		{R: 95.0 / 255.0, G: 215.0 / 255.0, B: 0.0},
		{R: 95.0 / 255.0, G: 215.0 / 255.0, B: 95.0 / 255.0},
		{R: 95.0 / 255.0, G: 215.0 / 255.0, B: 135.0 / 255.0},
		{R: 95.0 / 255.0, G: 215.0 / 255.0, B: 175.0 / 255.0},
		{R: 95.0 / 255.0, G: 215.0 / 255.0, B: 215.0 / 255.0},
		{R: 95.0 / 255.0, G: 215.0 / 255.0, B: 1.0},
		{R: 95.0 / 255.0, G: 1.0, B: 0.0},
		{R: 95.0 / 255.0, G: 1.0, B: 95.0 / 255.0},
		{R: 95.0 / 255.0, G: 1.0, B: 135.0 / 255.0},
		{R: 95.0 / 255.0, G: 1.0, B: 175.0 / 255.0},
		{R: 95.0 / 255.0, G: 1.0, B: 215.0 / 255.0},
		{R: 95.0 / 255.0, G: 1.0, B: 1.0},
		{R: 135.0 / 255.0, G: 0.0, B: 0.0},
		{R: 135.0 / 255.0, G: 0.0, B: 95.0 / 255.0},
		{R: 135.0 / 255.0, G: 0.0, B: 135.0 / 255.0},
		{R: 135.0 / 255.0, G: 0.0, B: 175.0 / 255.0},
		{R: 135.0 / 255.0, G: 0.0, B: 215.0 / 255.0},
		{R: 135.0 / 255.0, G: 0.0, B: 1.0},
		{R: 135.0 / 255.0, G: 95.0 / 255.0, B: 0.0},
		{R: 135.0 / 255.0, G: 95.0 / 255.0, B: 95.0 / 255.0},
		{R: 135.0 / 255.0, G: 95.0 / 255.0, B: 135.0 / 255.0},
		{R: 135.0 / 255.0, G: 95.0 / 255.0, B: 175.0 / 255.0},
		{R: 135.0 / 255.0, G: 95.0 / 255.0, B: 215.0 / 255.0},
		{R: 135.0 / 255.0, G: 95.0 / 255.0, B: 1.0},
		{R: 135.0 / 255.0, G: 135.0 / 255.0, B: 0.0},
		{R: 135.0 / 255.0, G: 135.0 / 255.0, B: 95.0 / 255.0},
		{R: 135.0 / 255.0, G: 135.0 / 255.0, B: 135.0 / 255.0},
		{R: 135.0 / 255.0, G: 135.0 / 255.0, B: 175.0 / 255.0},
		{R: 135.0 / 255.0, G: 135.0 / 255.0, B: 215.0 / 255.0},
		{R: 135.0 / 255.0, G: 135.0 / 255.0, B: 1.0},
		{R: 135.0 / 255.0, G: 175.0 / 255.0, B: 0.0},
		{R: 135.0 / 255.0, G: 175.0 / 255.0, B: 95.0 / 255.0},
		{R: 135.0 / 255.0, G: 175.0 / 255.0, B: 135.0 / 255.0},
		{R: 135.0 / 255.0, G: 175.0 / 255.0, B: 175.0 / 255.0},
		{R: 135.0 / 255.0, G: 175.0 / 255.0, B: 215.0 / 255.0},
		{R: 135.0 / 255.0, G: 175.0 / 255.0, B: 1.0},
		{R: 135.0 / 255.0, G: 215.0 / 255.0, B: 0.0},
		{R: 135.0 / 255.0, G: 215.0 / 255.0, B: 95.0 / 255.0},
		{R: 135.0 / 255.0, G: 215.0 / 255.0, B: 135.0 / 255.0},
		{R: 135.0 / 255.0, G: 215.0 / 255.0, B: 175.0 / 255.0},
		{R: 135.0 / 255.0, G: 215.0 / 255.0, B: 215.0 / 255.0},
		{R: 135.0 / 255.0, G: 215.0 / 255.0, B: 1.0},
		{R: 135.0 / 255.0, G: 1.0, B: 0.0},
		{R: 135.0 / 255.0, G: 1.0, B: 95.0 / 255.0},
		{R: 135.0 / 255.0, G: 1.0, B: 135.0 / 255.0},
		{R: 135.0 / 255.0, G: 1.0, B: 175.0 / 255.0},
		{R: 135.0 / 255.0, G: 1.0, B: 215.0 / 255.0},
		{R: 135.0 / 255.0, G: 1.0, B: 1.0},
		{R: 175.0 / 255.0, G: 0.0, B: 0.0},
		{R: 175.0 / 255.0, G: 0.0, B: 95.0 / 255.0},
		{R: 175.0 / 255.0, G: 0.0, B: 135.0 / 255.0},
		{R: 175.0 / 255.0, G: 0.0, B: 175.0 / 255.0},
		{R: 175.0 / 255.0, G: 0.0, B: 215.0 / 255.0},
		{R: 175.0 / 255.0, G: 0.0, B: 1.0},
		{R: 175.0 / 255.0, G: 95.0 / 255.0, B: 0.0},
		{R: 175.0 / 255.0, G: 95.0 / 255.0, B: 95.0 / 255.0},
		{R: 175.0 / 255.0, G: 95.0 / 255.0, B: 135.0 / 255.0},
		{R: 175.0 / 255.0, G: 95.0 / 255.0, B: 175.0 / 255.0},
		{R: 175.0 / 255.0, G: 95.0 / 255.0, B: 215.0 / 255.0},
		{R: 175.0 / 255.0, G: 95.0 / 255.0, B: 1.0},
		{R: 175.0 / 255.0, G: 135.0 / 255.0, B: 0.0},
		{R: 175.0 / 255.0, G: 135.0 / 255.0, B: 95.0 / 255.0},
		{R: 175.0 / 255.0, G: 135.0 / 255.0, B: 135.0 / 255.0},
		{R: 175.0 / 255.0, G: 135.0 / 255.0, B: 175.0 / 255.0},
		{R: 175.0 / 255.0, G: 135.0 / 255.0, B: 215.0 / 255.0},
		{R: 175.0 / 255.0, G: 135.0 / 255.0, B: 1.0},
		{R: 175.0 / 255.0, G: 175.0 / 255.0, B: 0.0},
		{R: 175.0 / 255.0, G: 175.0 / 255.0, B: 95.0 / 255.0},
		{R: 175.0 / 255.0, G: 175.0 / 255.0, B: 135.0 / 255.0},
		{R: 175.0 / 255.0, G: 175.0 / 255.0, B: 175.0 / 255.0},
		{R: 175.0 / 255.0, G: 175.0 / 255.0, B: 215.0 / 255.0},
		{R: 175.0 / 255.0, G: 175.0 / 255.0, B: 1.0},
		{R: 175.0 / 255.0, G: 215.0 / 255.0, B: 0.0},
		{R: 175.0 / 255.0, G: 215.0 / 255.0, B: 95.0 / 255.0},
		{R: 175.0 / 255.0, G: 215.0 / 255.0, B: 135.0 / 255.0},
		{R: 175.0 / 255.0, G: 215.0 / 255.0, B: 175.0 / 255.0},
		{R: 175.0 / 255.0, G: 215.0 / 255.0, B: 215.0 / 255.0},
		{R: 175.0 / 255.0, G: 215.0 / 255.0, B: 1.0},
		{R: 175.0 / 255.0, G: 1.0, B: 0.0},
		{R: 175.0 / 255.0, G: 1.0, B: 95.0 / 255.0},
		{R: 175.0 / 255.0, G: 1.0, B: 135.0 / 255.0},
		{R: 175.0 / 255.0, G: 1.0, B: 175.0 / 255.0},
		{R: 175.0 / 255.0, G: 1.0, B: 215.0 / 255.0},
		{R: 175.0 / 255.0, G: 1.0, B: 1.0},
		{R: 215.0 / 255.0, G: 0.0, B: 0.0},
		{R: 215.0 / 255.0, G: 0.0, B: 95.0 / 255.0},
		{R: 215.0 / 255.0, G: 0.0, B: 135.0 / 255.0},
		{R: 215.0 / 255.0, G: 0.0, B: 175.0 / 255.0},
		{R: 215.0 / 255.0, G: 0.0, B: 215.0 / 255.0},
		{R: 215.0 / 255.0, G: 0.0, B: 1.0},
		{R: 215.0 / 255.0, G: 95.0 / 255.0, B: 0.0},
		{R: 215.0 / 255.0, G: 95.0 / 255.0, B: 95.0 / 255.0},
		{R: 215.0 / 255.0, G: 95.0 / 255.0, B: 135.0 / 255.0},
		{R: 215.0 / 255.0, G: 95.0 / 255.0, B: 175.0 / 255.0},
		{R: 215.0 / 255.0, G: 95.0 / 255.0, B: 215.0 / 255.0},
		{R: 215.0 / 255.0, G: 95.0 / 255.0, B: 1.0},
		{R: 215.0 / 255.0, G: 135.0 / 255.0, B: 0.0},
		{R: 215.0 / 255.0, G: 135.0 / 255.0, B: 95.0 / 255.0},
		{R: 215.0 / 255.0, G: 135.0 / 255.0, B: 135.0 / 255.0},
		{R: 215.0 / 255.0, G: 135.0 / 255.0, B: 175.0 / 255.0},
		{R: 215.0 / 255.0, G: 135.0 / 255.0, B: 215.0 / 255.0},
		{R: 215.0 / 255.0, G: 135.0 / 255.0, B: 1.0},
		{R: 215.0 / 255.0, G: 175.0 / 255.0, B: 0.0},
		{R: 215.0 / 255.0, G: 175.0 / 255.0, B: 95.0 / 255.0},
		{R: 215.0 / 255.0, G: 175.0 / 255.0, B: 135.0 / 255.0},
		{R: 215.0 / 255.0, G: 175.0 / 255.0, B: 175.0 / 255.0},
		{R: 215.0 / 255.0, G: 175.0 / 255.0, B: 215.0 / 255.0},
		{R: 215.0 / 255.0, G: 175.0 / 255.0, B: 1.0},
		{R: 215.0 / 255.0, G: 215.0 / 255.0, B: 0.0},
		{R: 215.0 / 255.0, G: 215.0 / 255.0, B: 95.0 / 255.0},
		{R: 215.0 / 255.0, G: 215.0 / 255.0, B: 135.0 / 255.0},
		{R: 215.0 / 255.0, G: 215.0 / 255.0, B: 175.0 / 255.0},
		{R: 215.0 / 255.0, G: 215.0 / 255.0, B: 215.0 / 255.0},
		{R: 215.0 / 255.0, G: 215.0 / 255.0, B: 1.0},
		{R: 215.0 / 255.0, G: 1.0, B: 0.0},
		{R: 215.0 / 255.0, G: 1.0, B: 95.0 / 255.0},
		{R: 215.0 / 255.0, G: 1.0, B: 135.0 / 255.0},
		{R: 215.0 / 255.0, G: 1.0, B: 175.0 / 255.0},
		{R: 215.0 / 255.0, G: 1.0, B: 215.0 / 255.0},
		{R: 215.0 / 255.0, G: 1.0, B: 1.0},
		{R: 1.0, G: 0.0, B: 0.0},
		{R: 1.0, G: 0.0, B: 95.0 / 255.0},
		{R: 1.0, G: 0.0, B: 135.0 / 255.0},
		{R: 1.0, G: 0.0, B: 175.0 / 255.0},
		{R: 1.0, G: 0.0, B: 215.0 / 255.0},
		{R: 1.0, G: 0.0, B: 1.0},
		{R: 1.0, G: 95.0 / 255.0, B: 0.0},
		{R: 1.0, G: 95.0 / 255.0, B: 95.0 / 255.0},
		{R: 1.0, G: 95.0 / 255.0, B: 135.0 / 255.0},
		{R: 1.0, G: 95.0 / 255.0, B: 175.0 / 255.0},
		{R: 1.0, G: 95.0 / 255.0, B: 215.0 / 255.0},
		{R: 1.0, G: 95.0 / 255.0, B: 1.0},
		{R: 1.0, G: 135.0 / 255.0, B: 0.0},
		{R: 1.0, G: 135.0 / 255.0, B: 95.0 / 255.0},
		{R: 1.0, G: 135.0 / 255.0, B: 135.0 / 255.0},
		{R: 1.0, G: 135.0 / 255.0, B: 175.0 / 255.0},
		{R: 1.0, G: 135.0 / 255.0, B: 215.0 / 255.0},
		{R: 1.0, G: 135.0 / 255.0, B: 1.0},
		{R: 1.0, G: 175.0 / 255.0, B: 0.0},
		{R: 1.0, G: 175.0 / 255.0, B: 95.0 / 255.0},
		{R: 1.0, G: 175.0 / 255.0, B: 135.0 / 255.0},
		{R: 1.0, G: 175.0 / 255.0, B: 175.0 / 255.0},
		{R: 1.0, G: 175.0 / 255.0, B: 215.0 / 255.0},
		{R: 1.0, G: 175.0 / 255.0, B: 1.0},
		{R: 1.0, G: 215.0 / 255.0, B: 0.0},
		{R: 1.0, G: 215.0 / 255.0, B: 95.0 / 255.0},
		{R: 1.0, G: 215.0 / 255.0, B: 135.0 / 255.0},
		{R: 1.0, G: 215.0 / 255.0, B: 175.0 / 255.0},
		{R: 1.0, G: 215.0 / 255.0, B: 215.0 / 255.0},
		{R: 1.0, G: 215.0 / 255.0, B: 1.0},
		{R: 1.0, G: 1.0, B: 0.0},
		{R: 1.0, G: 1.0, B: 95.0 / 255.0},
		{R: 1.0, G: 1.0, B: 135.0 / 255.0},
		{R: 1.0, G: 1.0, B: 175.0 / 255.0},
		{R: 1.0, G: 1.0, B: 215.0 / 255.0},
		{R: 1.0, G: 1.0, B: 1.0},
		{R: 8.0 / 255.0, G: 8.0 / 255.0, B: 8.0 / 255.0},
		{R: 18.0 / 255.0, G: 18.0 / 255.0, B: 18.0 / 255.0},
		{R: 28.0 / 255.0, G: 28.0 / 255.0, B: 28.0 / 255.0},
		{R: 38.0 / 255.0, G: 38.0 / 255.0, B: 38.0 / 255.0},
		{R: 48.0 / 255.0, G: 48.0 / 255.0, B: 48.0 / 255.0},
		{R: 58.0 / 255.0, G: 58.0 / 255.0, B: 58.0 / 255.0},
		{R: 68.0 / 255.0, G: 68.0 / 255.0, B: 68.0 / 255.0},
		{R: 78.0 / 255.0, G: 78.0 / 255.0, B: 78.0 / 255.0},
		{R: 88.0 / 255.0, G: 88.0 / 255.0, B: 88.0 / 255.0},
		{R: 98.0 / 255.0, G: 98.0 / 255.0, B: 98.0 / 255.0},
		{R: 108.0 / 255.0, G: 108.0 / 255.0, B: 108.0 / 255.0},
		{R: 118.0 / 255.0, G: 118.0 / 255.0, B: 118.0 / 255.0},
		{R: 128.0 / 255.0, G: 128.0 / 255.0, B: 128.0 / 255.0},
		{R: 138.0 / 255.0, G: 138.0 / 255.0, B: 138.0 / 255.0},
		{R: 148.0 / 255.0, G: 148.0 / 255.0, B: 148.0 / 255.0},
		{R: 158.0 / 255.0, G: 158.0 / 255.0, B: 158.0 / 255.0},
		{R: 168.0 / 255.0, G: 168.0 / 255.0, B: 168.0 / 255.0},
		{R: 178.0 / 255.0, G: 178.0 / 255.0, B: 178.0 / 255.0},
		{R: 188.0 / 255.0, G: 188.0 / 255.0, B: 188.0 / 255.0},
		{R: 198.0 / 255.0, G: 198.0 / 255.0, B: 198.0 / 255.0},
		{R: 208.0 / 255.0, G: 208.0 / 255.0, B: 208.0 / 255.0},
		{R: 218.0 / 255.0, G: 218.0 / 255.0, B: 218.0 / 255.0},
		{R: 228.0 / 255.0, G: 228.0 / 255.0, B: 228.0 / 255.0},
		{R: 238.0 / 255.0, G: 238.0 / 255.0, B: 238.0 / 255.0},
	},
	// term24bit
	{}, // no palette
	// irc
	{
		{R: 1.0, G: 1.0, B: 1.0},
		{R: 0.0, G: 0.0, B: 0.0},
		{R: 0.0, G: 0.0, B: 127.0 / 255.0},
		{R: 0.0, G: 147.0 / 255.0, B: 0.0},
		{R: 1.0, G: 0.0, B: 0.0},
		{R: 127.0 / 255.0, G: 0.0, B: 0.0},
		{R: 156.0 / 255.0, G: 0.0, B: 156.0 / 255.0},
		{R: 252.0 / 255.0, G: 127.0 / 255.0, B: 0.0},
		{R: 1.0, G: 1.0, B: 0.0},
		{R: 0.0, G: 252.0 / 255.0, B: 0.0},
		{R: 0.0, G: 147.0 / 255.0, B: 147.0 / 255.0},
		{R: 0.0, G: 1.0, B: 1.0},
		{R: 0.0, G: 0.0, B: 252.0 / 255.0},
		{R: 1.0, G: 0.0, B: 1.0},
		{R: 127.0 / 255.0, G: 127.0 / 255.0, B: 127.0 / 255.0},
		{R: 210.0, G: 210.0, B: 210.0},
	},
}

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

// Comparison function between two colors
func ColorDistance(c1 colorful.Color, c2 colorful.Color) float64 {
	if mode == irc || mode == term16 {
		return c1.DistanceLuv(c2)
	} else {
		return c1.DistanceLab(c2)
	}
}

// Returns the palette index closest to the color in the current mode
func ColorToPalette(color colorful.Color) int {
	if mode == term24bit {
		log.Fatal("Should never happen")
	}
	result := 0
	dist := ColorDistance(color, colors[mode][0])
	for i := 0; i < len(colors[mode]); i++ {
		d := ColorDistance(color, colors[mode][i])
		if d < dist {
			dist = d
			result = i
		}
	}
	return result
}

func GetPixels(img image.Image) [][]colorful.Color {
	bounds := img.Bounds()
	img_size := bounds.Size()
	w, h := img_size.X, img_size.Y
	pixels := make([][]colorful.Color, h)
	for i := 0; i < h; i++ {
		pixels[i] = make([]colorful.Color, w)
	}

	di := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		dj := 0
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Workaround panic in go-colorful when alpha is 0
			_, _, _, alpha := img.At(x, y).RGBA()
			if alpha == 0 {
				// Default to white when transparent
				pixels[di][dj] = colorful.Color{R: 1.0, G: 1.0, B: 1.0}
			} else {
				pixels[di][dj] = colorful.MakeColor(img.At(x, y))
			}
			dj++
		}
		di++
	}

	return pixels
}

func Render(colors [][]colorful.Color, use_spaces bool) string {
	w, h := len(colors[0]), len(colors)
	var buffer bytes.Buffer
	step := 2
	if use_spaces {
		step = 1
	}

	for y := 0; y < h; y += step {
		var prev_col string
		var prev_bg_col string
		for x := 0; x < w; x++ {
			next_fg_col := ""
			if mode != term24bit {
				next_fg_col = strconv.Itoa(ColorToPalette(colors[y][x]))
			} else {
				r, g, b := colors[y][x].RGB255()
				next_fg_col = fmt.Sprintf("%d;%d;%d", r, g, b)
			}
			next_bg_col := ""
			// Process two vertical pixels per column unless we're printing spaces
			if !use_spaces {
				if y+1 < h {
					if mode != term24bit {
						next_bg_col = strconv.Itoa(ColorToPalette(colors[y+1][x]))
					} else {
						r, g, b := colors[y+1][x].RGB255()
						next_bg_col = fmt.Sprintf("%d;%d;%d", r, g, b)
					}
				}
			} else {
				next_bg_col = next_fg_col
				if mode != term24bit {
					next_fg_col = "0"
				} else {
					next_fg_col = "0;0;0"
				}
			}
			// This doesn't actually get printed
			next_col := next_fg_col + "," + next_bg_col

			if prev_col != next_col {
				// Color start header
				if mode == irc {
					buffer.WriteString("\x03")
				} else {
					// Foreground color selector is 38;
					buffer.WriteString("\x1B[38;")
					if mode == term24bit {
						buffer.WriteString("2;")
					} else {
						buffer.WriteString("5;")
					}
				}
				buffer.WriteString(next_fg_col)
				if next_bg_col != "" && prev_bg_col != next_bg_col {
					if mode != irc {
						// Background color selector is 48;
						buffer.WriteString("m\x1B[48;")
						if mode == term24bit {
							buffer.WriteString("2;")
						} else {
							buffer.WriteString("5;")
						}
					} else {
						// Irc colors are \x03FF,BB
						buffer.WriteString(",")
					}
					buffer.WriteString(next_bg_col)
					prev_bg_col = next_bg_col
				}
				prev_col = next_col
				if mode != irc {
					// ANSI escape sequences are terminated by 'm'
					buffer.WriteString("m")
				}
			}
			if use_spaces {
				// Print two spaces to keep aspect ratio
				buffer.WriteString("  ")
			} else {
				buffer.WriteString("▀")
			}
		}
		if mode != irc {
			// Reset colors at end of line
			buffer.WriteString("\x1B[0m")
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func main() {
	flagIRC := flag.Bool("irc", false, "Output IRC color codes")
	flag256 := flag.Bool("256", false, "Use 256 colors")
	flag24bit := flag.Bool("24bit", false, "Use 24-bit colors")

	flagAnimated := flag.Bool("animated", false, "Animated GIF playback")
	flagSpaces := flag.Bool("spaces", false, "Use 2 spaces per pixel instead of fitting two pixels in ▀")
	flagAutoresize := flag.Bool("autoresize", false, "Automatically downscale image so it fits your terminal")
	flagResizeW := flag.Int("width", 0, "Downscale image if greater than width")
	flagResizeH := flag.Int("height", 0, "Downscale image if greater than height")

	flag.Parse()

	mode = term16
	setMode := func(m int) {
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
		h -= 3 // Some padding for prompts and shit
	}
	if *flagSpaces {
		w /= 2
	} else {
		h *= 2
	}
	if *flagAnimated {
		if len(flag.Args()) > 1 {
			log.Fatal("Only one file supported with -animate")
		}
		anim := DecodeGIF(flag.Arg(0))
		frames := make([]string, len(anim.Image))
		for i := 0; i < len(anim.Image); i++ {
			if anim.Disposal[i] != gif.DisposalBackground {
				// TODO implement gif.DisposalNone and gif.DisposalPrevious
				// TODO make sure every gif frame has the correct output resolution
				// TODO prettier single-line progress indicator
				log.Print("[WARN] Unsupported gif frame disposal method (TODO)")
			}
			fmt.Printf("Processing frame %d/%d... ", i, len(anim.Image))
			var img image.Image
			img = anim.Image[i]
			if w != 0 || h != 0 {
				if h != 0 && img.Bounds().Size().Y > h {
					img = imaging.Resize(img, 0, h, imaging.NearestNeighbor)
				}
				if w != 0 && img.Bounds().Size().X > w {
					img = imaging.Resize(img, w, 0, imaging.NearestNeighbor)
				}
			}
			frames[i] = Render(GetPixels(img), *flagSpaces)
			fmt.Printf("Done!\n")
		}
		fmt.Print("\x1B[2J") // Clear screen
		frame := 0
		for {
			fmt.Print("\x1B[1;1H") // Move cursor to 1,1
			fmt.Println(frames[frame])
			delay := time.Duration(anim.Delay[frame]*10) * time.Millisecond
			time.Sleep(delay)
			frame++
			if frame >= len(anim.Image) {
				frame = 0
			}
		}
	} else {
		for _, file := range flag.Args() {
			img := DecodeImage(file)
			if w != 0 || h != 0 {
				if h != 0 && img.Bounds().Size().Y > h {
					img = imaging.Resize(img, 0, h, imaging.NearestNeighbor)
				}
				if w != 0 && img.Bounds().Size().X > w {
					img = imaging.Resize(img, w, 0, imaging.NearestNeighbor)
				}
			}
			fmt.Println(Render(GetPixels(img), *flagSpaces))
		}
	}
}
