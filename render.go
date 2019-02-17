package main

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/lucasb-eyer/go-colorful"
	"image"
	"image/color"
	"strconv"
)

type Pixel struct {
	color colorful.Color
	alpha uint32
}

// Rendering entrypoint
func RenderToText(img image.Image, grayscale bool, autocrop bool, use_spaces bool, width int, height int, mode RenderMode) string {
	if grayscale {
		img = Grayscale(img)
	}
	if autocrop {
		img = CropBorders(img)
	}
	if width != 0 || height != 0 {
		if height != 0 && img.Bounds().Size().Y > height {
			img = imaging.Resize(img, 0, height, imaging.NearestNeighbor)
		}
		if width != 0 && img.Bounds().Size().X > width {
			img = imaging.Resize(img, width, 0, imaging.NearestNeighbor)
		}
	}
	return Render(mode, use_spaces, GetPixels(img))
}

//
// Preprocessing filters
//

func Grayscale(img image.Image) image.Image {
	bounds := img.Bounds()
	gray := image.NewGray16(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray.Set(x, y, img.At(x, y))
		}
	}
	return gray
}

// PSA this code is repetitive and ugly
func CropBorders(img image.Image) image.Image {
	bounds := img.Bounds()
	ix, iy, w, h := bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y
	background := img.At(ix, iy)
	for {
		done := false
		for y := iy; y < h; y++ {
			if !IsTransparent(img.At(ix, y)) && img.At(ix, y) != background {
				done = true
				break
			}
		}
		if done {
			break
		}
		ix++
		if ix == w {
			break
		}
	}
	for {
		done := false
		for x := ix; x < w; x++ {
			if !IsTransparent(img.At(x, iy)) && img.At(x, iy) != background {
				done = true
				break
			}
		}
		if done {
			break
		}
		iy++
		if iy == h {
			break
		}
	}
	for {
		done := false
		for y := iy; y < h; y++ {
			if !IsTransparent(img.At(w-1, y)) && img.At(w-1, y) != background {
				done = true
				break
			}
		}
		if done {
			break
		}
		w--
		if w == 0 {
			break
		}
	}
	for {
		done := false
		for x := ix; x < w; x++ {
			if !IsTransparent(img.At(x, h-1)) && img.At(x, h-1) != background {
				done = true
				break
			}
		}
		if done {
			break
		}
		h--
		if h == 0 {
			break
		}
	}
	if ix == bounds.Min.X && iy == bounds.Min.Y && w == bounds.Max.X && h == bounds.Max.Y {
		return img
	}
	result := image.NewRGBA(image.Rect(0, 0, w-ix, h-iy))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			result.Set(x, y, img.At(ix + x, iy + y))
		}
	}
	return result
}

//
// Misc utility functions
//

// Extracts Pixels from an image.Image
func GetPixels(img image.Image) [][]Pixel {
	bounds := img.Bounds()
	img_size := bounds.Size()
	w, h := img_size.X, img_size.Y
	pixels := make([][]Pixel, h)
	for i := 0; i < h; i++ {
		pixels[i] = make([]Pixel, w)
	}

	di := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		dj := 0
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			orig_color := img.At(x, y)
			color, alpha_flag := colorful.MakeColor(orig_color)
			if !alpha_flag {
				color = colorful.Color{R: 1.0, G: 1.0, B: 1.0}
			}
			_, _, _, alpha := orig_color.RGBA()
			pixels[di][dj] = Pixel{
				color: color,
				alpha: alpha,
			}
			dj++
		}
		di++
	}

	return pixels
}

func ColorDistance(mode RenderMode, c1 colorful.Color, c2 colorful.Color) float64 {
	if mode == irc || mode == term16 { // heuristic; I think it looks nicer
		return c1.DistanceLuv(c2)
	} else {
		return c1.DistanceLab(c2)
	}
}

func IsTransparent(c color.Color) bool {
	_, _, _, alpha := c.RGBA()
	return alpha < TRANSPARENCY_THRESHOLD
}

// Returns the palette index closest to the color in the current mode
// converted to the proper string format
func ColorString(mode RenderMode, px Pixel) string {
	if px.alpha < TRANSPARENCY_THRESHOLD {
		return "" // sentinel for transparent colors
	}
	if mode == term24bit {
		r, g, b := px.color.RGB255()
		return fmt.Sprintf("%d;%d;%d", r, g, b)
	}
	result := 0
	dist := ColorDistance(mode, px.color, colors[mode][0])
	for i := 0; i < len(colors[mode]); i++ {
		d := ColorDistance(mode, px.color, colors[mode][i])
		if d < dist {
			dist = d
			result = i
		}
	}
	return strconv.Itoa(result)
}

//
// Escape squences
//

func StartFGColor(mode RenderMode) string {
	if mode == irc {
		return "\x03"
	} else {
		// Foreground color selector is 38;
		if mode == term24bit {
			return "\x1B[38;2;"
		} else {
			return "\x1B[38;5;"
		}
	}
}

func StartBGColor(mode RenderMode) string {
	if mode == irc {
		return ","
	} else {
		// Background color selector is 48;
		if mode == term24bit {
			return "\x1B[48;2;"
		} else { // 256
			return "\x1B[48;5;"
		}
	}
}

func EndColor(mode RenderMode) string {
	if mode != irc {
		// ANSI escape sequences are terminated by 'm'
		return "m"
	}
	return ""
}

func Clear(mode RenderMode) string {
	if mode != irc {
		return "\x1B[0m"
	}
	return ""
}

func Render(mode RenderMode, use_spaces bool, colors [][]Pixel) string {
	ix, iy, w, h := 0, 0, len(colors[0]), len(colors)
	var buffer bytes.Buffer
	ch := ""
	step := 2
	if use_spaces {
		// Two spaces to keep aspect ratio
		ch = "  "
		step = 1
	}

	for y := iy; y < h; y += step {
		var prev_fg_col string
		var prev_bg_col string
		for x := ix; x < w; x++ {
			next_fg_col := ColorString(mode, colors[y][x])
			next_bg_col := ""
			// Process two vertical pixels per column unless we're printing spaces
			if !use_spaces {
				ch = "▀"
				if y+1 < h {
					next_bg_col = ColorString(mode, colors[y+1][x])
				}
				if next_fg_col == "" {
					if next_bg_col == "" {
						ch = " "
					} else {
						ch = "▄"
					}
					next_fg_col = next_bg_col
					next_bg_col = ""
				}
			} else {
				next_bg_col = next_fg_col
				next_fg_col = ""
			}

			if next_fg_col == "" && next_bg_col == "" {
				if prev_fg_col != "" || prev_bg_col != "" {
					buffer.WriteString(Clear(mode))
				}
				prev_fg_col = ""
				prev_bg_col = ""
			} else if prev_fg_col != next_fg_col || prev_bg_col != next_bg_col {
				if mode == irc || prev_fg_col != next_fg_col {
					if mode == irc && next_fg_col == "" {
						next_fg_col = "0"
					}
					buffer.WriteString(StartFGColor(mode))
					buffer.WriteString(next_fg_col)
					buffer.WriteString(EndColor(mode))
					prev_fg_col = next_fg_col
				}
				if next_bg_col != "" && prev_bg_col != next_bg_col {
					buffer.WriteString(StartBGColor(mode))
					buffer.WriteString(next_bg_col)
					buffer.WriteString(EndColor(mode))
					prev_bg_col = next_bg_col
				}
			}
			buffer.WriteString(ch)
		}
		buffer.WriteString(Clear(mode))
		buffer.WriteString("\n")
	}
	return buffer.String()
}
