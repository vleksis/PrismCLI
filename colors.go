package main

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

const (
	// ANSI tags to render text
	resetCode  = "\033[0m"
	fgColorFmt = "\033[38;2;%d;%d;%dm"
	bgColorFmt = "\033[48;2;%d;%d;%dm"

	// RGB input regular expression
	rgbRegexp = `^[rR][gG][bB]\((\d+), *(\d+), *(\d+)\)$`
)

type RGB struct {
	R, G, B uint8
}

type ConversionError struct{}

func (e *ConversionError) Error() string {
	return "invalid RGB format"
}

// transform "rgb(0,1,2)" into RGB{0, 1, 2}
func stringToRGB(s string) (RGB, error) {
	re := regexp.MustCompile(rgbRegexp)
	matches := re.FindStringSubmatch(s)
	if len(matches) != 4 {
		return RGB{}, &ConversionError{}
	}

	r, err := strconv.ParseUint(matches[1], 10, 8)
	if err != nil {
		return RGB{0, 0, 0}, &ConversionError{}
	}
	g, err := strconv.ParseUint(matches[2], 10, 8)
	if err != nil {
		return RGB{0, 0, 0}, &ConversionError{}
	}
	b, err := strconv.ParseUint(matches[3], 10, 8)
	if err != nil {
		return RGB{0, 0, 0}, &ConversionError{}
	}
	return RGB{uint8(r), uint8(g), uint8(b)}, nil
}

func stringToForegrounder(s string) (func(int) string, error) {
	if f, ok := ColorToForegrounder[s]; ok {
		return f, nil
	}
	rgb, err := stringToRGB(s)
	if err != nil {
		return nil, err
	}
	return NewMonoForegrounder(rgb), nil
}

func stringToBackgrounder(s string) (func(int) string, error) {
	if f, ok := ColorToBackgrounder[s]; ok {
		return f, nil
	}
	rgb, err := stringToRGB(s)
	if err != nil {
		return nil, err
	}
	return NewMonoBackgrounder(rgb), nil
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

func BoldStyler(int) string {
	return "\033[1m"
}

func ItalicStyler(int) string {
	return "\033[3m"
}

func UnderlineStyler(int) string {
	return "\033[4m"
}

func BlinkStyler(int) string {
	return "\033[5m"
}

func StrikethroughStyler(int) string {
	return "\033[9"
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

func (w *CustomWriter) SetStylers(stylers ...func(int) string) *CustomWriter {
	var sb strings.Builder
	for _, fn := range stylers {
		sb.WriteString(fn(w.index))
	}
	w.styler = func(index int) string {
		return sb.String()
	}
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
