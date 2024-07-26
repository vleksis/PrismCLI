package main

import (
	"bufio"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := NewCustomWriter(os.Stdout)
	writer.SetForegrounder(NewMonoForegrounder(ColorToRGB["red"]))

	for {
		input, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		writer.Write(input)
	}

}
