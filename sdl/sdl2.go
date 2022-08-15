// ToDo
// rewrite the way paddle is being drawn to make it easier to rotate []
// rotating paddle []
// collision detection needs some improvements []
// Imperfect AI []
package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight int = 800, 600

type gameState int

const (
	start gameState = iota
	play
)

var state = start

var nums = [][]byte{
	{
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
	},
	{
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
	},
	{
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
	},
}

func drawScore(pos pos, size, scoreIndex int, pixels []byte) {
	startX := int(pos.x)
	startY := int(pos.y)

	for i := 0; i < len(nums[scoreIndex]); i++ {
		if nums[scoreIndex][i] == 1 {
			for y := startY; y < startY+size; y++ {
				for x := startX; x < startX+size; x++ {
					setPixel(x, y, color{255, 255, 255}, pixels)
				}
			}
		} else if scoreIndex > 0 {
			for y := startY; y < startY+size; y++ {
				for x := startX; x < startX+size; x++ {
					setPixel(x, y, color{0, 0, 0}, pixels)
				}
			}
		}

		startX += size
		if (i+1)%3 == 0 {
			startY += size
			startX -= size * 3
		}
	}
}

func lerp(a, b, pct float32) float32 {
	return a + pct*(b-a)
}

type color struct {
	r, g, b byte
}

type pos struct {
	x, y float32
}

type ball struct {
	pos
	radius float32
	xv     float32
	yv     float32
	color  color
}

func (ball *ball) drawOrClear(clear bool, pixels []byte) {
	var curColor color
	if clear {
		curColor = color{0, 0, 0}
	} else {
		curColor = ball.color
	}
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x+x), int(ball.y+y), curColor, pixels)
			}
		}
	}
}

func getCenter() pos {
	return pos{float32(winWidth) / 2, float32(winHeight) / 2}
}

// update pos and collision detection
func (ball *ball) update(paddle1, paddle2 *paddle, pixels []byte) {
	ball.x += ball.xv
	ball.y += ball.yv

	if ball.y-ball.radius <= 0 || int(ball.y+ball.radius) > winHeight {
		ball.yv = -ball.yv
	}

	if ball.x < 0 {
		//paddle2.score++
		//ball.pos = getCenter()
		ball.pos = pos{140, 140}
		state = start
	} else if int(ball.x) > winWidth {
		//paddle1.score++
		//ball.pos = getCenter()
		ball.pos = pos{120, 140}
		state = start
	}

	// outer side collision detection
	if ball.x-ball.radius == paddle1.x+paddle1.w/2 {
		if ball.y >= paddle1.y-paddle1.h/2 && ball.y <= paddle1.y+paddle1.h/2 {
			ball.xv = -ball.xv
			// for when the ball gets stuck on the paddle and keeps wiggling
			ball.x = paddle1.x + paddle1.w/2.0 + ball.radius
		}
	} else if ball.x-ball.radius < paddle1.x+paddle1.w/2 && ball.x+ball.radius > paddle1.x-paddle1.w/2 {
		// bottom/top collision detection
		if ball.yv < 0 && ball.y-ball.radius <= paddle1.y+paddle1.h/2 {
			ball.yv = -ball.yv
		}
	}

	if ball.x+ball.radius > paddle2.x-paddle2.w/2 {
		if ball.y >= paddle2.y-paddle2.h/2 && ball.y <= paddle2.y+paddle1.h/2 {
			ball.xv = -ball.xv
			ball.x = paddle2.x - paddle2.w/2.0 - ball.radius
		}
	}
}

type paddle struct {
	pos
	w     float32
	h     float32
	speed float32
	color color
	score int
}

// x and y are the middle points
func (paddle *paddle) drawOrClear(clear bool, pixels []byte) {
	startX := int(paddle.x - paddle.w/2)
	startY := int(paddle.y - paddle.h/2)

	var curColor color
	if clear {
		curColor = color{0, 0, 0}
	} else {
		curColor = paddle.color
	}

	for y := startY; y < startY+int(paddle.h); y++ {
		for x := startX; x < startX+int(paddle.w); x++ {
			setPixel(x, y, curColor, pixels)
		}
	}
	numX := lerp(paddle.x, getCenter().x, 0.2)
	drawScore(pos{numX, 35}, 10, paddle.score, pixels)
}

func (paddle *paddle) update(keyState []uint8) {
	if keyState[sdl.SCANCODE_UP] != 0 && paddle.y >= paddle.h/2+paddle.speed {
		paddle.y -= paddle.speed
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 && paddle.y+paddle.h/2+paddle.speed <= float32(winHeight) {
		paddle.y += paddle.speed
	}
}

func (paddle *paddle) aiUpdate(ball *ball) {
	if ball.y >= paddle.h/2 && ball.y <= float32(winHeight)-paddle.h/2 {
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

	// each pixel takes up 4 bytes
	// 4th one is for alpha(?), we don't use it tho
	pixels := make([]byte, winWidth*winHeight*4)

	player1 := paddle{pos: pos{100, 100}, w: 20, h: 100, speed: 10, color: color{255, 255, 255}, score: 0}
	player2 := paddle{pos: pos{700, 100}, w: 20, h: 100, speed: 10, color: color{255, 255, 255}, score: 0}
	ball := ball{pos: pos{x: 100, y: 170}, radius: 20, xv: 5, yv: -5, color: color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		player1.drawOrClear(true, pixels)
		player2.drawOrClear(true, pixels)
		ball.drawOrClear(true, pixels)

		player1.update(keyState)
		player2.aiUpdate(&ball)

		if state == play {
			player1.update(keyState)
			player2.aiUpdate(&ball)
			ball.update(&player1, &player2, pixels)
		} else if state == start {
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				if player1.score == 3 || player2.score == 3 {
					player1.score = 0
					player2.score = 0
				}
				state = play
			}
		}

		player1.drawOrClear(false, pixels)
		player2.drawOrClear(false, pixels)
		ball.drawOrClear(false, pixels)
		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		sdl.Delay(16)
	}

	tex.Destroy()
	window.Destroy()
	renderer.Destroy()
}
