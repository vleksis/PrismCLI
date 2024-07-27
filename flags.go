package main

import (
	"flag"
	"io"
	"log"
	"os"
)

const (
	// presets
	defaultPreset   = "default"
	darkModePreset  = "dark"
	lightModePreset = "light"
	rainbowPreset   = "rainbow"
	customPreset    = "custom"
)

var (
	// flags

	preset = flag.String("preset", defaultPreset,
		"Text style preset (default, dark, light, rainbow)")
	fgColor = flag.String("fg", "",
		"Foreground color (e.g., red, green, blue or RGB(255,0,0))")
	bgColor = flag.String("bg", "",
		"Background color (e.g., red, green, blue or RGB(0,0,0))")

	bold = flag.Bool("bold", false,
		"Apply bold style")
	italic = flag.Bool("italic", false,
		"Apply italic style")
	underline = flag.Bool("underline", false,
		"Apply underline style")
	blink = flag.Bool("blink", false,
		"Apply blink style")
	strikethrough = flag.Bool("strikethrough", false,
		"Apply strikethrough style")

	inputFile = flag.String("input", "",
		"Input file containing text to colorize")
	outputFile = flag.String("output", "",
		"Output file to save the colorized text")
	configFile = flag.String("config", "",
		"Configuration file for default settings")

	helpFlag = flag.Bool("h", false, "Show help message")

	// helper global variables
	openedFiles  = make([]*os.File, 0)
	stylersFlags = []struct {
		flag   *bool
		styler func(int) string
	}{
		{bold, BoldStyler},
		{italic, ItalicStyler},
		{underline, UnderlineStyler},
		{blink, BlinkStyler},
		{strikethrough, StrikethroughStyler},
	}

	writer CustomWriter = NewCustomWriter(nil)
	reader io.Reader

	// messages
	helpText = `
Usage: mycli [options] <text>
A CLI tool for colorizing text with styles.

Options:
  -preset string
    	Text style preset (default, dark, light, rainbow) (default "default")
  -fg string
    	Foreground color (e.g., red, green, blue or RGB(255,0,0)) (default "white")
  -bg string
    	Background color (e.g., red, green, blue or RGB(0,0,0)) (default "black")
  -bold
    	Apply bold style
  -italic
    	Apply italic style
  -underline
    	Apply underline style
  -strikethrough
    	Apply strikethrough style
  -input string
    	Input file containing text to colorize
  -output string
    	Output file to save the colorized text
  -config string
    	Configuration file for default settings
  -h, --help
    	Show this help message and exit

Examples:
  mycli -fg red -bg blue -bold -underline "Hello, World!"
  mycli -preset dark -input text.txt -output colored_text.txt
`
)

func init() {
	flag.BoolVar(helpFlag, "help", false, "Show help message")
}

func processFlags() {
	flag.Parse()

	if err := validateFlags(); err != nil {
		log.Fatal(err)
	}

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}
	processInputFlag()
	processOutputFlag()
	processForegrounderFlag()
	processBackgroundFlag()
	processStylerFlags()
}

func validateFlags() error {
	// todo: implement validation
	return nil
}

func processInputFlag() error {
	if *inputFile == "" {
		reader = os.Stdin
		return nil
	}
	f, err := os.Open(*inputFile)
	if err == nil {
		reader = f
		openedFiles = append(openedFiles, f)
	}
	return err

}

func processOutputFlag() error {
	if *outputFile == "" {
		writer.base = os.Stdout
		return nil
	}
	f, err := os.Create(*outputFile)
	if err == nil {
		writer.base = f
		openedFiles = append(openedFiles, f)
	}
	return err
}

func processForegrounderFlag() error {
	foregrounder, err := stringToForegrounder(*fgColor)
	if err != nil {
		return err
	}
	writer.SetForegrounder(foregrounder)
	return nil
}

func processBackgroundFlag() error {
	backgrounder, err := stringToBackgrounder(*bgColor)
	if err != nil {
		return err
	}
	writer.SetBackgrounder(backgrounder)
	return nil
}

func processStylerFlags() {
	var pool []func(int) string
	for _, f := range stylersFlags {
		if *f.flag {
			pool = append(pool, f.styler)
		}
	}
	writer.SetStylers(pool...)
}

func closeAllFiles() {
	for _, f := range openedFiles {
		if err := f.Close(); err != nil {
			log.Println("Error closing file:", err)
		}
	}
}
