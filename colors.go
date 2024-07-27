package main

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"strings"
)

const (
	// ANSI tags to render text
	resetCode  = "\033[0m"
	fgColorFmt = "\033[38;2;%d;%d;%dm"
	bgColorFmt = "\033[48;2;%d;%d;%dm"
)

type RGB struct {
	R, G, B uint8
}

func rgbToANSIForeground(code RGB) string {
	return fmt.Sprintf(fgColorFmt, code.R, code.G, code.B)
}

func rgbToANSIBackground(code RGB) string {
	return fmt.Sprintf(bgColorFmt, code.R, code.G, code.B)
}

func rainbowRGB(freq, shift float64) func(int) RGB {
	return func(i int) RGB {
		return RGB{
			uint8(math.Sin(freq*float64(i)+0*math.Pi/3+shift)*127 + 128),
			uint8(math.Sin(freq*float64(i)+2*math.Pi/3+shift)*127 + 128),
			uint8(math.Sin(freq*float64(i)+4*math.Pi/3+shift)*127 + 128),
		}
	}
}

func NewRainbowForegrounder(freq, shift float64) func(int) string {
	return func(index int) string {
		return rgbToANSIForeground(rainbowRGB(freq, shift)(index))
	}
}

func NewRainbowBackgrounder(freq, init float64) func(int) string {
	return func(index int) string {
		return rgbToANSIBackground(rainbowRGB(freq, init)(index))
	}
}

func RandomRGB() RGB {
	return RGB{
		uint8(rand.Int()),
		uint8(rand.Int()),
		uint8(rand.Int()),
	}
}

func NewChaosForegrounder() func(int) string {
	return func(index int) string {
		return rgbToANSIForeground(RandomRGB())
	}
}

func NewChaosBackgrounder() func(int) string {
	return func(index int) string {
		return rgbToANSIBackground(RandomRGB())
	}
}

func NewMonoForegrounder(code RGB) func(int) string {
	return func(int) string {
		return rgbToANSIForeground(code)
	}
}

func NewMonoBackgrounder(code RGB) func(int) string {
	return func(int) string {
		return rgbToANSIBackground(code)
	}
}

type CustomWriter struct {
	base         io.Writer
	foregrounder func(int) string
	backgrounder func(int) string
	styler       func(int) string
	index        int
}

func NewCustomWriter(w io.Writer) CustomWriter {
	return CustomWriter{
		base:         w,
		foregrounder: func(int) string { return "" },
		backgrounder: func(int) string { return "" },
		styler:       func(int) string { return "" },
		index:        0,
	}
}

func (w *CustomWriter) SetForegrounder(f func(int) string) *CustomWriter {
	w.foregrounder = f
	return w
}

func (w *CustomWriter) SetBackgrounder(f func(int) string) *CustomWriter {
	w.backgrounder = f
	return w
}

func (w *CustomWriter) SetStyler(f func(int) string) *CustomWriter {
	w.styler = f
	return w
}

func (w *CustomWriter) Write(p []byte) (n int, err error) {
	var formattedText strings.Builder
	for _, s := range p {
		if s == '\n' {
			formattedText.WriteByte(s)
			continue
		}

		formattedText.WriteString(w.foregrounder(w.index))
		formattedText.WriteString(w.backgrounder(w.index))
		formattedText.WriteString(w.styler(w.index))
		formattedText.WriteByte(s)
		formattedText.WriteString(resetCode)

		w.index++
	}
	return w.base.Write([]byte(formattedText.String()))
}
