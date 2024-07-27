package main

var ColorToRGB map[string]RGB = map[string]RGB{
	"white":      RGB{255, 255, 255},
	"black":      RGB{0, 0, 0},
	"red":        RGB{255, 0, 0},
	"green":      RGB{0, 255, 0},
	"blue":       RGB{0, 0, 255},
	"yellow":     RGB{255, 255, 0},
	"cyan":       RGB{0, 255, 255},
	"magenta":    RGB{255, 0, 255},
	"gray":       RGB{128, 128, 128},
	"orange":     RGB{255, 165, 0},
	"purple":     RGB{128, 0, 128},
	"brown":      RGB{165, 42, 42},
	"pink":       RGB{255, 192, 203},
	"lime":       RGB{50, 205, 50},
	"teal":       RGB{0, 128, 128},
	"olive":      RGB{128, 128, 0},
	"maroon":     RGB{128, 0, 0},
	"navy":       RGB{0, 0, 128},
	"beige":      RGB{245, 245, 220},
	"coral":      RGB{255, 127, 80},
	"salmon":     RGB{250, 128, 114},
	"turquoise":  RGB{64, 224, 208},
	"indigo":     RGB{75, 0, 130},
	"khaki":      RGB{240, 230, 140},
	"plum":       RGB{221, 160, 221},
	"gold":       RGB{255, 215, 0},
	"silver":     RGB{192, 192, 192},
	"periwinkle": RGB{204, 204, 255},
	"lavender":   RGB{230, 230, 250},
	"orchid":     RGB{218, 112, 214},
	"wheat":      RGB{245, 222, 179},
}

var ColorToForegrounder map[string]func(int) string
var ColorToBackgrounder map[string]func(int) string

func init() {
	ColorToForegrounder = make(map[string]func(int) string)
	ColorToBackgrounder = make(map[string]func(int) string)

	for color, rgb := range ColorToRGB {
		ColorToForegrounder[color] = NewMonoForegrounder(rgb)
		ColorToBackgrounder[color] = NewMonoBackgrounder(rgb)
	}

	ColorToForegrounder["rainbow"] = NewRainbowForegrounder(0.03, -1)
	ColorToBackgrounder["rainbow"] = NewRainbowBackgrounder(0.03, 1)

	ColorToForegrounder["chaos"] = NewChaosForegrounder()
	ColorToBackgrounder["chaos"] = NewChaosBackgrounder()

}
