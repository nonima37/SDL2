package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight int = 800, 600

type color struct {
	r, g, b byte
}

type pos struct {
	x, y float32
}

type ball struct {
	pos
	radius int
	xv     float32
	yv     float32
	color  color
}

func (ball *ball) draw(pixels []byte) {
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x)+x, int(ball.y)+y, color{255, 255, 255}, pixels)
			}
		}
	}
}

// update pos and collision detection
func (ball *ball) update(paddle1, paddle2 *paddle) {
	ball.x += ball.xv
	ball.y += ball.yv

	if int(ball.y)-ball.radius < 0 || int(ball.y)+ball.radius > winHeight {
		ball.yv = -ball.yv
	}

	if ball.x < paddle1.x || ball.x > paddle2.x {
		ball.x = 300
		ball.x = 300
	}

	if ball.x-float32(ball.radius) < float32(paddle1.x)+float32(paddle1.w/2) {
		if int(ball.y) >= int(paddle1.y)-paddle1.h/2 && int(ball.y) <= int(paddle1.y)+paddle1.h/2 {
			ball.xv = -ball.xv
		}
	}

	if ball.x+float32(ball.radius) > float32(paddle2.x)-float32(paddle2.w/2) {
		if int(ball.y) >= int(paddle2.y)-paddle2.h/2 && int(ball.y) <= int(paddle2.y)+paddle1.h/2 {
			ball.xv = -ball.xv
		}
	}
}

type paddle struct {
	pos
	w     int
	h     int
	color color
}

// x and y are the middle points
func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x) - paddle.w/2
	startY := int(paddle.y) - paddle.h/2

	for y := 0; y < paddle.h; y++ {
		for x := 0; x < paddle.w; x++ {
			setPixel(startX+x, startY+y, paddle.color, pixels)
		}
	}
}

func (paddle *paddle) update(keyState []uint8, pixels []byte) {
	if keyState[sdl.SCANCODE_UP] != 0 && paddle.y >= float32(paddle.h/2)+10 {
		paddle.y -= 10
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 && paddle.y+float32(paddle.h)/2+10 <= float32(winHeight) {
		paddle.y += 10
	}
}

func (paddle *paddle) aiUpdate(ball *ball) {
	if int(ball.y) >= paddle.h/2 && int(ball.y) <= winHeight-paddle.h/2 {
		paddle.y = ball.y
	}
}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y*winWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

func clear(pixels []byte) {
	// this can be done more efficiently
	for i := range pixels {
		pixels[i] = 0
	}
}

func main() {
	window, err := sdl.CreateWindow("window", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		fmt.Println(err)
		return
	}

	pixels := make([]byte, winWidth*winHeight*4)

	player1 := paddle{pos: pos{100, 100}, w: 20, h: 100, color: color{255, 255, 255}}
	player2 := paddle{pos: pos{700, 100}, w: 20, h: 100, color: color{255, 255, 255}}
	ball := ball{pos: pos{x: 50, y: 70}, radius: 10, xv: 10, yv: 10, color: color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		clear(pixels)

		player1.update(keyState, pixels)
		player2.aiUpdate(&ball)
		ball.update(&player1, &player2)

		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()
		sdl.Delay(16)
	}

	tex.Destroy()
	window.Destroy()
	renderer.Destroy()

}
