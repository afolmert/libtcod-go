package main

import (
	"container/vector"
	"compress/zlib"
	//	"fmt"
	"gob"
	"io/ioutil"
	"os"
	"path"
	"strings"
	. "tcod"
)

type Transition uint8

type MoveTransition uint8

// Transitions
const (
	INSTANT = iota
	FADE_IN
	MOVE_IN
	RANDOMIZE
	NB_TRANSITIONS
)

// Move transitions
const (
	MOVE_UP = iota
	MOVE_DOWN
	MOVE_LEFT
	MOVE_RIGHT
	NB_MOVE_TRANSITIONS
)

const TRANSITION_TIME = 700 // time in milliseconds

type IAsciiArt interface {
	Render(console IConsole, transition Transition, first bool)
}

type AsciiArtGallery struct {
	secondary      IConsole
	firstTime      int // time in milliseconds when first render occurred
	moveTransition MoveTransition
	arts           []IAsciiArt
	artfiles       vector.StringVector
	last           int
	loadedArtsChan chan bool
	loadedArts     bool
}

type AsciiArt struct {
	name             string
	chars            []string
	colors           [][]Color
	width, height    int
	offsetX, offsetY int
}


var asciiArts []IAsciiArt

var aas *AsciiArtGallery

func NewAsciiArtGallery() *AsciiArtGallery {
	result := &AsciiArtGallery{}
	result.loadedArts = false
	result.loadedArtsChan = make(chan bool)
	result.goLoadAsciiArts()
	return result
}

func (self *AsciiArtGallery) isArtsLoaded() bool {
	// if arts not loaded then wait for loading
	if !self.loadedArts {
		if _, ok := <-self.loadedArtsChan; ok {
			self.loadedArts = true
		}
	}
	return self.loadedArts
}

// loads ascii arts from data directory
func (self *AsciiArtGallery) goLoadAsciiArts() {
	self.loadedArts = false
	go func() {
		files := vector.StringVector{}
		entries, err := ioutil.ReadDir("data/ascii")
		if err != nil {
			panic("Cannot read files from data/ascii")
		}
		for _, f := range entries {
			if strings.HasSuffix(f.Name, ".dat") {
				files.Push(path.Join("data/ascii", f.Name))
			}
		}

		self.arts = make([]IAsciiArt, files.Len())
		self.artfiles = files
		self.loadedArtsChan <- true
	}()
}

func (self *AsciiArtGallery) getAsciiArt(id int) IAsciiArt {
	var result IAsciiArt
	if self.arts[id] == nil {
		art, err := NewAsciiArtFromFile(self.artfiles[id])
		if err != nil {
			panic("Cannot load ascii file " + self.artfiles[id])
		}
		self.arts[id] = art
	}
	result = self.arts[id]
	return result
}

// select random ASCII art but not the one before
func (self *AsciiArtGallery) randomAsciiArt() IAsciiArt {
	if self.isArtsLoaded() {
		var id int

		// find id to return
		if len(self.arts) == 0 {
			panic("No arts found!")
		} else if len(self.arts) == 1 {
			id = 0
		} else {
			for {
				id = random.GetInt(0, len(self.arts)-1)
				if self.last != id {
					break
				}
			}
			self.last = id
		}
		return self.getAsciiArt(id)

	}
	return nil
}


func NewAsciiArtFromFile(fname string) (art IAsciiArt, err os.Error) {

	fin, err := os.Open(fname, os.O_RDONLY, 0664)
	if err != nil {
		return nil, err
	}
	i, err := zlib.NewReader(fin)
	defer i.Close()
	if err != nil {
		return nil, err
	}
	d := gob.NewDecoder(i)
	a := AsciiArt{}
	d.Decode(&a)
	return &a, nil

}


func (self *AsciiArt) PutChar(console IConsole, x, y int) {
	color := self.colors[y+self.offsetY][x+self.offsetX]
	// char := int(self.chars[y+self.offsetY][x+self.offsetX])
	console.PutChar(x, y, 'x', BKGND_SET)
	console.SetFore(x, y, color)
	console.SetBack(x, y, COLOR_BLACK, BKGND_SET)

}

func (self *AsciiArt) Draw(console IConsole) {
	console.SetBackgroundColor(COLOR_BLACK)
	console.Clear()
	for x := 0; x < min(console.GetWidth(), self.width); x++ {
		for y := 0; y < min(console.GetHeight(), self.height); y++ {
			self.PutChar(console, x, y)
		}
	}
}


func (self *AsciiArt) Render(console IConsole, transition Transition, first bool) {
	if first {
		aas.firstTime = int(SysElapsedMilliseconds())
	}
	elapsed := int(SysElapsedMilliseconds()) - aas.firstTime

	switch transition {

	case INSTANT:
		{
			self.Draw(console)

		}
	case RANDOMIZE:
		{
			if first {
				console.SetBackgroundColor(COLOR_BLACK)
			}
			for i := 0; i < 500; i++ {
				x := random.GetInt(0, console.GetWidth()-1)
				y := random.GetInt(0, console.GetHeight()-1)

				// if in our picture then draw it
				if x < self.width && y < self.height {
					self.PutChar(console, x, y)
				} else {
					// else draw black background
					console.PutChar(x, y, ' ', BKGND_SET)
					console.SetFore(x, y, COLOR_BLACK)
					console.SetBack(x, y, COLOR_BLACK, BKGND_SET)
				}
			}

		}
	case FADE_IN:
		{
			if first {
				if aas.secondary != nil {
					aas.secondary.Delete()
				}
				aas.secondary = NewConsole(console.GetWidth(), console.GetHeight())
				console.Blit(0, 0, console.GetWidth(), console.GetHeight(), aas.secondary, 0, 0, 1.0, 1.0)
			}
			console.SetBackgroundColor(COLOR_BLACK)
			console.Clear()
			// fade out old
			if elapsed <= TRANSITION_TIME {
				alpha := 1.0 - float(elapsed)/TRANSITION_TIME
				aas.secondary.Blit(0, 0, console.GetWidth(), console.GetHeight(), console, 0, 0, alpha, alpha)
				// fade in new
			} else {
				self.Draw(aas.secondary)
				alpha := minf((float(elapsed)-TRANSITION_TIME)/TRANSITION_TIME, 1)
				aas.secondary.Blit(0, 0, console.GetWidth(), console.GetHeight(), console, 0, 0, alpha, alpha)
			}

		}
	case MOVE_IN:
		{
			if first {
				if aas.secondary != nil {
					aas.secondary.Delete()
				}
				aas.secondary = NewConsole(console.GetWidth(), console.GetHeight())
				aas.moveTransition = MoveTransition(random.GetInt(0, NB_MOVE_TRANSITIONS-1))
				self.Draw(aas.secondary)
			}
			widthPart := max(1, min(console.GetWidth(), int(float(elapsed)/TRANSITION_TIME*float(console.GetWidth()))))
			heightPart := max(1, min(console.GetHeight(), int(float(elapsed)/TRANSITION_TIME*float(console.GetHeight()))))

			switch aas.moveTransition {
			case MOVE_LEFT:
				aas.secondary.Blit(
					0, 0, widthPart, console.GetHeight(),
					console, 0, 0, 1.0, 1.0)
			case MOVE_RIGHT:
				aas.secondary.Blit(
					console.GetWidth()-widthPart, 0, widthPart, console.GetHeight(),
					console, console.GetWidth()-widthPart, 0, 1.0, 1.0)
			case MOVE_DOWN:
				aas.secondary.Blit(0, 0, console.GetWidth(), heightPart,
					console, 0, 0, 1.0, 1.0)
			case MOVE_UP:
				aas.secondary.Blit(
					0, console.GetHeight()-heightPart, console.GetWidth(), heightPart,
					console, 0, console.GetHeight()-heightPart, 1.0, 1.0)

			}

		}
	default:
		{
			// do nothing
			self.Render(console, FADE_IN, first)
			// TODO
			// rain/driple -> transitions -> similar to MOVE_IN transitions but for each streak sepreately
			// with different speed
			// fade in
			// have 2 consoles -> save current as screenshot
			// // and fill then new one and blit it
			// move in -> <-
			//  rain driple like matrix effect
		}
	}
}
