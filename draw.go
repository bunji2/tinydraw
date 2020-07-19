package main

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	defaultGridWidth  = 20
	defaultGridHeight = 20
	defaultMargin     = 2
)

var gridWidth, gridHeight, margin int

var black = color.Black
var white = color.White
var red = color.RGBA{0xFF, 0x00, 0x00, 0xFF}
var blue = color.RGBA{0x00, 0x00, 0xFF, 0xFF}
var green = color.RGBA{0x00, 0xFF, 0x00, 0xFF}
var yellow = color.RGBA{0xFF, 0xFF, 0x00, 0xFF}
var pink = color.RGBA{0xFF, 0x00, 0xFF, 0xFF}
var skyblue = color.RGBA{0x00, 0xFF, 0xFF, 0xFF}
var gray = color.RGBA{127, 127, 127, 255}

// Data は描画オブジェクトの型
type Data struct {
	gridVisible bool
	width       int
	height      int
	bgColor     color.Color
	fgColor     color.Color
	rgb         *image.RGBA
}

// NewColor 色のオブジェクトを生成する
func NewColor(r, g, b uint8) (c color.Color) {
	c = color.RGBA{r, g, b, 0xFF}
	return
}

// SetParams はグリッドの縦横のピクセル値を決める
func SetParams(params ...int) {
	var gw, gh, m int
	if len(params) > 1 {
		gw = params[0]
		gh = params[1]
		if len(params) > 2 {
			m = params[2]
		}
	}
	if gw != 0 && gh != 0 && (gridWidth == 0 || gridHeight == 0) {
		gridWidth, gridHeight = gw, gh
	}
	if m != 0 && margin == 0 {
		margin = m
	}
}

// NewDraw は描画オブジェクトを生成する
func NewDraw(w, h int) (d *Data) {
	if gridWidth == 0 || gridHeight == 0 {
		gridWidth, gridHeight = defaultGridWidth, defaultGridHeight
		// [MEMO] 以後、gridWidth と gridHeght はこの値が使われる。変更もできない。
	}
	if margin == 0 {
		margin = defaultMargin
		// [MEMO] 以後、margin はこの値が使われる。変更もできない。
	}
	d = &Data{
		width:   w,
		height:  h,
		bgColor: white,
		fgColor: black,
	}
	d.rgb = image.NewRGBA(image.Rect(0, 0, d.width*gridWidth+margin*2, d.height*gridHeight+margin*2))
	d.fillAll(d.bgColor)
	return
}

func (d *Data) saveFile(outFile string) {
	w, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}
	defer w.Close()
	png.Encode(w, d.rgb)
}

func (d *Data) drawGridSquares() {
	for x := 0; x <= d.width; x++ {
		d.drawGridLine(x, 0, x, d.height)
	}
	for y := 0; y <= d.height; y++ {
		d.drawGridLine(0, y, d.width, y)
	}
}

func (d *Data) setFgColor(c color.Color) {
	d.fgColor = c
}

func getColorOfStr(c string) (r color.Color) {
	switch c {
	case "black":
		r = black
	case "white":
		r = white
	case "yellow":
		r = yellow
	case "skyblue":
		r = skyblue
	case "red":
		r = red
	case "green":
		r = green
	case "blue":
		r = blue
	case "pink":
		r = pink
	case "gray":
		r = gray
	}
	return
}

func (d *Data) setBgColor(c color.Color) {
	d.bgColor = c
	d.fillAll(d.bgColor)
}

func (d *Data) fillAll(c color.Color) {
	// 矩形を取得
	rect := d.rgb.Rect

	// 全部埋める
	for h := rect.Min.Y; h < rect.Max.Y; h++ {
		for v := rect.Min.X; v < rect.Max.X; v++ {
			d.rgb.Set(v, h, c)
		}
	}
}

func (d *Data) fillSquare(x1, y1, x2, y2 int) {
	//fmt.Println("#", x1, y1, x2, y2)
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			//fmt.Println("#fill", x, y)
			d.rgb.Set(x+margin, y+margin, d.fgColor)
		}
	}
}

func (d *Data) fillGridSquare(x1, y1, x2, y2 int) {
	d.fillSquare(x1*gridWidth, y1*gridHeight, x2*gridWidth, y2*gridHeight)
}

func (d *Data) drawLine(x1, y1, x2, y2 int) {
	if x1 == x2 && y1 == y2 {
		//fmt.Println("#1", x1, y1, x2, y2)
		d.rgb.Set(x1+margin, y1+margin, d.fgColor)
	} else if x1 == x2 {
		//fmt.Println("#2:Y", x1, y1, x2, y2)
		for y := y1; y <= y2; y++ {
			d.rgb.Set(x1+margin, y+margin, d.fgColor)
		}
		return
	} else if y1 == y2 {
		//fmt.Println("#3:X", x1, y1, x2, y2)
		for x := x1; x <= x2; x++ {
			d.rgb.Set(x+margin, y1+margin, d.fgColor)
		}
		return
	} else {
		//fmt.Println("#4", x1, y1, x2, y2)
		dx := float32(x2 - x1)
		dy := float32(y2 - y1)
		a := dy / dx
		f := func(x float32) float32 {
			return a*(x-float32(x1)) + float32(y1)
		}
		//l := math.Sqrt(math.Power(dx, 2.0)+math.Power(dy, 2.0))
		ox1, ox2 := x1, x2
		if ox1 > ox2 {
			ox1, ox2 = ox2, ox1
		}

		for x := float32(ox1); x <= float32(ox2); x++ {
			y := f(x)
			d.rgb.Set(int(x)+margin, int(y)+margin, d.fgColor)
		}

		g := func(y float32) float32 {
			return (y-float32(y1))/a + float32(x1)
		}
		//l := math.Sqrt(math.Power(dx, 2.0)+math.Power(dy, 2.0))
		oy1, oy2 := y1, y2
		if oy1 > oy2 {
			oy1, oy2 = oy2, oy1
		}

		for y := float32(oy1); y <= float32(oy2); y++ {
			x := g(y)
			d.rgb.Set(int(x)+margin, int(y)+margin, d.fgColor)
		}

	}
}

func (d *Data) drawGridLine(x1, y1, x2, y2 int) {
	d.drawLine(x1*gridWidth, y1*gridHeight, x2*gridWidth, y2*gridHeight)
}

func (d *Data) drawText(x, y int, text string) {
	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}

	dr := &font.Drawer{
		Dst:  d.rgb,
		Src:  image.NewUniform(d.fgColor),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	dr.DrawString(text)
}

func (d *Data) drawGridText(gx, gy int, text string) {
	d.drawText(gx*gridWidth, gy*gridHeight, text)
}
