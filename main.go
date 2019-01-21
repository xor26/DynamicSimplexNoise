package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight int = 1200, 800

type color struct {
	r, g, b byte
}


func lerp(b1, b2 byte, percent float32) byte {
	return byte(float32(b1) + percent*(float32(b2)-float32(b1)))
}

func colorLerp(c1, c2 color, percent float32) color{
	red := lerp(c1.r, c2.r, percent)
	green := lerp(c1.g, c2.g, percent)
	blue := lerp(c1.b, c2.b, percent)
	return color{red, green, blue}
}

func getGradient(c1 color, c2 color) []color {
	result := make([]color, 256)
	for i := range result{
		percent := float32(i)/float32(255)
		result[i] = colorLerp(c1, c2, percent)
	}

	return result
}

func clamp(min, max, v int) int {
	if v < min {
		return min
	}

	if v > max {
		return max
	}

	return v
}

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Pong", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tex.Destroy()
	pixels := make([]byte, winWidth*winHeight*4)

	frequency := float32(10)
	lacunarity := float32(3)
	gain := float32(2)
	octaves := 2
	noise := make([]float32, 1200*800)
	gradient := getGradient(color{255, 0,0 }, color{0,255,255})

	keyState := sdl.GetKeyboardState()
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		if  keyState[sdl.SCANCODE_F] != 0 {
			frequency += 5
		}
		if  keyState[sdl.SCANCODE_D] != 0 {
			if frequency > 0 {
				frequency -= 5
			}
		}

		if  keyState[sdl.SCANCODE_L] != 0 {
			lacunarity += 3
		}
		if  keyState[sdl.SCANCODE_K] != 0 {
			lacunarity -= 3
		}

		if  keyState[sdl.SCANCODE_G] != 0 {
			gain += 1
		}

		if  keyState[sdl.SCANCODE_F] != 0 {
			if gain > 0 {
				gain -= 1
			}
		}

		if  keyState[sdl.SCANCODE_O] != 0 {
			octaves += 1
		}

		if  keyState[sdl.SCANCODE_I] != 0 {
			if octaves > 0 {
				octaves -= 1
			}
		}
		fmt.Println("Frequency: ", frequency, " Octaves: ", octaves, " Gain", gain, "Lacunarity", lacunarity)
		regenerateNoisePermutationTable()

		min := float32(9999.0)
		max := float32(-9999.0)
		i := 0
		for x := 0; x < 800; x++ {
			for y := 0; y < 1200; y++ {
				//noise[i] = makeNoise(float32(x), float32(y),frequency, lacunarity, gain, octaves)
				noise[i] = makeTurbulentNoise(float32(x), float32(y),frequency, lacunarity, gain, octaves)
				if min > noise[i] {
					min = noise[i]
				}
				if max < noise[i] {
					max = noise[i]
				}
				i++
			}
		}
		scale := 255.0 / (max - min)
		offset := min * scale
		for i := range noise {
			byteNoise := byte(noise[i]*scale + offset)
			pixelColor := gradient[clamp(0, 255, int(byteNoise) )]

			pixels[4*i] = pixelColor.r
			pixels[4*i+1] = pixelColor.g
			pixels[4*i+2] = pixelColor.b
		}

		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()
		sdl.Delay(60)

	}
}
