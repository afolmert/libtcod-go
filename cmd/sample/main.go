package main

//
import (
	"container/vector"
	"crypto/rand"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	. "tcod"
)

//
// misc functions
//
//

func sin(f float32) float32 {
	return float32(math.Sin(float64(f)))
}

func cos(f float32) float32 {
	return float32(math.Cos(float64(f)))
}

func atoi(s string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		return 0
	} else {
		return result
	}
	return 0
}


func randByte() byte {
	b := make([]byte, 1)
	rand.Read(b)
	return b[0]
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
	return a
}

func minf(a, b float32) float32 {
	if a < b {
		return a
	} else {
		return b
	}
	return a
}

func max(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
	return a
}

func maxf(a, b float32) float32 {
	if a < b {
		return b
	} else {
		return a
	}
	return a
}

func randColor() Color {
	return Color{randByte(), randByte(), randByte()}
}


func sqr(i int) int {
	return i * i
}

func sqrf(i float32) float32 {
	return i * i
}

func sqrt(i float32) float32 {
	return float32(math.Sqrt(float64(i)))
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func absf(i float32) float32 {
	if i < 0 {
		return -i
	}
	return i
}

//
// Demo
//

type Demo interface {
	Render(first bool, key *Key)
}

//
// Global variables
//
var demoConsole *Console
var rootConsole *RootConsole
var guiConsole *Console
var gui *Gui

var random *Random

const SCREEN_WIDTH = 108
const SCREEN_HEIGHT = 66

const DEMO_SCREEN_WIDTH = 75
const DEMO_SCREEN_HEIGHT = 60

const DEMO_SCREEN_X = 30
const DEMO_SCREEN_Y = 3

const GUI_SCREEN_WIDTH = SCREEN_WIDTH
const GUI_SCREEN_HEIGHT = SCREEN_HEIGHT

var curSample = 0
var first = true

type Sample struct {
	name   string
	demo   Demo
	button *RadioButton
}

var samples []Sample


///
/// Ascii Art sample
///
//
const (
	TOPLEFT = iota
	TOPRIGHT
	BOTTOMLEFT
	BOTTOMRIGHT
)


const ASCII_ART_SLIDE_DELAY = 3 // every 3 seconds new slide shows

type AsciiArtDemo struct {
	asciiArt      IAsciiArt
	lastSlideTime float32 // last
	transition    Transition
}


func NewAsciiArtDemo() *AsciiArtDemo {
	return &AsciiArtDemo{}
}



func (self *AsciiArtDemo) Render(first bool, key *Key) {
	if aas == nil {
		aas = NewAsciiArtGallery()
	}

	if self.asciiArt == nil || SysElapsedSeconds()-self.lastSlideTime > ASCII_ART_SLIDE_DELAY {
		self.asciiArt = aas.randomAsciiArt()
		first = true
		if self.asciiArt == nil {
			demoConsole.SetDefaultBackground(COLOR_BLACK)
			demoConsole.SetDefaultForeground(COLOR_GREY)
			demoConsole.Clear()
			demoConsole.PrintRectEx(0, 0, demoConsole.GetWidth(), demoConsole.GetHeight(), BKGND_SET, LEFT, "Loading images, please wait...")
		} else {
			self.transition = Transition(random.GetInt(INSTANT+1, NB_TRANSITIONS-1))
			self.lastSlideTime = SysElapsedSeconds()
		}
	}

	if self.asciiArt != nil {
		self.asciiArt.Render(demoConsole, self.transition, first)
	}
}


// Old Colors demo
type ColorsDemo struct {
	cols             [4]Color // random corner colors
	dirr, dirg, dirb [4]int
}

func NewColorsDemo() *ColorsDemo {
	result := &ColorsDemo{}
	result.cols = [4]Color{
		Color{50, 40, 150},
		Color{240, 85, 5},
		Color{50, 35, 240},
		Color{10, 200, 130}}
	result.dirr = [4]int{1, -1, 1, 1}
	result.dirg = [4]int{1, -1, -1, 1}
	result.dirb = [4]int{1, 1, 1, -1}
	return result
}


func (self *ColorsDemo) Render(first bool, key *Key) {
	var c, x, y int
	var textColor Color
	// slighty modify the corner colors
	if first {
		SysSetFps(0)
		demoConsole.Clear()
	}
	// slighty modify the corner colors
	for c = 0; c < 4; c++ {
		// move each corner color
		component := random.GetInt(0, 2)
		switch component {
		case 0:
			self.cols[c].R = byte(int(self.cols[c].R) + 5*self.dirr[c])
			if self.cols[c].R == 255 {
				self.dirr[c] = -1
			} else if self.cols[c].R == 0 {
				self.dirr[c] = 1
			}
		case 1:
			self.cols[c].G = byte(int(self.cols[c].G) + 5*self.dirg[c])
			if self.cols[c].G == 255 {
				self.dirg[c] = -1
			} else if self.cols[c].G == 0 {
				self.dirg[c] = 1
			}
		case 2:
			self.cols[c].B = byte(int(self.cols[c].B) + 5*self.dirb[c])
			if self.cols[c].B == 255 {
				self.dirb[c] = -1
			} else if self.cols[c].B == 0 {
				self.dirb[c] = 1
			}
		}
	}

	// scan the whole screen, interpolating corner colors
	for x = 0; x < DEMO_SCREEN_WIDTH; x++ {
		xcoef := float32(x) / (DEMO_SCREEN_WIDTH - 1)
		// get the current column top and bottom colors
		top := self.cols[TOPLEFT].Lerp(self.cols[TOPRIGHT], xcoef)
		bottom := self.cols[BOTTOMLEFT].Lerp(self.cols[BOTTOMRIGHT], xcoef)
		for y = 0; y < DEMO_SCREEN_HEIGHT; y++ {
			ycoef := float32(y) / (DEMO_SCREEN_HEIGHT - 1)
			// get the current cell color
			curColor := top.Lerp(bottom, ycoef)
			demoConsole.SetCharBackground(x, y, curColor, BKGND_SET)
		}
	}

	// print the text
	// get the background color at the text position
	textColor = demoConsole.GetCharBackground(DEMO_SCREEN_WIDTH/2, 5)
	// and invert it
	textColor.R = 255 - textColor.R
	textColor.G = 255 - textColor.G
	textColor.B = 255 - textColor.B
	demoConsole.SetDefaultForeground(textColor)
	// put random text (for performance tests)
	for x = 0; x < DEMO_SCREEN_WIDTH; x++ {
		for y = 0; y < DEMO_SCREEN_HEIGHT; y++ {
			var c int
			col := demoConsole.GetCharBackground(x, y)
			col = col.Lerp(COLOR_BLACK, 0.5)
			c = random.GetInt('a', 'z')
			demoConsole.SetDefaultForeground(col)
			demoConsole.PutChar(x, y, c, BKGND_NONE)
		}
	}
	// the background behind the text is slightly darkened using the BKGND_MULTIPLY flag
	demoConsole.SetDefaultBackground(COLOR_GREY)
	demoConsole.PrintRectEx(DEMO_SCREEN_WIDTH/2, 5, DEMO_SCREEN_WIDTH-2, DEMO_SCREEN_HEIGHT-1, BKGND_MULTIPLY, CENTER,
		"The Doryen library uses 24 bits colors, for both background and foreground.")
}
//


// Offscreen demo
///
///
///
type OffscreenDemo struct {
	secondary  *Console // second screen
	screenshot *Console // screenshot screen
	init       bool     // draw the secondary screen only the first time
	counter    int
	x, y       int // secondary screen position
	xdir, ydir int // movement direction
}

func NewOffscreenDemo() *OffscreenDemo {
	return &OffscreenDemo{
		init:    false,
		counter: 0,
		x:       0,
		y:       0,
		xdir:    1,
		ydir:    1}
}


func (self *OffscreenDemo) Render(first bool, key *Key) {
	if !self.init {
		self.init = true
		self.secondary = NewConsole(DEMO_SCREEN_WIDTH/2, DEMO_SCREEN_HEIGHT/2)
		self.screenshot = NewConsole(DEMO_SCREEN_WIDTH, DEMO_SCREEN_HEIGHT)
		self.secondary.PrintFrame(0, 0, DEMO_SCREEN_WIDTH/2, DEMO_SCREEN_HEIGHT/2, false, BKGND_SET, "Offscreen console")
		self.secondary.PrintRectEx(DEMO_SCREEN_WIDTH/4, 2, DEMO_SCREEN_WIDTH/2-2, DEMO_SCREEN_HEIGHT/2,
			BKGND_NONE, CENTER, "You can render to an offscreen console and blit in on another one, simulating alpha transparency.")
	}
	if first {
		SysSetFps(30) // limited to 30 fps
		// get a "screenshot" of the current sample screen
		demoConsole.Blit(0, 0, DEMO_SCREEN_WIDTH, DEMO_SCREEN_HEIGHT,
			self.screenshot, 0, 0, 1.0, 1.0)
	}
	self.counter++
	if self.counter%20 == 0 {
		// move the secondary screen every 2 seconds
		self.x += self.xdir
		self.y += self.ydir
		if self.x == DEMO_SCREEN_WIDTH/2+5 {
			self.xdir = -1
		} else if self.x == -5 {
			self.xdir = 1
		}
		if self.y == DEMO_SCREEN_HEIGHT/2+5 {
			self.ydir = -1
		} else if self.y == -5 {
			self.ydir = 1
		}
	}
	// restore the initial screen
	self.screenshot.Blit(0, 0, DEMO_SCREEN_WIDTH, DEMO_SCREEN_HEIGHT,
		demoConsole, 0, 0, 1.0, 1.0)
	// blit the overlapping screen
	self.secondary.Blit(0, 0, DEMO_SCREEN_WIDTH/2, DEMO_SCREEN_HEIGHT/2,
		demoConsole, self.x, self.y, 1.0, 0.6)

}


///
/// Lines demo
///
type LinesDemo struct {
	bk        *Console  // colored background
	bkFlag    BkgndFlag // current blending mode
	init      bool
	flagNames []string
}


func NewLinesDemo() *LinesDemo {
	result := &LinesDemo{
		init:   false,
		bkFlag: BKGND_SET}
	result.flagNames = []string{
		"BKGND_NONE",
		"BKGND_SET",
		"BKGND_MULTIPLY",
		"BKGND_LIGHTEN",
		"BKGND_DARKEN",
		"BKGND_SCREEN",
		"BKGND_COLOR_DODGE",
		"BKGND_COLOR_BURN",
		"BKGND_ADD",
		"BKGND_ADDALPHA",
		"BKGND_BURN",
		"BKGND_OVERLAY",
		"BKGND_ALPHA"}
	return result
}


func lineListener(x, y int, demo interface{}) bool {
	if x >= 0 && y >= 0 && x < DEMO_SCREEN_WIDTH && y < DEMO_SCREEN_HEIGHT {
		demoConsole.SetCharBackground(x, y, COLOR_LIGHT_BLUE, BkgndFlag(demo.(*LinesDemo).bkFlag))
	}
	return true
}


func (self *LinesDemo) Render(first bool, key *Key) {
	var xo, yo, xd, yd, x, y int            // segment starting, ending, current position
	var alpha float32                       // alpha value when blending mode = BKGND_ALPHA
	var angle, cos_angle, sin_angle float32 // segment angle data
	var recty int                           // gradient vertical position
	if key.Vk == K_ENTER || key.Vk == K_KPENTER {
		// switch to the next blending mode
		self.bkFlag++
		if (self.bkFlag & 0xff) > BKGND_ALPH {
			self.bkFlag = BKGND_NONE
		}
	}
	if (self.bkFlag & 0xff) == BKGND_ALPH {
		// for the alpha mode, update alpha every frame
		alpha = (1.0 + cos(float32(SysElapsedSeconds()))*2) / 2.0
		self.bkFlag = BkgndAlpha(alpha)
	} else if (self.bkFlag & 0xff) == BKGND_ADDA {
		// for the add alpha mode, update alpha every frame
		alpha = (1.0 + cos(float32(SysElapsedSeconds()))*2) / 2.0
		self.bkFlag = BkgndAddAlpha(alpha)
	}
	if !self.init {
		self.bk = NewConsole(DEMO_SCREEN_WIDTH, DEMO_SCREEN_HEIGHT)
		// initialize the colored background
		for x = 0; x < DEMO_SCREEN_WIDTH; x++ {
			for y = 0; y < DEMO_SCREEN_HEIGHT; y++ {
				var col Color
				col.R = (uint8)(x * 255 / (DEMO_SCREEN_WIDTH - 1))
				col.G = (uint8)((x + y) * 255 / (DEMO_SCREEN_WIDTH - 1 + DEMO_SCREEN_HEIGHT - 1))
				col.B = (uint8)(y * 255 / (DEMO_SCREEN_HEIGHT - 1))
				self.bk.SetCharBackground(x, y, col, BKGND_SET)
			}
		}
		self.init = true
	}
	if first {
		SysSetFps(30) // limited to 30 fps
		demoConsole.SetDefaultForeground(COLOR_WHITE)
	}
	// blit the background
	self.bk.Blit(0, 0, DEMO_SCREEN_WIDTH, DEMO_SCREEN_HEIGHT, demoConsole, 0, 0, 1.0, 1.0)
	// render the gradient
	recty = (int)((DEMO_SCREEN_HEIGHT - 2) * ((1.0 + cos(float32(SysElapsedSeconds()))) / 2.0))
	for x = 0; x < DEMO_SCREEN_WIDTH; x++ {
		var col Color
		col.R = uint8(x * 255 / DEMO_SCREEN_WIDTH)
		col.G = uint8(x * 255 / DEMO_SCREEN_WIDTH)
		col.B = uint8(x * 255 / DEMO_SCREEN_WIDTH)
		demoConsole.SetCharBackground(x, recty, col, BkgndFlag(self.bkFlag))
		demoConsole.SetCharBackground(x, recty+1, col, BkgndFlag(self.bkFlag))
		demoConsole.SetCharBackground(x, recty+2, col, BkgndFlag(self.bkFlag))
	}
	// calculate the segment ends
	angle = float32(SysElapsedSeconds()) * 2.0
	cos_angle = cos(angle)
	sin_angle = sin(angle)
	xo = (int)(DEMO_SCREEN_WIDTH / 2 * (1 + cos_angle))
	yo = (int)(DEMO_SCREEN_HEIGHT/2 + sin_angle*DEMO_SCREEN_WIDTH/2)
	xd = (int)(DEMO_SCREEN_WIDTH / 2 * (1 - cos_angle))
	yd = (int)(DEMO_SCREEN_HEIGHT/2 - sin_angle*DEMO_SCREEN_WIDTH/2)
	// render the line
	Line(xo, yo, xd, yd, self, lineListener)
	// print the current flag
	demoConsole.PrintEx(2, 2, BKGND_NONE, LEFT, "%s (ENTER to change)", self.flagNames[self.bkFlag&0xff])
}


///
/// NoiseDemo
///

const (
	PERLIN = iota
	SIMPLEX
	WAVELET
	FBM_PERLIN
	TURBULENCE_PERLIN
	FBM_SIMPLEX
	TURBULENCE_SIMPLEX
	FBM_WAVELET
	TURBULENCE_WAVELET
)


type NoiseDemo struct {
	funName    []string
	fun        int
	noise      *Noise
	dx, dy     float32
	octaves    float32
	hurst      float32
	lacunarity float32
	img        *Image
	zoom       float32
}

func NewNoiseDemo() *NoiseDemo {
	result := &NoiseDemo{
		fun:        PERLIN,
		noise:      nil,
		dx:         0.0,
		dy:         0.0,
		octaves:    4.0,
		hurst:      NOISE_DEFAULT_HURST,
		lacunarity: NOISE_DEFAULT_LACUNARITY,
		img:        nil,
		zoom:       3.0}

	result.funName = []string{
		"1 : perlin noise       ",
		"2 : simplex noise      ",
		"3 : wavelet noise      ",
		"4 : perlin fbm         ",
		"5 : perlin turbulence  ",
		"6 : simplex fbm        ",
		"7 : simplex turbulence ",
		"8 : wavelet fbm        ",
		"9 : wavelet turbulence ",
	}

	return result
}


func (self *NoiseDemo) Render(first bool, key *Key) {
	var x, y, curfun int
	if self.noise == nil {
		self.noise = NewNoiseWithOptions(2, self.hurst, self.lacunarity, random)
		self.img = NewImage(DEMO_SCREEN_WIDTH*2, DEMO_SCREEN_HEIGHT*2)
	}
	if first {
		SysSetFps(30) // limited to 30 fps
	}
	demoConsole.Clear()
	// texture animation
	self.dx += 0.01
	self.dy += 0.01
	// render the 2d noise fun
	for y = 0; y < 2*DEMO_SCREEN_HEIGHT; y++ {
		for x = 0; x < 2*DEMO_SCREEN_WIDTH; x++ {
			var value float32
			var c uint8
			var col Color
			f := make([]float32, 2)
			f[0] = self.zoom*float32(x)/(2*DEMO_SCREEN_WIDTH) + self.dx
			f[1] = self.zoom*float32(y)/(2*DEMO_SCREEN_HEIGHT) + self.dy
			value = 0.0
			switch self.fun {
			case PERLIN:
				value = self.noise.GetEx(f, NOISE_PERLIN)
			case SIMPLEX:
				value = self.noise.GetEx(f, NOISE_SIMPLEX)
			case WAVELET:
				value = self.noise.GetEx(f, NOISE_WAVELET)
			case FBM_PERLIN:
				value = self.noise.GetFbmEx(f, self.octaves, NOISE_PERLIN)
			case TURBULENCE_PERLIN:
				value = self.noise.GetTurbulenceEx(f, self.octaves, NOISE_PERLIN)
			case FBM_SIMPLEX:
				value = self.noise.GetFbmEx(f, self.octaves, NOISE_SIMPLEX)
			case TURBULENCE_SIMPLEX:
				value = self.noise.GetTurbulenceEx(f, self.octaves, NOISE_SIMPLEX)
			case FBM_WAVELET:
				value = self.noise.GetFbmEx(f, self.octaves, NOISE_WAVELET)
			case TURBULENCE_WAVELET:
				value = self.noise.GetTurbulenceEx(f, self.octaves, NOISE_WAVELET)
			}
			c = (uint8)((value + 1.0) / 2.0 * 255)
			// use a bluish color
			col.R = c / 2
			col.R = col.G
			col.B = c
			self.img.PutPixel(x, y, col)
		}
	}
	// blit the noise image with subcell resolution
	self.img.Blit2x(demoConsole, 0, 0, 0, 0, -1, -1)

	// draw a transparent rectangle
	demoConsole.SetDefaultBackground(COLOR_GREY)
	demoConsole.Rect(2, 2, 23, If(self.fun <= WAVELET, 10, 13).(int), false, BKGND_MULTIPLY)
	for y = 2; y < 2+(If(self.fun <= WAVELET, 10, 13).(int)); y++ {
		for x = 2; x < 2+23; x++ {
			col := demoConsole.GetCharForeground(x, y)
			col = col.Multiply(COLOR_GREY)
			demoConsole.SetCharForeground(x, y, col)
		}
	}

	// draw the text
	for curfun = PERLIN; curfun <= TURBULENCE_WAVELET; curfun++ {
		if curfun == self.fun {
			demoConsole.SetDefaultForeground(COLOR_WHITE)
			demoConsole.SetDefaultBackground(COLOR_LIGHT_BLUE)
			demoConsole.PrintEx(2, 2+curfun, BKGND_SET, LEFT, self.funName[curfun])
		} else {
			demoConsole.SetDefaultForeground(COLOR_GREY)
			demoConsole.PrintEx(2, 2+curfun, BKGND_NONE, LEFT, self.funName[curfun])
		}
	}
	// draw parameters
	demoConsole.SetDefaultForeground(COLOR_WHITE)
	demoConsole.PrintEx(2, 11, BKGND_NONE, LEFT, "Y/H : zoom (%2.1f)", self.zoom)
	if self.fun > WAVELET {
		demoConsole.PrintEx(2, 12, BKGND_NONE, LEFT, "E/D : hurst (%2.1f)", self.hurst)
		demoConsole.PrintEx(2, 13, BKGND_NONE, LEFT, "R/F : lacunarity (%2.1f)", self.lacunarity)
		demoConsole.PrintEx(2, 14, BKGND_NONE, LEFT, "T/G : octaves (%2.1f)", self.octaves)
	}
	// handle keypress
	if key.Vk == K_NONE {
		return
	}
	if key.C >= '1' && key.C <= '9' {
		// change the noise fun
		self.fun = int(key.C - '1')
	} else if key.C == 'E' || key.C == 'e' {
		// increase hurst
		self.hurst += 0.1
		self.noise = NewNoiseWithOptions(2, self.hurst, self.lacunarity, random)
	} else if key.C == 'D' || key.C == 'd' {
		// decrease hurst
		self.hurst -= 0.1
		self.noise = NewNoiseWithOptions(2, self.hurst, self.lacunarity, random)
	} else if key.C == 'R' || key.C == 'r' {
		// increase lacunarity
		self.lacunarity += 0.5
		self.noise = NewNoiseWithOptions(2, self.hurst, self.lacunarity, random)
	} else if key.C == 'F' || key.C == 'f' {
		// decrease lacunarity
		self.lacunarity -= 0.5
		self.noise = NewNoiseWithOptions(2, self.hurst, self.lacunarity, random)
	} else if key.C == 'T' || key.C == 't' {
		// increase octaves
		self.octaves += 0.5
	} else if key.C == 'G' || key.C == 'g' {
		// decrease octaves
		self.octaves -= 0.5
	} else if key.C == 'Y' || key.C == 'y' {
		// increase zoom
		self.zoom += 0.2
	} else if key.C == 'H' || key.C == 'h' {
		// decrease zoom
		self.zoom -= 0.2
	}
}


///
/// Fov demo
///

const TORCH_RADIUS = 10.0
const SQUARED_TORCH_RADIUS = TORCH_RADIUS * TORCH_RADIUS

type FovDemo struct {
	smap         []string
	px, py       int // player position
	recomputeFov bool
	torch        bool
	lightWalls   bool
	map_         *Map
	darkWall     Color
	lightWall    Color
	darkGround   Color
	lightGround  Color
	noise        *Noise
	algoNum      int
	algoNames    []string
	torchx       []float32 // torch light position in the perlin noise
}


func NewFovDemo() *FovDemo {
	result := &FovDemo{
		px:           20,
		py:           10,
		recomputeFov: true,
		torch:        true,
		lightWalls:   true,
		map_:         nil,
		darkWall:     Color{0, 0, 100},
		lightWall:    Color{130, 110, 50},
		darkGround:   Color{50, 50, 150},
		lightGround:  Color{200, 180, 50},
		noise:        nil,
		algoNum:      0,
		torchx:       make([]float32, 1)}

	result.torchx[0] = 0.0
	result.algoNames = []string{
		"BASIC      ", "DIAMOND    ", "SHADOW     ",
		"PERMISSIVE0", "PERMISSIVE1", "PERMISSIVE2", "PERMISSIVE3", "PERMISSIVE4",
		"PERMISSIVE5", "PERMISSIVE6", "PERMISSIVE7", "PERMISSIVE8", "RESTRICTIVE"}

	result.smap = []string{
		"#############################################################################",
		"##           ##########      ##########################                  ####",
		"##           ########    #     ######################## ################ ####",
		"##           #########  ###        #################### ################ ####",
		"##           #####      #####             ############# ################ ####",
		"##           ###       ########    ###### ############# ################ ####",
		"##  ###########      #################### #############               ## ####",
		"##  ############    ######                  ######################## ### ####",
		"##  ####   #######  ######   #     #     #  ######################## ### ####",
		"## #####   ######      ###                  ######################## ### ####",
		"## #####                                                         ### ### ####",
		"## #####                                    #################### ### ### ####",
		"## ##############################           #################### ### ### ####",
		"## ##############################           #################### ### ### ####",
		"## ##############################           #################### ### ### ####",
		"## ##############################                                ### ### ####",
		"## #####         #       #                  ######################## ### ####",
		"## #####     #   #       #                  ##     #        ##       ### ####",
		"## #####     #   #       #                  ##     #        ##       ### ####",
		"## ##### #####   #       #                  ##                       ### ####",
		"## ##### ######  #       #                  ##     #        ##       ### ####",
		"##            #  #       ################   ##     #        ##       ### ####",
		"##            #  #       ##                 ##     ##################### ####",
		"##            #          # #                ##     #        ##       ### ####",
		"##            ########## #  #               ##     #        ##       ### ####",
		"##            #          #   #              ##                       ### ####",
		"##            #          #    #             ##     #         #       ### ####",
		"###############          #     #            ##     #         #       ### ####",
		"# ######                 #      #           ##     #         #       ### ####",
		"# #####################  #       #          ##     #         #       ### ####",
		"# #                      #        #         ############################ ####",
		"# #                      #         ####                              ### ####",
		"# #                                                                  ### ####",
		"# #                                         ###################  ####### ####",
		"# # ####################################    ##              ###  ####### ####",
		"# #      ###############################    ##  ###############  ####### ####",
		"# ######                              ##    ##  ###############  ####### ####",
		"#        ############################ ##    ##                   ####### ####",
		"# ################################### ##    ##                   ####### ####",
		"# #                                   ##    ##  ###############  ####### ####",
		"# # ####################################    ##  ###############  ####### ####",
		"# # ####################################    ##              ###  ####### ####",
		"# #                                   ##    ###################  ####### ####",
		"# ################################### ##    ############################ ####",
		"# ################################### ##                              ## ####",
		"# #                                   ##     #                     #  ## ####",
		"# # ####################################                              ## ####",
		"# # ####################################                              ## ####",
		"# #                                   ##                 #            ## ####",
		"# ################################### ##                              ## ####",
		"#                                     ##     #                     #  ## ####",
		"########################################                              ## ####",
		"########                                    ########          ########## ####",
		"########                                    ########          ########## ####",
		"####       ######      ###   #     #     #  ########          ########## ####",
		"#### ###   ########## ####                  ########          ########## ####",
		"#### ###   ##########   ###########=##################     ############# ####",
		"#### ##################   #####          #############     ############# ####",
		"#### ###             #### #####          #############                   ####",
		"####           #     ####                ####################################",
		"########       #     #### #####          ####################################"}
	return result
}


func (self *FovDemo) Render(first bool, key *Key) {
	var x, y int
	// torch position & intensity variation
	var dx, dy, di float32 = 0, 0, 0

	if self.map_ == nil {
		self.map_ = NewMap(DEMO_SCREEN_WIDTH, DEMO_SCREEN_HEIGHT)
		for y = 0; y < DEMO_SCREEN_HEIGHT; y++ {
			for x = 0; x < DEMO_SCREEN_WIDTH; x++ {
				if self.smap[y][x] == ' ' {
					self.map_.SetProperties(x, y, true, true) // ground
				} else if self.smap[y][x] == '=' {
					self.map_.SetProperties(x, y, true, false) // window
				}
			}
		}
		self.noise = NewNoiseWithOptions(1, 1.0, 1.0, random) // 1d noise for the torch flickering
	}
	if first {
		SysSetFps(30) // limited to 30 fps
		// we draw the foreground only the first time.
		//   during the player movement, only the @ is redrawn.
		//   the rest impacts only the background color
		// draw the help text & player @
		demoConsole.Clear()
		demoConsole.SetDefaultForeground(COLOR_WHITE)
		demoConsole.PrintEx(1, 0, BKGND_NONE, LEFT, "IJKL : move around\nT : torch fx %s\nW : light walls %s\n+-: algo %s",
			If(self.torch, "on ", "off").(string), If(self.lightWalls, "on ", "off").(string), self.algoNames[self.algoNum])
		demoConsole.SetDefaultForeground(COLOR_BLACK)
		demoConsole.PutChar(self.px, self.py, '@', BKGND_NONE)
		// draw windows
		for y = 0; y < DEMO_SCREEN_HEIGHT; y++ {
			for x = 0; x < DEMO_SCREEN_WIDTH; x++ {
				if self.smap[y][x] == '=' {
					demoConsole.PutChar(x, y, CHAR_DHLINE, BKGND_NONE)
				}
			}
		}
	}
	if self.recomputeFov {
		self.recomputeFov = false
		self.map_.ComputeFov(self.px, self.py, If(self.torch, int(TORCH_RADIUS), 0).(int), self.lightWalls, FovAlgorithm(self.algoNum))
	}
	if self.torch {
		tdx := make([]float32, 1)
		// slightly change the perlin noise parameter
		self.torchx[0] += 0.2
		// randomize the light position between -1.5 and 1.5
		tdx[0] = self.torchx[0] + 20.0
		dx = self.noise.GetEx(tdx, NOISE_SIMPLEX) * 1.5
		tdx[0] += 30.0
		dy = self.noise.GetEx(tdx, NOISE_SIMPLEX) * 1.5
		di = 0.2 * self.noise.GetEx(self.torchx, NOISE_SIMPLEX)
	}
	for y = 0; y < DEMO_SCREEN_HEIGHT; y++ {
		for x = 0; x < DEMO_SCREEN_WIDTH; x++ {
			visible := self.map_.IsInFov(x, y)
			wall := self.smap[y][x] == '#'
			if !visible {
				demoConsole.SetCharBackground(x, y,
					If(wall, self.darkWall, self.darkGround).(Color), BKGND_SET)

			} else {
				if !self.torch {
					demoConsole.SetCharBackground(x, y,
						If(wall, self.lightWall, self.lightGround).(Color), BKGND_SET)
				} else {
					base := If(wall, self.darkWall, self.darkGround).(Color)
					light := If(wall, self.lightWall, self.lightGround).(Color)
					r := (x-self.px+int(dx))*(x-self.px+int(dx)) + (y-self.py+int(dy))*(y-self.py+int(dy)) // cell distance to torch (squared)
					if r < SQUARED_TORCH_RADIUS {
						l := float32(SQUARED_TORCH_RADIUS-r)/float32(SQUARED_TORCH_RADIUS) + di
						l = ClampF(0.0, 1.0, l)
						base = base.Lerp(light, l)
					}
					if x >= 0 && x < DEMO_SCREEN_WIDTH && y >= 0 && y < DEMO_SCREEN_HEIGHT {
						demoConsole.SetCharBackground(x, y, base, BKGND_SET)
					}
				}
			}
		}
	}
	if key.C == 'I' || key.C == 'i' {
		if self.smap[self.py-1][self.px] == ' ' {
			demoConsole.PutChar(self.px, self.py, ' ', BKGND_NONE)
			self.py--
			demoConsole.PutChar(self.px, self.py, '@', BKGND_NONE)
			self.recomputeFov = true
		}
	} else if key.C == 'K' || key.C == 'k' {
		if self.smap[self.py+1][self.px] == ' ' {
			demoConsole.PutChar(self.px, self.py, ' ', BKGND_NONE)
			self.py++
			demoConsole.PutChar(self.px, self.py, '@', BKGND_NONE)
			self.recomputeFov = true
		}
	} else if key.C == 'J' || key.C == 'j' {
		if self.smap[self.py][self.px-1] == ' ' {
			demoConsole.PutChar(self.px, self.py, ' ', BKGND_NONE)
			self.px--
			demoConsole.PutChar(self.px, self.py, '@', BKGND_NONE)
			self.recomputeFov = true
		}
	} else if key.C == 'L' || key.C == 'l' {
		if self.smap[self.py][self.px+1] == ' ' {
			demoConsole.PutChar(self.px, self.py, ' ', BKGND_NONE)
			self.px++
			demoConsole.PutChar(self.px, self.py, '@', BKGND_NONE)
			self.recomputeFov = true
		}
	} else if key.C == 'T' || key.C == 't' {
		self.torch = !self.torch
		demoConsole.SetDefaultForeground(COLOR_WHITE)
		demoConsole.PrintEx(1, 0, BKGND_NONE, LEFT, "IJKL : move around\nT : torch fx %s\nW : light walls %s\n+-: algo %s",
			If(self.torch, "on ", "off").(string), If(self.lightWalls, "on ", "off").(string), self.algoNames[self.algoNum])
		demoConsole.SetDefaultForeground(COLOR_BLACK)
	} else if key.C == 'W' || key.C == 'w' {
		self.lightWalls = !self.lightWalls
		demoConsole.SetDefaultForeground(COLOR_WHITE)
		demoConsole.PrintEx(1, 0, BKGND_NONE, LEFT, "IJKL : move around\nT : torch fx %s\nW : light walls %s\n+-: algo %s",
			If(self.torch, "on ", "off").(string), If(self.lightWalls, "on ", "off").(string), self.algoNames[self.algoNum])
		demoConsole.SetDefaultForeground(COLOR_BLACK)
		self.recomputeFov = true
	} else if key.C == '+' || key.C == '-' {
		self.algoNum += If(key.C == '+', 1, -1).(int)
		self.algoNum = Clamp(0, NB_FOV_ALGORITHMS-1, self.algoNum)
		demoConsole.SetDefaultForeground(COLOR_WHITE)
		demoConsole.PrintEx(1, 0, BKGND_NONE, LEFT, "IJKL : move around\nT : torch fx %s\nW : light walls %s\n+-: algo %s",
			If(self.torch, "on ", "off").(string), If(self.lightWalls, "on ", "off").(string), self.algoNames[self.algoNum])
		demoConsole.SetDefaultForeground(COLOR_BLACK)
		self.recomputeFov = true
	}
}


///
///
///
//
type PathDemo struct {
	smap            []string
	px, py          int // player position
	dx, dy          int // destination
	map_            *Map
	darkWall        Color
	darkGround      Color
	lightGround     Color
	path            *Path
	usingAstar      bool
	dijkstraDist    float32
	dijkstra        *Dijkstra
	recalculatePath bool
	busy            float32
	oldChar         int
}

func NewPathDemo() *PathDemo {
	result := &PathDemo{

		px:              20,
		py:              10,
		dx:              24,
		dy:              1,
		map_:            nil,
		darkWall:        Color{0, 0, 100},
		darkGround:      Color{50, 50, 150},
		lightGround:     Color{200, 180, 50},
		path:            nil,
		usingAstar:      true,
		dijkstraDist:    0,
		dijkstra:        nil,
		recalculatePath: false,
		busy:            0,
		oldChar:         ' '}

	result.smap = []string{
		"###########################################################################",
		"##           ##########      ##########################                  ##",
		"##           ########    #     ######################## ################ ##",
		"##           #########  ###        #################### ################ ##",
		"##           #####      #####             ############# ################ ##",
		"##           ###       ########    ###### ############# ################ ##",
		"##  ###########      #################### #############               ## ##",
		"##  ############    ######                  ######################## ### ##",
		"##  ####   #######  ######   #     #     #  ######################## ### ##",
		"## #####   ######      ###                  ######################## ### ##",
		"## #####                                                         ### ### ##",
		"## #####                                    #################### ### ### ##",
		"## ##############################           #################### ### ### ##",
		"## ##############################           #################### ### ### ##",
		"## ##############################           #################### ### ### ##",
		"## ##############################                                ### ### ##",
		"## #####         #       #                  ######################## ### ##",
		"## #####     #   #       #                  ##     #        ##       ### ##",
		"## #####     #   #       #                  ##     #        ##       ### ##",
		"## ##### #####   #       #                  ##                       ### ##",
		"## ##### ######  #       #                  ##     #        ##       ### ##",
		"##            #  #       ################   ##     #        ##       ### ##",
		"##            #  #       ##                 ##     ##################### ##",
		"##            #          # #                ##     #        ##       ### ##",
		"##            ########## #  #               ##     #        ##       ### ##",
		"##            #          #   #              ##                       ### ##",
		"##            #          #    #             ##     #         #       ### ##",
		"###############          #     #            ##     #         #       ### ##",
		"# ######                 #      #           ##     #         #       ### ##",
		"# #####################  #       #          ##     #         #       ### ##",
		"# #                      #        #         ############################ ##",
		"# #                      #         ####                              ### ##",
		"# #                                                                  ### ##",
		"# #                                         ###################  ####### ##",
		"# # ####################################    ##              ###  ####### ##",
		"# #      ###############################    ##  ###############  ####### ##",
		"# ######                              ##    ##  ###############  ####### ##",
		"#        ############################ ##    ##                   ####### ##",
		"# ################################### ##    ##                   ####### ##",
		"# #                                   ##    ##  ###############  ####### ##",
		"# # ####################################    ##  ###############  ####### ##",
		"# # ####################################    ##              ###  ####### ##",
		"# #                                   ##    ###################  ####### ##",
		"# ################################### ##    ############################ ##",
		"# ################################### ##                              ## ##",
		"# #                                   ##     #                     #  ## ##",
		"# # ####################################                              ## ##",
		"# # ####################################                              ## ##",
		"# #                                   ##                 #            ## ##",
		"# ################################### ##                              ## ##",
		"#                                     ##     #                     #  ## ##",
		"########################################                              ## ##",
		"########                                    ########          ########## ##",
		"########                                    ########          ########## ##",
		"####       ######      ###   #     #     #  ########          ########## ##",
		"#### ###   ########## ####                  ########          ########## ##",
		"#### ###   ##########   ###########=##################     ############# ##",
		"#### ##################   #####          #############     ############# ##",
		"#### ###             #### #####          #############                   ##",
		"####           #     ####                ##################################",
		"########       #     #### #####          ##################################"}
	return result
}


func (self *PathDemo) Render(first bool, key *Key) {
	var mouse Mouse
	var mx, my, x, y, i int
	if self.map_ == nil {
		// initialize the map
		self.map_ = NewMap(DEMO_SCREEN_WIDTH, DEMO_SCREEN_HEIGHT)
		for y = 0; y < DEMO_SCREEN_HEIGHT; y++ {
			for x = 0; x < DEMO_SCREEN_WIDTH; x++ {
				if self.smap[y][x] == ' ' { // ground
					self.map_.SetProperties(x, y, true, true)
				} else if self.smap[y][x] == '=' {
					self.map_.SetProperties(x, y, true, false)
				} // window
			}
		}
		self.path = NewPathUsingMap(self.map_, 1.41)
		self.dijkstra = NewDijkstraUsingMap(self.map_, 1.41)
	}
	if first {
		SysSetFps(30) // limited to 30 fps
		// we draw the foreground only the first time.
		// during the player movement, only the @ is redrawn.
		// the rest impacts only the background color
		// draw the help text & player @
		demoConsole.Clear()
		demoConsole.SetDefaultForeground(COLOR_WHITE)
		demoConsole.PutChar(self.dx, self.dy, '+', BKGND_NONE)
		demoConsole.PutChar(self.px, self.py, '@', BKGND_NONE)
		demoConsole.PrintEx(1, 1, BKGND_NONE, LEFT, "IJKL / mouse :\nmove destination\nTAB : Aself.dijkstra")
		demoConsole.PrintEx(1, 4, BKGND_NONE, LEFT, "Using : A*")
		// draw windows
		for y = 0; y < DEMO_SCREEN_HEIGHT; y++ {
			for x = 0; x < DEMO_SCREEN_WIDTH; x++ {
				if self.smap[y][x] == '=' {
					demoConsole.PutChar(x, y, CHAR_DHLINE, BKGND_NONE)
				}
			}
		}

		self.recalculatePath = true
	}
	if self.recalculatePath {
		if self.usingAstar {
			self.path.Compute(self.px, self.py, self.dx, self.dy)
		} else {
			var x, y int
			self.dijkstraDist = 0.0
			// compute the distance grid
			self.dijkstra.Compute(self.px, self.py)
			// get the maximum distance (needed for ground shading only)
			for y = 0; y < DEMO_SCREEN_HEIGHT; y++ {
				for x = 0; x < DEMO_SCREEN_WIDTH; x++ {
					d := self.dijkstra.GetDistance(x, y)
					if d > self.dijkstraDist {
						self.dijkstraDist = d
					}
				}
			}
			// compute the self.path
			self.dijkstra.PathSet(self.dx, self.dy)
		}
		self.recalculatePath = false
		self.busy = 0.2
	}
	// draw the dungeon
	for y = 0; y < DEMO_SCREEN_HEIGHT; y++ {
		for x = 0; x < DEMO_SCREEN_WIDTH; x++ {
			wall := self.smap[y][x] == '#'
			demoConsole.SetCharBackground(x, y, If(wall, self.darkWall, self.darkGround).(Color), BKGND_SET)
		}
	}
	// draw the self.path
	if self.usingAstar {
		for i = 0; i < self.path.Size(); i++ {
			x, y := self.path.Get(i)
			demoConsole.SetCharBackground(x, y, self.lightGround, BKGND_SET)
		}
	} else {
		var x, y, i int
		for y = 0; y < DEMO_SCREEN_HEIGHT; y++ {
			for x = 0; x < DEMO_SCREEN_WIDTH; x++ {
				wall := self.smap[y][x] == '#'
				if !wall {
					d := self.dijkstra.GetDistance(x, y)
					demoConsole.SetCharBackground(x, y, self.lightGround.Lerp(self.darkGround, 0.9*d/self.dijkstraDist), BKGND_SET)
				}
			}
		}
		for i = 0; i < self.dijkstra.Size(); i++ {
			x, y := self.dijkstra.Get(i)
			demoConsole.SetCharBackground(x, y, self.lightGround, BKGND_SET)
		}
	}
	// move the creature
	self.busy -= SysGetLastFrameLength()
	if self.busy <= 0.0 {
		self.busy = 0.2
		if self.usingAstar {
			if !self.path.IsEmpty() {
				demoConsole.PutChar(self.px, self.py, ' ', BKGND_NONE)
				self.px, self.py = self.path.Walk(true)
				demoConsole.PutChar(self.px, self.py, '@', BKGND_NONE)
			}
		} else {
			if !self.dijkstra.IsEmpty() {
				demoConsole.PutChar(self.px, self.py, ' ', BKGND_NONE)
				self.px, self.py = self.dijkstra.PathWalk()
				demoConsole.PutChar(self.px, self.py, '@', BKGND_NONE)
				self.recalculatePath = true
			}
		}
	}
	if (key.C == 'I' || key.C == 'i') && self.dy > 0 {
		// destination move north
		demoConsole.PutChar(self.dx, self.dy, self.oldChar, BKGND_NONE)
		self.dy--
		self.oldChar = demoConsole.GetChar(self.dx, self.dy)
		demoConsole.PutChar(self.dx, self.dy, '+', BKGND_NONE)
		if self.smap[self.dy][self.dx] == ' ' {
			self.recalculatePath = true
		}
	} else if (key.C == 'K' || key.C == 'k') && self.dy < DEMO_SCREEN_HEIGHT-1 {
		// destination move south
		demoConsole.PutChar(self.dx, self.dy, self.oldChar, BKGND_NONE)
		self.dy++
		self.oldChar = demoConsole.GetChar(self.dx, self.dy)
		demoConsole.PutChar(self.dx, self.dy, '+', BKGND_NONE)
		if self.smap[self.dy][self.dx] == ' ' {
			self.recalculatePath = true
		}
	} else if (key.C == 'J' || key.C == 'j') && self.dx > 0 {
		// destination move west
		demoConsole.PutChar(self.dx, self.dy, self.oldChar, BKGND_NONE)
		self.dx--
		self.oldChar = demoConsole.GetChar(self.dx, self.dy)
		demoConsole.PutChar(self.dx, self.dy, '+', BKGND_NONE)
		if self.smap[self.dy][self.dx] == ' ' {
			self.recalculatePath = true
		}
	} else if (key.C == 'L' || key.C == 'l') && self.dx < DEMO_SCREEN_WIDTH-1 {
		// destination move east
		demoConsole.PutChar(self.dx, self.dy, self.oldChar, BKGND_NONE)
		self.dx++
		self.oldChar = demoConsole.GetChar(self.dx, self.dy)
		demoConsole.PutChar(self.dx, self.dy, '+', BKGND_NONE)
		if self.smap[self.dy][self.dx] == ' ' {
			self.recalculatePath = true
		}
	} else if key.Vk == K_TAB {
		self.usingAstar = !self.usingAstar
		if self.usingAstar {
			demoConsole.PrintEx(1, 4, BKGND_NONE, LEFT, "Using : A*      ")
		} else {
			demoConsole.PrintEx(1, 4, BKGND_NONE, LEFT, "Using : self.dijkstra")
		}
		self.recalculatePath = true
	}
	mouse = MouseGetStatus()
	mx = mouse.Cx - DEMO_SCREEN_X
	my = mouse.Cy - DEMO_SCREEN_Y
	if mx >= 0 && mx < DEMO_SCREEN_WIDTH && my >= 0 && my < DEMO_SCREEN_HEIGHT && (self.dx != mx || self.dy != my) {
		demoConsole.PutChar(self.dx, self.dy, self.oldChar, BKGND_NONE)
		self.dx = mx
		self.dy = my
		self.oldChar = demoConsole.GetChar(self.dx, self.dy)
		demoConsole.PutChar(self.dx, self.dy, '+', BKGND_NONE)
		if self.smap[self.dy][self.dx] == ' ' {
			self.recalculatePath = true
		}
	}
}


///
///
///
type SampleMap [DEMO_SCREEN_WIDTH][DEMO_SCREEN_HEIGHT]byte

type BspDemo struct {
	randomRoom  bool // a room fills a random part of the node or the maximum available space ?
	roomWalls   bool // if true, there is always a wall on north & west side of a room
	bspDepth    int
	minRoomSize int
	bsp         *Bsp
	generate    bool
	refresh     bool
	sampleMap   SampleMap
	darkWall    Color
	darkGround  Color
}

//
//typedef	char map_t[DEMO_SCREEN_WIDTH][DEMO_SCREEN_HEIGHT];
//
// TODO check if memset is available ?
//
// TODO !! memset(map,'#',sizeof(char)*DEMO_SCREEN_WIDTH*DEMO_SCREEN_HEIGHT) // TODO
func (self *SampleMap) Fill(c byte) {
	for x := 0; x < len(self); x++ {
		for y := 0; y < len(self[x]); y++ {
			self[x][y] = c
		}
	}
}

//

// draw a vertical line
func (self *SampleMap) VLine(x, y1, y2 int) {
	y := y1
	var dy int
	if y1 > y2 {
		dy = -1
	} else {
		dy = 1
	}
	(*self)[x][y] = ' '
	if y1 == y2 {
		return
	}
	for {
		y += dy
		(*self)[x][y] = ' '
		if y == y2 {
			break
		}
	}
}
//
// draw a vertical line up until we reach an empty space
func (self *SampleMap) VLineUp(x, y int) {
	for y >= 0 && (*self)[x][y] != ' ' {
		(*self)[x][y] = ' '
		y--
	}
}

// draw a vertical line down until we reach an empty space
func (self *SampleMap) VLineDown(x, y int) {
	for y < DEMO_SCREEN_HEIGHT && (*self)[x][y] != ' ' {
		(*self)[x][y] = ' '
		y++
	}
}

// draw a horizontal line
func (self *SampleMap) HLine(x1, y, x2 int) {
	var x int = x1
	var dx int
	if x1 > 2 {
		dx = -1
	} else {
		dx = 1
	}
	(*self)[x][y] = ' '
	if x1 == x2 {
		return
	}
	for {
		x += dx
		(*self)[x][y] = ' '
		if x == x2 {
			break
		}
	}
}

// draw a horizontal line left until we reach an empty space
func (self *SampleMap) HLineLeft(x, y int) {
	for x >= 0 && (*self)[x][y] != ' ' {
		(*self)[x][y] = ' '
		x--
	}
}

// draw a horizontal line right until we reach an empty space
func (self *SampleMap) HLineRight(x, y int) {
	for x < DEMO_SCREEN_WIDTH && (*self)[x][y] != ' ' {
		(*self)[x][y] = ' '
		x++
	}
}

// the class building the dungeon from the bsp nodes
func TraverseNode(node *Bsp, userData interface{}) bool {
	var d *BspDemo = userData.(*BspDemo)
	if node.IsLeaf() {
		// calculate the room size
		minx := node.X + 1
		maxx := node.X + node.W - 1
		miny := node.Y + 1
		maxy := node.Y + node.H - 1
		var x, y int
		if !d.roomWalls {
			if minx > 1 {
				minx--
			}
			if miny > 1 {
				miny--
			}
		}
		if maxx == DEMO_SCREEN_WIDTH-1 {
			maxx--
		}
		if maxy == DEMO_SCREEN_HEIGHT-1 {
			maxy--
		}
		if d.randomRoom {
			minx = random.GetInt(minx, maxx-d.minRoomSize+1)
			miny = random.GetInt(miny, maxy-d.minRoomSize+1)
			maxx = random.GetInt(minx+d.minRoomSize-1, maxx)
			maxy = random.GetInt(miny+d.minRoomSize-1, maxy)
		}
		// resize the node to fit the room
		//	printf("node %dx%d %dx%d => room %dx%d %dx%d\n",node.x,node.y,node.w,node.h,minx,miny,maxx-minx+1,maxy-miny+1)
		node.X = minx
		node.Y = miny
		node.W = maxx - minx + 1
		node.H = maxy - miny + 1
		// dig the room
		for x = minx; x <= maxx; x++ {
			for y = miny; y <= maxy; y++ {
				d.sampleMap[x][y] = ' '
			}
		}
	} else {
		//	printf("lvl %d %dx%d %dx%d\n",node.level, node.x,node.y,node.w,node.h)
		// resize the node to fit its sons
		left := node.Left()
		right := node.Right()
		node.X = min(left.X, right.X)
		node.Y = min(left.Y, right.Y)
		node.W = max(left.X+left.W, right.X+right.W) - node.X
		node.H = max(left.Y+left.H, right.Y+right.H) - node.Y
		// create a corridor between the two lower nodes
		if node.Horizontal {
			// vertical corridor
			if left.X+left.W-1 < right.X || right.X+right.W-1 < left.X {
				// no overlapping zone. we need a Z shaped corridor
				x1 := random.GetInt(left.X, left.X+left.W-1)
				x2 := random.GetInt(right.X, right.X+right.W-1)
				y := random.GetInt(left.Y+left.H, right.Y)
				d.sampleMap.VLineUp(x1, y-1)
				d.sampleMap.HLine(x1, y, x2)
				d.sampleMap.VLineDown(x2, y+1)
			} else {
				// straight vertical corridor
				minx := max(left.X, right.X)
				maxx := min(left.X+left.W-1, right.X+right.W-1)
				x := random.GetInt(minx, maxx)
				d.sampleMap.VLineDown(x, right.Y)
				d.sampleMap.VLineUp(x, right.Y-1)
			}
		} else {
			// horizontal corridor
			if left.Y+left.H-1 < right.Y || right.Y+right.H-1 < left.Y {
				// no overlapping zone. we need a Z shaped corridor
				y1 := random.GetInt(left.Y, left.Y+left.H-1)
				y2 := random.GetInt(right.Y, right.W+right.H-1)
				x := random.GetInt(left.X+left.W, right.X)
				d.sampleMap.HLineLeft(x-1, y1)
				d.sampleMap.VLine(x, y1, y2)
				d.sampleMap.HLineRight(x+1, y2)
			} else {
				// straight horizontal corridor
				miny := max(left.Y, right.Y)
				maxy := min(left.Y+left.H-1, right.Y+right.H-1)
				y := random.GetInt(miny, maxy)
				d.sampleMap.HLineLeft(right.X-1, y)
				d.sampleMap.HLineRight(right.X, y)
			}
		}
	}
	return true
}


func NewBspDemo() *BspDemo {
	return &BspDemo{
		randomRoom:  false,
		roomWalls:   true,
		bspDepth:    8,
		minRoomSize: 4,
		bsp:         nil,
		generate:    true,
		refresh:     false,
		sampleMap:   SampleMap{},
		darkWall:    Color{0, 0, 100},
		darkGround:  Color{50, 50, 150}}
}



func (self *BspDemo) Render(first bool, key *Key) {
	var x, y int

	if self.generate || self.refresh {
		// dungeon generation
		if self.bsp == nil {
			// create the bsp
			self.bsp = NewBspWithSize(0, 0, DEMO_SCREEN_WIDTH, DEMO_SCREEN_HEIGHT)
		} else {
			// restore the nodes size
			self.bsp.Resize(0, 0, DEMO_SCREEN_WIDTH, DEMO_SCREEN_HEIGHT)
		}
		self.sampleMap.Fill('#')
		if self.generate {
			// build a new random bsp tree
			self.bsp.RemoveSons()
			self.bsp.SplitRecursive(random, self.bspDepth,
				self.minRoomSize+If(self.roomWalls, 1, 0).(int),
				self.minRoomSize+If(self.roomWalls, 1, 0).(int), 1.5, 1.5)
		}
		// create the dungeon from the bsp
		self.bsp.TraverseInvertedLevelOrder(TraverseNode, self)
		self.generate = false
		self.refresh = false
	}
	demoConsole.Clear()
	demoConsole.SetDefaultForeground(COLOR_WHITE)
	demoConsole.PrintEx(1, 1, BKGND_NONE, LEFT,
		"ENTER : rebuild bsp\nSPACE : rebuild dungeon\n+-: bsp depth %d\n: room size %d\n1 : random room size %s",
		self.bspDepth, self.minRoomSize,
		If(self.randomRoom, "ON", "OFF").(string))
	if self.randomRoom {
		demoConsole.PrintEx(1, 6, BKGND_NONE, LEFT, "2 : room walls %s",
			If(self.roomWalls, "ON", "OFF").(string))
	}
	// render the level
	for y = 0; y < DEMO_SCREEN_HEIGHT; y++ {
		for x = 0; x < DEMO_SCREEN_WIDTH; x++ {
			wall := self.sampleMap[x][y] == '#'
			demoConsole.SetCharBackground(x, y, If(wall, self.darkWall, self.darkGround).(Color), BKGND_SET)
		}
	}
	if key.Vk == K_ENTER || key.Vk == K_KPENTER {
		self.generate = true
	} else if key.C == ' ' {
		self.refresh = true
	} else if key.C == '+' {
		self.bspDepth++
		self.generate = true
	} else if key.C == '-' && self.bspDepth > 1 {
		self.bspDepth--
		self.generate = true
	} else if key.C == '*' {
		self.minRoomSize++
		self.generate = true
	} else if key.C == '/' && self.minRoomSize > 2 {
		self.minRoomSize--
		self.generate = true
	} else if key.C == '1' || key.Vk == K_1 || key.Vk == K_KP1 {
		self.randomRoom = !self.randomRoom
		if !self.randomRoom {
			self.roomWalls = true
		}
		self.refresh = true
	} else if key.C == '2' || key.Vk == K_2 || key.Vk == K_KP2 {
		self.roomWalls = !self.roomWalls
		self.refresh = true
	}
}


///
///
///
type ImageDemo struct {
	img, circle *Image
	blue        Color
	green       Color
}

func NewImageDemo() *ImageDemo {
	result := &ImageDemo{}
	result.blue = Color{0, 0, 255}
	result.green = Color{0, 255, 0}
	return result
}


func (self *ImageDemo) Render(first bool, key *Key) {
	var x, y, scalex, scaley, angle float32
	var elapsed int
	if self.img == nil {
		self.img = LoadImage("data/img/skull.png")
		self.img.SetKeyColor(COLOR_BLACK)
		self.circle = LoadImage("data/img/circle.png")
	}
	if first {
		SysSetFps(30) // limited to 30 fps
	}
	demoConsole.SetDefaultBackground(COLOR_BLACK)
	demoConsole.Clear()
	x = float32(float32(DEMO_SCREEN_WIDTH)/2 + cos(float32(SysElapsedSeconds()))*10.0)
	y = (float32)(float32(DEMO_SCREEN_HEIGHT) / 2)
	scalex = float32(0.2 + 1.8*(1.0+cos(float32(SysElapsedSeconds())/2))/2.0)
	scaley = scalex
	angle = float32(SysElapsedSeconds())
	elapsed = int(float32(SysElapsedMilliseconds()) / 2000)
	if elapsed&1 != 0 {
		// split the color channels of circle.png
		// the red channel
		demoConsole.SetDefaultBackground(COLOR_RED)
		demoConsole.Rect(0, 3, 15, 15, false, BKGND_SET)
		self.circle.BlitRect(demoConsole, 0, 3, -1, -1, BKGND_MULTIPLY)
		// the green channel
		demoConsole.SetDefaultBackground(self.green)
		demoConsole.Rect(15, 3, 15, 15, false, BKGND_SET)
		self.circle.BlitRect(demoConsole, 15, 3, -1, -1, BKGND_MULTIPLY)
		// the blue channel
		demoConsole.SetDefaultBackground(self.blue)
		demoConsole.Rect(30, 3, 15, 15, false, BKGND_SET)
		self.circle.BlitRect(demoConsole, 30, 3, -1, -1, BKGND_MULTIPLY)
	} else {
		// render circle.png with normal blitting
		self.circle.BlitRect(demoConsole, 0, 3, -1, -1, BKGND_SET)
		self.circle.BlitRect(demoConsole, 15, 3, -1, -1, BKGND_SET)
		self.circle.BlitRect(demoConsole, 30, 3, -1, -1, BKGND_SET)
	}
	self.img.Blit(demoConsole, x, y, BKGND_SET, scalex, scaley, angle)
}


///
///
/// Mouse demo
//


type MouseDemo struct {
	lbut, rbut, mbut bool
	chars            [DEMO_SCREEN_WIDTH][DEMO_SCREEN_HEIGHT]int
	secondary        IConsole
	asciiArt         IAsciiArt
}

func NewMouseDemo() *MouseDemo {
	result := &MouseDemo{
		lbut: false,
		rbut: false,
		mbut: false}
	for x := 0; x < DEMO_SCREEN_WIDTH; x++ {
		for y := 0; y < DEMO_SCREEN_HEIGHT; y++ {
			result.chars[x][y] = random.GetInt('0', 'z')
		}
	}
	return result
}


func (self *MouseDemo) Render(first bool, key *Key) {
	var mouse Mouse
	if aas == nil {
		aas = NewAsciiArtGallery()
	}

	if first {
		demoConsole.SetDefaultBackground(COLOR_BLACK)
		demoConsole.SetDefaultForeground(COLOR_WHITE)
		demoConsole.Clear()
		MouseMove(320, 200)
		MouseShowCursor(true)
		SysSetFps(30) // limited to 30 fps

		self.secondary = NewConsole(38, 16)
	}

	if self.asciiArt == nil {
		self.asciiArt = aas.randomAsciiArt()
	}
	self.asciiArt.Render(demoConsole, INSTANT, first)

	mouse = MouseGetStatus()

	for x := 0; x < DEMO_SCREEN_WIDTH; x++ {
		for y := 0; y < DEMO_SCREEN_HEIGHT; y++ {
			cx := abs(x + DEMO_SCREEN_X - mouse.Cx)
			cy := abs(y + DEMO_SCREEN_Y - mouse.Cy)

			fore := demoConsole.GetCharForeground(x, y)

			dist := sqrt(float32((sqr(cx)) + sqr(cy)))
			dist = minf(float32(DEMO_SCREEN_WIDTH), dist)

			fore = fore.Lighten(0.6 * sqrf((1 - dist/DEMO_SCREEN_WIDTH)))

			demoConsole.SetCharForeground(x, y, fore)
		}
	}

	// get colors and chars under cursor
	var char, foreColor, backColor string

	if mouse.Cx >= DEMO_SCREEN_X && mouse.Cx-DEMO_SCREEN_X < DEMO_SCREEN_WIDTH &&
		mouse.Cy >= DEMO_SCREEN_Y && mouse.Cy-DEMO_SCREEN_Y < DEMO_SCREEN_HEIGHT {
		cx := abs(DEMO_SCREEN_X - mouse.Cx)
		cy := abs(DEMO_SCREEN_Y - mouse.Cy)
		char = fmt.Sprintf("%s (%d)", string(demoConsole.GetChar(cx, cy)), demoConsole.GetChar(cx, cy))
		foreColor = fmt.Sprintf("%d", demoConsole.GetCharForeground(cx, cy))
		backColor = fmt.Sprintf("%d", demoConsole.GetCharBackground(cx, cy))
	}

	if mouse.LButtonPressed {
		self.lbut = !self.lbut
	}
	if mouse.RButtonPressed {
		self.rbut = !self.rbut
	}
	if mouse.MButtonPressed {
		self.mbut = !self.mbut
	}
	self.secondary.SetDefaultForeground(COLOR_WHITE)
	self.secondary.SetDefaultBackground(COLOR_BLACK)
	self.secondary.Clear()

	self.secondary.PrintEx(1, 1, BKGND_NONE, LEFT,
		`Mouse position : %4dx%4d
Mouse cell     : %4dx%4d
Mouse movement : %4dx%4d
Left button    : %s (toggle %s)
Right button   : %s (toggle %s)
Middle button  : %s (toggle %s)

Fore color     : %s
Back color     : %s
Char           : %s

1 : Hide cursor
2 : Show cursor
`,
		mouse.X, mouse.Y,
		mouse.Cx, mouse.Cy,
		mouse.Dx, mouse.Dy,
		If(mouse.LButton, " ON", "OFF").(string), If(self.lbut, " ON", "OFF").(string),
		If(mouse.RButton, " ON", "OFF").(string), If(self.rbut, " ON", "OFF").(string),
		If(mouse.MButton, " ON", "OFF").(string), If(self.mbut, " ON", "OFF").(string),
		foreColor, backColor, char)

	if key.C == '1' {
		MouseShowCursor(false)
	} else if key.C == '2' {
		MouseShowCursor(true)
	}
	// restore the initial screen
	self.secondary.Blit(0, 0, self.secondary.GetWidth(), self.secondary.GetHeight(),
		demoConsole, 4, 4, 0.8, 0.8)

}
//


///
/// Name
///
type NameDemo struct {
	nbSets int
	curSet int
	delay  float32
	sets   []string
	names  vector.StringVector
}

func NewNameDemo() *NameDemo {
	return &NameDemo{
		names: vector.StringVector{}}
}


func (self *NameDemo) Render(first bool, key *Key) {
	if len(self.sets) == 0 {
		var files []string
		files = SysGetDirectoryContent("data/namegen", "*.cfg")
		// parse all the files
		for _, f := range files {
			NamegenParse("data/namegen/"+f, random)
		}
		// get the sets list
		self.sets = NamegenGetSets()
		self.nbSets = len(self.sets)
	}
	if first {
		SysSetFps(30) // limited to 30 fps
	}

	for self.names.Len() >= DEMO_SCREEN_HEIGHT-10 {
		self.names.Delete(0)
	}

	demoConsole.Clear()
	demoConsole.SetDefaultBackground(COLOR_BLACK)
	demoConsole.SetDefaultForeground(COLOR_DARK_GREEN)
	demoConsole.PrintEx(1, 1, BKGND_NONE, LEFT, "%s\n\n+ : next generator\n- : prev generator",
		self.sets[self.curSet])
	for i, name := range self.names {
		if len(name) < DEMO_SCREEN_WIDTH {
			demoConsole.PrintEx(DEMO_SCREEN_WIDTH-2, 2+i, BKGND_NONE, RIGHT, name)
		}
	}

	self.delay += SysGetLastFrameLength()
	if self.delay >= 0.5 {
		self.delay -= 0.5
		// add a new name to the list
		self.names.Push(NamegenGenerate(self.sets[self.curSet]))
	}
	if key.C == '+' {
		self.curSet++
		if self.curSet == self.nbSets {
			self.curSet = 0
		}
		self.names.Push("======")
	} else if key.C == '-' {
		self.curSet--
		if self.curSet < 0 {
			self.curSet = self.nbSets - 1
		}
		self.names.Push("======")
	}
}

//
// Parser demo
//
const EXAMPLE_FILE = "data/cfg/sample.cfg"

type ParserDemo struct {
	text string
}

func NewParserDemo() *ParserDemo {
	return &ParserDemo{}
}


func parse(fname string) string {
	statesList := []string{"hungry", "very hunger", "starving"}

	p := NewParser()
	ps := p.RegisterStruct("item_type")
	ps.AddProperty("cost", TYPE_INT, true)
	ps.AddProperty("weight", TYPE_FLOAT, true)
	ps.AddProperty("deal_damage", TYPE_BOOL, true)
	ps.AddProperty("damages", TYPE_DICE, true)
	ps.AddProperty("damaged_color", TYPE_COLOR, true)
	ps.AddProperty("damage_type", TYPE_STRING, true)
	ps.AddListProperty("features", TYPE_STRING, true)
	ps.AddListProperty("versions", TYPE_INT, true)
	ps.AddListProperty("advances", TYPE_FLOAT, true)
	ps.AddFlag("abstract")


	ps.AddValueList("states", statesList, true)

	ps = p.RegisterStruct("video")
	ps.AddProperty("mode", TYPE_STRING, true)
	ps.AddProperty("fullscreen", TYPE_BOOL, true)

	ps_input := p.RegisterStruct("input")
	ps_mouse := p.RegisterStruct("mouse")
	ps_mouse.AddProperty("sensibility", TYPE_FLOAT, true)

	ps_input.AddStructure(ps_mouse)

	props := p.Run(EXAMPLE_FILE)

	text := ""
	for _, p := range props {
		text += fmt.Sprintf("%s: %v\n", p.Name, p.Value)
	}
	return text
}

func (self *ParserDemo) Render(first bool, key *Key) {
	if first {
		input, err := ioutil.ReadFile(EXAMPLE_FILE)
		if err != nil {
			panic(fmt.Sprintf("Problem reading file %s", EXAMPLE_FILE))
		}
		output := parse(EXAMPLE_FILE)

		self.text = "Input file:\n\n"
		self.text += string(input)
		self.text += "\n\n"
		self.text += "Parsed output:\n\n"
		self.text += output
	}
	demoConsole.SetDefaultBackground(COLOR_BLACK)
	demoConsole.SetDefaultForeground(COLOR_DARK_GREEN)
	demoConsole.Clear()

	demoConsole.PrintRectEx(1, 1, demoConsole.GetWidth()-1, demoConsole.GetHeight()-1, BKGND_SET, LEFT, self.text)

}


//
//
//

const (
	SAMPLE_TRUE_COLORS = iota
	SAMPLE_MOUSE_SUPPORT
	SAMPLE_OFFSCREEN_CONSOLE
	SAMPLE_LINE_DRAWING
	SAMPLE_NOISE
	SAMPLE_PATH_FINDING
	SAMPLE_FIELD_OF_VIEW
	SAMPLE_BSP_TOOLKIT
	SAMPLE_IMAGE_TOOLKIT
	SAMPLE_NAME_GENERATOR
	SAMPLE_PARSER_DEMO
	SAMPLE_ASCII_SLIDESHOW
)

func Initialize() {

	// change dir to program dir 
	program := os.Args[0]
	dir, _ := path.Split(program)
	os.Chdir(dir)

	random = GetRandomInstance()
	samples = []Sample{
		Sample{name: "True colors       ", demo: NewColorsDemo()},
		Sample{name: "Mouse support     ", demo: NewMouseDemo()},
		Sample{name: "Offscreen console ", demo: NewOffscreenDemo()},
		Sample{name: "Line drawing      ", demo: NewLinesDemo()},
		Sample{name: "Noise             ", demo: NewNoiseDemo()},
		Sample{name: "Path finding      ", demo: NewPathDemo()},
		Sample{name: "Field of view     ", demo: NewFovDemo()},
		Sample{name: "Bsp toolkit       ", demo: NewBspDemo()},
		Sample{name: "Image toolkit     ", demo: NewImageDemo()},
		Sample{name: "Name generator    ", demo: NewNameDemo()},
		Sample{name: "Parser demo       ", demo: NewParserDemo()},
		Sample{name: "Ascii slideshow   ", demo: NewAsciiArtDemo()}}
}



func switchDemoCbk(w IWidget, data interface{}) {
	curSample = data.(int)
	first = true
}

func setColors(w IWidget) {
	w.SetDefaultForeground(COLOR_GREY, COLOR_BLACK)
	w.SetDefaultBackground(COLOR_BLACK, COLOR_GREY)
}

func setButtonColors(w IWidget) {
	b := w.(*RadioButton)
	b.SetDefaultForeground(COLOR_GREY, COLOR_GREY.Lighten(0.4))
	b.SetDefaultBackground(COLOR_BLACK, COLOR_BLACK.Lighten(0.4))
	b.SetUseSelectionColor(true)
	b.SetSelectionColor(COLOR_WHITE, COLOR_LIGHT_BLUE)
}

func buildGui() {
	gui = NewGui(guiConsole)
	//gui = NewGui(rootConsole)
	// status bar
	s := gui.NewStatusBarDim(0, 0, GUI_SCREEN_WIDTH, 1)
	setColors(s)
	vbox := gui.NewVBox(0, 2, 1)
	setColors(vbox)
	tools := gui.NewToolBarWithWidth(1, 1, 15, "Demo", "Tools to modify the heightmap")
	tools.SetShouldPrintFrame(false)
	setColors(s)
	for i := range samples {
		samples[i].button = gui.NewRadioButton(samples[i].name, "Show "+strings.Trim(samples[i].name, " ")+" demo", switchDemoCbk, i)
		setButtonColors(samples[i].button)
		tools.AddWidget(samples[i].button)
	}
	vbox.AddWidget(tools)
}

func Run() {

	curSample = SAMPLE_MOUSE_SUPPORT // index of the current sample
	first = true                     // first time we render a sample
	//var i int
	key := Key{Vk: K_NONE, C: 0}
	font := "data/fonts/prestige10x10_gs_tc.png"
	//	font := "terminal2.png"
	nbCharHoriz := 0
	nbCharVertic := 0
	var argn int
	fullscreenWidth := 0
	fullscreenHeight := 0
	fontFlags := FONT_TYPE_GREYSCALE | FONT_LAYOUT_TCOD
	fontNewFlags := 0

	//fullscreen := false
	creditsEnd := false

	//	SysSetFps(25)

	// FONTS = [
	//     ('fonts/terminal10x18.png',
	//      libtcod.FONT_LAYOUT_ASCII_INROW),
	//     ('fonts/terminal8x15.png',
	//      libtcod.FONT_LAYOUT_ASCII_INROW),
	//     ('fonts/terminal8x8.png',
	//      libtcod.FONT_LAYOUT_ASCII_INCOL),
	// ]
	//
	font = "data/fonts/terminal8x8.png"
	fontFlags = FONT_LAYOUT_ASCII_INCOL

	// initialize the rootConsole console (open the game window)
	for argn = 1; argn < flag.NArg(); argn++ {
		if flag.Arg(argn) == "-font" && argn+1 < flag.NArg() {
			argn++
			font = flag.Arg(argn)
			fontFlags = 0
		} else if flag.Arg(argn) == "-font-nb-char" && argn+2 < flag.NArg() {
			argn++
			nbCharHoriz = atoi(flag.Arg(argn))
			argn++
			nbCharVertic = atoi(flag.Arg(argn))
			fontFlags = 0
		} else if flag.Arg(argn) == "-fullscreen-resolution" && argn+2 < flag.NArg() {
			argn++
			fullscreenWidth = atoi(flag.Arg(argn))
			argn++
			fullscreenHeight = atoi(flag.Arg(argn))
		} else if flag.Arg(argn) == "-fullscreen" {
			//fullscreen = true
		} else if flag.Arg(argn) == "-font-in-row" {
			fontFlags = 0
			fontNewFlags |= FONT_LAYOUT_ASCII_INROW
		} else if flag.Arg(argn) == "-font-greyscale" {
			fontFlags = 0
			fontNewFlags |= FONT_TYPE_GREYSCALE
		} else if flag.Arg(argn) == "-font-tcod" {
			fontFlags = 0
			fontNewFlags |= FONT_LAYOUT_TCOD
		} else if flag.Arg(argn) == "-help" {
			fmt.Print("options :\n")
			fmt.Print("-font <filename> : use a custom font\n")
			fmt.Print("-font-nb-char <nbCharHoriz> <nbCharVertic> : number of characters in the font\n")
			fmt.Print("-font-in-row : the font layout is in row instead of columns\n")
			fmt.Print("-font-tcod : the font uses TCOD layout instead of ASCII\n")
			fmt.Print("-font-greyscale : antialiased font using greyscale bitmap\n")
			fmt.Print("-fullscreen : start in fullscreen\n")
			fmt.Print("-fullscreen-resolution <screen_width> <screen_height> : force fullscreen resolution\n")
			os.Exit(0)
		} else {
			argn++ // ignore parameter
		}
	}

	if fontFlags == 0 {
		fontFlags = fontNewFlags
	}


	if fullscreenWidth > 0 {
		SysForceFullscreenResolution(fullscreenWidth, fullscreenHeight)
	}

	rootConsole = NewRootConsoleWithFont(SCREEN_WIDTH, SCREEN_HEIGHT, "Go demo", false, font, fontFlags, nbCharHoriz,
	nbCharVertic, RENDERER_SDL)

	guiConsole = NewConsole(GUI_SCREEN_WIDTH, GUI_SCREEN_HEIGHT)

	// initialize the offscreen console for the samples
	demoConsole = NewConsole(DEMO_SCREEN_WIDTH, DEMO_SCREEN_HEIGHT)

	buildGui()

	for {


		rootConsole.Clear() // THIS IS CRITICAL !!
		guiConsole.SetDefaultBackground(COLOR_BLACK)
		guiConsole.SetDefaultForeground(COLOR_WHITE)
		guiConsole.Clear()
		//		gui.UpdateWidgets(key)
		gui.RenderWidgets()

		// render gui
		guiConsole.Blit(0, 0, GUI_SCREEN_WIDTH, GUI_SCREEN_HEIGHT, rootConsole, 0, 0, 1.0, 1.0)

		if !creditsEnd {
			creditsEnd = rootConsole.RenderCredits(2, 49, false)
		}

		// print the help message
		rootConsole.SetDefaultForeground(COLOR_GREY)
		rootConsole.PrintEx(2, 19, BKGND_NONE, LEFT, "%c%c: select demo", CHAR_ARROW_N, CHAR_ARROW_S)
		rootConsole.PrintEx(2, 20, BKGND_NONE, LEFT, "alt-enter: %s",
			If(rootConsole.IsFullscreen(), "window ", "fullscreen ").(string))

		rootConsole.PrintEx(2, 60, BKGND_NONE, LEFT, "last frame: %4d ms", int(SysGetLastFrameLength()*1000))
		rootConsole.PrintEx(2, 61, BKGND_NONE, LEFT, "fps: %10d fps", SysGetFps())
		rootConsole.PrintEx(2, 62, BKGND_NONE, LEFT, "elapsed: %8dms", SysElapsedMilliseconds())
		// render current sample
		samples[curSample].demo.Render(first, &key)
		first = false
		//
		// blit the sample console on the rootConsole console
		demoConsole.Blit(0, 0, DEMO_SCREEN_WIDTH, DEMO_SCREEN_HEIGHT, // the source console & zone to blit
			rootConsole, DEMO_SCREEN_X, DEMO_SCREEN_Y, // the destination console & position
			1.0, 1.0) // alpha coefs  )

		// update the game screen
		rootConsole.Flush()

		// did the user hit a key ?
		key = rootConsole.CheckForKeypress(KEY_PRESSED)
		gui.UpdateWidgets(key)
		if key.Vk == K_DOWN {
			// down arrow : next sample
			curSample = (curSample + 1) % len(samples)
			first = true
		} else if key.Vk == K_UP {
			// up arrow : previous sample
			curSample--
			if curSample < 0 {
				curSample = len(samples) - 1
			}
			first = true
		} else if key.Vk == K_ENTER && key.LAlt {
			// ALT-ENTER : switch fullscreen
			rootConsole.SetFullscreen(!rootConsole.IsFullscreen())
		} else if key.Vk == K_PRINTSCREEN {
			// save screenshot
			SysSaveScreenshot()
		}


		samples[curSample].button.Select()

		if rootConsole.IsWindowClosed() {
			break
		}
	}
}


func main() {
	Initialize()
	Run()
}
