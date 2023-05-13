package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
)

const imagesDir = "images/"

// Width and height of the image
const size = 2 << 10

type Rule func(left, center, right bool) bool

func generateRule(name int) (Rule, error) {
	if name > 255 || name < 0 {
		return nil, errors.New("rule name must be between 0 and 255 inclusive")
	}

	return func(left, center, right bool) bool {
		if !left && !center && !right {
			return name%2 == 1
		} else if !left && !center && right {
			return (name/2)%2 == 1
		} else if !left && center && !right {
			return (name/4)%2 == 1
		} else if !left && center && right {
			return (name/8)%2 == 1
		} else if left && !center && !right {
			return (name/16)%2 == 1
		} else if left && !center && right {
			return (name/32)%2 == 1
		} else if left && center && !right {
			return (name/64)%2 == 1
		} else if left && center && right {
			return (name/128)%2 == 1
		}

		// should never get here
		return false
	}, nil

}

type simResults [][]bool

func simulateRule(rule Rule) simResults {
	data := make([][]bool, size)
	for i := range data {
		data[i] = make([]bool, size)
	}
	data[0][size/2] = true

	for y := 1; y < size; y++ {
		for x := 0; x < size; x++ {
			l := data[y-1][max(0, x-1)]
			c := data[y-1][x]
			r := data[y-1][min(size-1, x+1)]
			data[y][x] = rule(l, c, r)
		}
	}

	return data
}

func createImage(data simResults, offset int, filename string) {
	img := image.NewRGBA(image.Rect(0, 0, size, size-offset))
	black := color.RGBA{0, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}

	for y, row := range data {
		if y < offset {
			// skip this data
			continue
		}
		for x, cell := range row {
			if cell {
				img.Set(x, y, black)
			} else {
				img.Set(x, y, white)
			}
		}
	}

	f, err := os.Create(filepath.Join(imagesDir, filename))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}

func main() {

	for ruleName := 0; ruleName < 256; ruleName++ {
		fmt.Println("Starting simulation for rule" + strconv.Itoa(ruleName))
		rule, err := generateRule(ruleName)
		if err != nil {
			panic(err)
		}
		data := simulateRule(rule)

		var offset int = 0

		createImage(data, offset, "rule"+strconv.Itoa(ruleName)+".png")
	}

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
