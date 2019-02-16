package main

import (
	"fmt"
	"image/png"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/sabith-th/games_with_go/noise"
	"github.com/sabith-th/games_with_go/vector3"
	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight, winDepth int = 800, 600, 100

type mouseState struct {
	leftButton  bool
	rightButton bool
	x, y        int
}

func getMouseState() mouseState {
	mouseX, mouseY, mouseButtonState := sdl.GetMouseState()
	leftButton := mouseButtonState & sdl.ButtonLMask()
	rightButton := mouseButtonState & sdl.ButtonRMask()
	var result mouseState
	result.x = int(mouseX)
	result.y = int(mouseY)
	result.leftButton = !(leftButton == 0)
	result.rightButton = !(rightButton == 0)
	return result
}

type balloon struct {
	tex  *sdl.Texture
	pos  vector3.Vector3
	dir  vector3.Vector3
	w, h int
}

type balloonArray []*balloon

func (balloons balloonArray) Len() int {
	return len(balloons)
}

func (balloons balloonArray) Swap(i, j int) {
	balloons[i], balloons[j] = balloons[j], balloons[i]
}

func (balloons balloonArray) Less(i, j int) bool {
	diff := balloons[i].pos.Z - balloons[j].pos.Z
	return diff < -0.5
}

func (balloon *balloon) update(elapsedTime float32) {
	p := vector3.Add(balloon.pos, vector3.Mult(balloon.dir, elapsedTime))

	if p.X < 0 || p.X > float32(winWidth) {
		balloon.dir.X = -balloon.dir.X
	}

	if p.Y < 0 || p.Y > float32(winHeight) {
		balloon.dir.Y = -balloon.dir.Y
	}

	if p.Z < 0 || p.Z > float32(winDepth) {
		balloon.dir.Z = -balloon.dir.Z
	}

	balloon.pos = vector3.Add(balloon.pos, vector3.Mult(balloon.dir, elapsedTime))
}

func (balloon *balloon) draw(renderer *sdl.Renderer) {
	scale := (balloon.pos.Z/200 + 1) / 2
	newW := int32(float32(balloon.w) * scale)
	newH := int32(float32(balloon.h) * scale)
	x := int32(balloon.pos.X - float32(newW)/2)
	y := int32(balloon.pos.Y - float32(newH)/2)
	rect := &sdl.Rect{X: x, Y: y, W: newW, H: newH}
	renderer.Copy(balloon.tex, nil, rect)
}

type rgba struct {
	r, g, b byte
}

func lerp(b1, b2 byte, pct float32) byte {
	return byte(float32(b1) + pct*(float32(b2)-float32(b1)))
}

func rgbalerp(c1, c2 rgba, pct float32) rgba {
	return rgba{lerp(c1.r, c2.r, pct), lerp(c1.g, c2.g, pct), lerp(c1.b, c2.b, pct)}
}

func getGradient(c1, c2 rgba) []rgba {
	result := make([]rgba, 256)
	for i := range result {
		pct := float32(i) / float32(255)
		result[i] = rgbalerp(c1, c2, pct)
	}
	return result
}

func clamp(min, max, v int) int {
	if v < min {
		v = min
	} else if v > max {
		v = max
	}
	return v
}

func rescaleAndDraw(noise []float32, min, max float32, gradient []rgba, w, h int) []byte {
	result := make([]byte, w*h*4)
	scale := 255.0 / (max - min)
	offset := min * scale

	for i := range noise {
		noise[i] = noise[i]*scale - offset
		c := gradient[clamp(0, 255, int(noise[i]))]
		p := i * 4
		result[p] = c.r
		result[p+1] = c.g
		result[p+2] = c.b
	}
	return result
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixel(x, y int, c rgba, pixels []byte) {
	index := (y*winWidth + x) * 4
	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

func pixelsToTexture(renderer *sdl.Renderer, pixels []byte, w, h int) *sdl.Texture {
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(w), int32(h))
	if err != nil {
		panic(err)
	}
	tex.Update(nil, pixels, w*4)
	return tex
}

func loadBalloons(renderer *sdl.Renderer, numBalloons int) []*balloon {
	balloonStrs := []string{"images/balloon_red.png", "images/balloon_blue.png", "images/balloon_green.png"}
	balloonTextures := make([]*sdl.Texture, len(balloonStrs))

	for i, bstr := range balloonStrs {
		infile, err := os.Open(bstr)
		if err != nil {
			panic(err)
		}
		defer infile.Close()

		img, err := png.Decode(infile)
		if err != nil {
			panic(err)
		}

		w := img.Bounds().Max.X
		h := img.Bounds().Max.Y

		balloonPixels := make([]byte, w*h*4)
		bIndex := 0
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				balloonPixels[bIndex] = byte(r / 256)
				bIndex++
				balloonPixels[bIndex] = byte(g / 256)
				bIndex++
				balloonPixels[bIndex] = byte(b / 256)
				bIndex++
				balloonPixels[bIndex] = byte(a / 256)
				bIndex++
			}
		}
		tex := pixelsToTexture(renderer, balloonPixels, w, h)
		err = tex.SetBlendMode(sdl.BLENDMODE_BLEND)
		if err != nil {
			panic(err)
		}
		balloonTextures[i] = tex
	}
	balloons := make([]*balloon, numBalloons)
	for i := range balloons {
		tex := balloonTextures[i%3]
		pos := vector3.Vector3{X: rand.Float32() * float32(winWidth), Y: rand.Float32() * float32(winHeight), Z: rand.Float32() * float32(winDepth)}
		dir := vector3.Vector3{X: rand.Float32() * 0.5, Y: rand.Float32() * 0.5, Z: rand.Float32() * 0.5}
		_, _, w, h, err := tex.Query()
		if err != nil {
			panic(err)
		}
		balloons[i] = &balloon{tex, pos, dir, int(w), int(h)}
	}
	return balloons
}

func main() {

	window, err := sdl.CreateWindow("Balloons", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
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
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	cloudNoise, min, max := noise.MakeNoise(noise.FBM, 0.009, 0.5, 3, 3, winWidth, winHeight)
	cloudGradient := getGradient(rgba{0, 0, 255}, rgba{255, 255, 255})
	cloudPixels := rescaleAndDraw(cloudNoise, min, max, cloudGradient, winWidth, winHeight)
	cloudTexture := pixelsToTexture(renderer, cloudPixels, winWidth, winHeight)

	balloons := loadBalloons(renderer, 20)
	var elapsedTime float32
	currentMouseState := getMouseState()
	prevMouseState := currentMouseState

	for {
		frameStart := time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		currentMouseState = getMouseState()

		if !currentMouseState.leftButton && prevMouseState.leftButton {
			fmt.Println("Left Click")
		}
		renderer.Copy(cloudTexture, nil, nil)

		for _, balloon := range balloons {
			balloon.update(elapsedTime)
		}
		sort.Stable(balloonArray(balloons))
		for _, balloon := range balloons {
			balloon.draw(renderer)
		}

		renderer.Present()

		elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		if elapsedTime < 5 {
			sdl.Delay(5 - uint32(elapsedTime))
			elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		}

		prevMouseState = currentMouseState
	}

}
