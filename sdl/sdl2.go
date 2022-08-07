// ToDo
// More efficient screen update []
// Keeping the score [x]
// Imperfect AI []
package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight int = 800, 600

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

func drawScore(startx, starty, size, scoreindex int, pixels []byte) {
	for i := 0; i < len(nums[scoreindex]); i++ {
		if nums[scoreindex][i] == 1 {
			for y := starty; y < starty+size; y++ {
				for x := startx; x < startx+size; x++ {
					setPixel(x, y, color{255, 255, 255}, pixels)
				}
			}
		} else if scoreindex > 0 {
			for y := starty; y < starty+size; y++ {
				for x := startx; x < startx+size; x++ {
					setPixel(x, y, color{0, 0, 0}, pixels)
				}
			}
		}

		startx += size
		if (i+1)%3 == 0 {
			starty += size
			startx -= size * 3
		}

	}
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

func (ball *ball) draw(pixels []byte) {
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x+x), int(ball.y+y), color{255, 255, 255}, pixels)
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

	if ball.y-ball.radius < 0 || int(ball.y+ball.radius) > winHeight {
		ball.yv = -ball.yv
	}

	if ball.x < paddle1.x {
		paddle2.score++
		ball.pos = getCenter()
	}

	if ball.x-ball.radius < paddle1.x+paddle1.w/2 {
		if ball.y >= paddle1.y-paddle1.h/2 && ball.y <= paddle1.y+paddle1.h/2 {
			ball.xv = -ball.xv
		}
	}

	if ball.x+ball.radius > paddle2.x-paddle2.w/2 {
		if ball.y >= paddle2.y-paddle2.h/2 && ball.y <= paddle2.y+paddle1.h/2 {
			ball.xv = -ball.xv
		}
	}
}

func (ball *ball) clear(pixels []byte) {
	//sdasdas
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x+x), int(ball.y+y), color{0, 0, 0}, pixels)
			}
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
func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x - paddle.w/2)
	startY := int(paddle.y - paddle.h/2)

	for y := startY; y < startY+int(paddle.h); y++ {
		for x := startX; x < startX+int(paddle.w); x++ {
			setPixel(x, y, paddle.color, pixels)
		}
	}
}

func (paddle *paddle) update(keyState []uint8) {
	if keyState[sdl.SCANCODE_UP] != 0 && paddle.y >= paddle.h/2+10 {
		paddle.y -= paddle.speed
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 && paddle.y+paddle.h/2+10 <= float32(winHeight) {
		paddle.y += paddle.speed
	}
}

// this is just the draw() function but let's leave it for now
func (paddle *paddle) clear(pixels []byte) {
	startX := int(paddle.x) - int(paddle.w/2)
	startY := int(paddle.y) - int(paddle.h/2)

	for y := 0; y < int(paddle.h); y++ {
		for x := 0; x < int(paddle.w); x++ {
			// why doesn't this work??
			// pixels[((startY+y)*winWidth + startX + x)] = 200
			// we'll use this for now
			setPixel(startX+x, startY+y, color{0, 0, 0}, pixels)
		}
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
	// 4th one is fo alpha(?), we don't use it
	pixels := make([]byte, winWidth*winHeight*4)

	player1 := paddle{pos: pos{100, 100}, w: 20, h: 100, speed: 20, color: color{255, 255, 255}, score: 0}
	player2 := paddle{pos: pos{700, 100}, w: 20, h: 100, speed: 20, color: color{255, 255, 255}, score: 0}
	ball := ball{pos: pos{x: 150, y: 70}, radius: 10, xv: 6, yv: 6, color: color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()

	for {
		drawScore(200, 30, 20, player2.score, pixels)

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		player1.clear(pixels)
		player2.clear(pixels)
		ball.clear(pixels)

		player1.update(keyState)
		player2.aiUpdate(&ball)
		ball.update(&player1, &player2, pixels)

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
