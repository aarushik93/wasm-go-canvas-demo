//Wasming
// compile: 89
package main

import (
	"fmt"
	"math"
	"math/rand"
	"syscall/js"

	"github.com/aarushik93/wasm-example/utils"
)

var (
	width         float64
	height        float64
	sctx          js.Value
	wctx          js.Value
	window        js.Value
	doc           js.Value
	shapesCanvas js.Value
	wordCanvas    js.Value
	words         = []string{"GOPHER", "RUSHI", "WASM"}
	quality       = 10
	counter       = 0
	drawNext      = true

	renderer       js.Func
	rendererReplay js.Func

	shapes     []*Shape
	possiblePos []*Coords

	cHeight = 800
	cWidth = 370
)

func setUp() {
	window = js.Global()
	doc = window.Get("document")


	rendererReplay := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		replay()
		return nil
	})

	possiblePos = make([]*Coords, 0, int(height*width)*2)
	shapes = make([]*Shape, 0, int(height*width))

	js.Global().Set("replay", rendererReplay)
	
	shapesCanvas = doc.Call("getElementById", "sCanvas")
	sctx = shapesCanvas.Call("getContext", "2d")

	wordCanvas = doc.Call("getElementById", "wCanvas")
	wctx = wordCanvas.Call("getContext", "2d")

	width = window.Get("innerWidth").Float()
	height = window.Get("innerHeight").Float()

	shapesCanvas.Set("height", height)
	shapesCanvas.Set("width", width)

	wordCanvas.Set("height", height)
	wordCanvas.Set("width", width)

	shapesCanvas.Call("setAttribute", "width", width)
	shapesCanvas.Call("setAttribute", "height", height)

	wordCanvas.Call("setAttribute", "width", width)
	wordCanvas.Call("setAttribute", "height", height)

	sctx.Set("fillStyle", "white")

	sctx.Set("strokeStyle", "violet")
	sctx.Set("lineWidth", 5)


	wctx.Set("textAlign", "center")
	wctx.Set("textBaseLine", "middle")

}

func main() {
	done := make(chan struct{}, 0)
	// Init Canvas stuff
	setUp()

	defer rendererReplay.Release()
	<-done

}

func newShape() {
	var newPos *Coords

	for newPos == nil {
		newPos = possiblePos[rand.Intn(len(possiblePos))]
	}

	draw := true

	for _, s := range shapes {
		if utils.Distance(s.c.x, s.c.y, newPos.x, newPos.y) <s.h{
			draw = false
			break
		}
	}

	if draw {
		shapes = append(shapes, &Shape{
			c:         newPos,
			h:         1,
			ctx:       sctx,
			isGrowing: true,
		})
	}
}

func loop() {
	wctx.Call("clearRect", 0, 0, width, height)

	for i := 0; i < quality; i++ {
		newShape()
	}
	for _, s := range shapes {
		for _, s1 := range shapes {
			if s.c.x != s1.c.x || s.c.y != s1.c.y || s.h != s.h {
				if utils.Distance(s.c.x, s.c.y, s1.c.x, s1.c.y) < s.h+s1.h {
					s.isGrowing = false
					break
				}
			}
		}
		s.Show()
		s.Grow()

	}

	window.Call("requestAnimationFrame", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		loop()
		return nil
	}))
}

func replay() {

	drawNext = true

	go wctx.Call("clearRect", 0, 0, width, height)
	go sctx.Call("clearRect", 0, 0, width, height)

	word := words[rand.Intn(len(words))]

	wctx.Set("font", fmt.Sprintf("normal 900 %+vpx sans-serif", width/float64(len(word))+2))
	wctx.Call("fillText", word, width/2, height/2)

	go func() {
		shapes = shapes[:0]
		possiblePos = possiblePos[:0]
		for i := 0; i < wordCanvas.Get("width").Int(); i += 5 {
			for g := 0; g < wordCanvas.Get("height").Int(); g += 5 {
				img := wctx.Call("getImageData", i, g, 1, 1)
				imgData := img.Get("data")
				for j := 0; j < imgData.Length(); j++ {
					if imgData.Index(0).Int() > 0 || imgData.Index(1).Int() > 0 || imgData.Index(2).Int() > 0 || imgData.Index(3).Int() > 0 {
						possiblePos = append(possiblePos, &Coords{
							x: float64(i),
							y: float64(g),
						})

					}
				}

			}
		}

		wctx.Call("clearRect", 0, 0, width, height)
		loop()
	}()


}

type Coords struct {
	x float64
	y float64
}

type Shape struct {
	c         *Coords
	h         float64
	ctx       js.Value
	isGrowing bool
	r         int
	g         int
	b         int
	a         int
}

func (s *Shape) Grow() {
	if s.isGrowing {
		s.h += 1
	}
}

func (s *Shape) Show() {
	s.ctx.Call("fill")
	s.ctx.Call("beginPath")
	s.ctx.Call("arc", s.c.x, s.c.y, s.h, 0, 2*math.Pi)
	s.ctx.Call("stroke")
	s.ctx.Call("closePath")
}
