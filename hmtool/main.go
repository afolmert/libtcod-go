/*
libtcod height map tool

Actually self is currently more a curiosity than a real useful tool.
It allows you to tweak a heightmap by hand and generate the corresponding C,C++ or python code.

The heightmap tool source code is public domain. Do whatever you want with it.
*/
package main

import (
	"container/list"
	"fmt"
	"github.com/afolmert/libtcod-go/tcod"
	"io/ioutil"
	"math"
	"os"
	"path"
	"strconv"
)

//
//
// Operations
//
//
const HM_WIDTH = 100
const HM_HEIGHT = 80

// Operations
//
//
type OpType int

// OpType
const (
	NORM = iota
	ADDFBM
	SCALEFBM
	ADDHILL
	ADDLEVEL
	SMOOTH
	RAIN
	NOISELERP
	VORONOI
)

type CodeType int

// CodeType
const (
	C = iota
	CPP
	PY
	GO
	NB_CODE
)

var HEADER1 []string = []string{
	// C header
	"#include <stdlib.h>\n" +
		"#include \"libtcod.h\"\n" +
		"// size of the heightmap\n" +
		"#define HM_WIDTH 100\n" +
		"#define HM_HEIGHT 80\n",
	// CPP header
	"#include \"libtcod.hpp\"\n" +
		"// size of the heightmap\n" +
		"#define HM_WIDTH 100\n" +
		"#define HM_HEIGHT 80\n",
	// PY header
	"#!/usr/bin/python\n" +
		"import math\n" +
		"import libtcodpy as libtcod\n" +
		"# size of the heightmap\n" +
		"HM_WIDTH=100\n" +
		"HM_HEIGHT=80\n",
	// GO header
	"package main\n" +
		"import \"tcod\"\n" +
		"import \"math\"\n" +
		"// size of the heightmap\n" +
		"const HM_WIDTH=100\n" +
		"const HM_HEIGHT=80\n\n" +
		"func min(a, b float32) float32 {\n" +
		"  if a < b { \n" +
		"	  return a  \n" +
		"  } \n" +
		"  return b  \n" +
		"} \n\n " +
		"  \n " +
		"func sin(f float32) float32 { \n " +
		"	return float32(math.Sin(float3264(f))) \n " +
		"} \n " +
		" \n " +
		"func cos(f float32) float32 { \n " +
		"	return float32(math.Cos(float3264(f))) \n " +
		"} \n " +
		" \n ",
}

var HEADER2 []string = []string{
	// C header 2
	"// function building the heightmap\n" +
		"void buildMap(TCOD_heightmap_t *hm) {\n",
	// CPP header 2
	"// function building the heightmap\n" +
		"void buildMap(TCODHeightMap *hm) {\n",
	// PY header 2
	"# function building the heightmap\n" +
		"def buildMap(hm) :\n",
	// GO header 2
	"// function building the heightmap\n" +
		"func buildMap(hm *tcod.HeightMap) {\n",
}

var FOOTER1 []string = []string{
	// C footer
	"}\n" +
		"// test code to print the heightmap\n" +
		"// to compile this file on Linux :\n" +
		"//  gcc hm.c -o hm -I include/ -L . -ltcod\n" +
		"// to compile this file on Windows/mingw32 :\n" +
		"//  gcc hm.c -o hm.exe -I include/ -L lib -ltcod-mingw\n" +
		"int main(int argc, char *argv[]) {\n" +
		"\tint x,y;\n" +
		"\tTCOD_heightmap_t *hm=TCOD_heightmap_new(HM_WIDTH,HM_HEIGHT);\n",
	// CPP footer
	"}\n" +
		"// test code to print the heightmap\n" +
		"// to compile this file on Linux :\n" +
		"//  g++ hm.cpp -o hm -I include/ -L . -ltcod -ltcod++\n" +
		"// to compile this file on Windows/mingw32 :\n" +
		"//  g++ hm.cpp -o hm.exe -I include/ -L lib -ltcod-mingw\n" +
		"int main(int argc, char *argv[]) {\n" +
		"\tTCODHeightMap hm(HM_WIDTH,HM_HEIGHT);\n" +
		"\tbuildMap(&hm);\n" +
		"\tTCODConsole::initRoot(HM_WIDTH,HM_HEIGHT,\"height map test\",false);\n" +
		"\tfor (int x=0; x < HM_WIDTH; x ++ ) {\n" +
		"\t\tfor (int y=0;y < HM_HEIGHT; y++ ) {\n" +
		"\t\t\tfloat32 z = hm.getValue(x,y);\n" +
		"\t\t\tuint8 val=(uint8)(z*255);\n" +
		"\t\t\tTCODColor c(val,val,val);\n" +
		"\t\t\tTCODConsole::root->setBack(x,y,c);\n" +
		"\t\t}\n" +
		"\t}\n" +
		"\tTCODConsole::root->flush();\n" +
		"\tTCODConsole::waitForKeypress(true);\n" +
		"\treturn 0;\n" +
		"}\n",
	// PY footer
	"# test code to print the heightmap\n" +
		"hm=libtcod.heightmap_new(HM_WIDTH,HM_HEIGHT)\n" +
		"buildMap(hm)\n" +
		"libtcod.console_init_root(HM_WIDTH,HM_HEIGHT,\"height map test\",False)\n" +
		"for x in range(HM_WIDTH) :\n" +
		"    for y in range(HM_HEIGHT) :\n" +
		"        z = libtcod.heightmap_get_value(hm,x,y)\n" +
		"        val=int(z*255) & 0xFF\n" +
		"        c=libtcod.Color(val,val,val)\n" +
		"        libtcod.console_set_back(None,x,y,c,libtcod.BKGND_SET)\n" +
		"libtcod.console_flush()\n" +
		"libtcod.console_wait_for_keypress(True)\n",
	// GO footer
	"}\n" +
		"// test code to print the heightmap\n" +
		"// to compile self file on Linux :\n" +
		"//  8g hm.go && 8l hm.8 \n" +
		"func main() {\n" +
		"\thm := tcod.NewHeightMap(HM_WIDTH,HM_HEIGHT)\n" +
		"\tbuildMap(hm)\n" +
		"\troot := tcod.NewRootConsole(HM_WIDTH,HM_HEIGHT,\"height map test\",false)\n" +
		"\tfor x := 0; x < HM_WIDTH; x++ {\n" +
		"\t\tfor y := 0; y < HM_HEIGHT; y++ {\n" +
		"\t\t\tz := hm.GetValue(x,y)\n" +
		"\t\t\tval:=uint8(z*255)\n" +
		"\t\t\tc:= tcod.Color{val,val,val}\n" +
		"\t\t\troot.SetBack(x,y,c,tcod.BKGND_SET)\n" +
		"\t\t}\n" +
		"\t}\n" +
		"\troot.Flush()\n" +
		"\troot.WaitForKeypress(true)\n" +
		"}\n",
}

var FOOTER2 []string = []string{
	// C footer
	"\tbuildMap(hm);\n" +
		"\tTCOD_console_init_root(HM_WIDTH,HM_HEIGHT,\"height map test\",false);\n" +
		"\tfor (x=0; x < HM_WIDTH; x ++ ) {\n" +
		"\t\tfor (y=0;y < HM_HEIGHT; y++ ) {\n" +
		"\t\t\tfloat32 z = TCOD_heightmap_get_value(hm,x,y);\n" +
		"\t\t\tuint8 val=(uint8)(z*255);\n" +
		"\t\t\tTCOD_color_t c={val,val,val};\n" +
		"\t\t\tTCOD_console_set_back(NULL,x,y,c,TCOD_BKGND_SET);\n" +
		"\t\t}\n" +
		"\t}\n" +
		"\tTCOD_console_flush();\n" +
		"\tTCOD_console_wait_for_keypress(true);\n" +
		"\treturn 0;\n" +
		"}\n",
	// CPP footer
	"",
	// PY footer
	"",
	// GO footer
	"",
}

var hm, hmold *tcod.HeightMap
var noise *tcod.Noise
var rnd, backupRnd *tcod.Random

var gui *tcod.Gui
var root *tcod.RootConsole
var guicon *tcod.Console

var greyscale bool = false
var slope bool = false
var normal bool = false
var isNormalized bool = true
var oldNormalized bool = true

var msg string
var msgDelay float32 = 0.0
var hillRadius float32 = 0.1
var hillVariation float32 = 0.5
var addFbmDelta float32 = 0.0
var scaleFbmDelta float32 = 0.0
var seed uint32 = 0xdeadbeef

var sandHeight float32 = 0.12
var grassHeight float32 = 0.315
var snowHeight float32 = 0.785

var landMassLabel *tcod.Label
var minZLabel *tcod.Label
var maxZLabel *tcod.Label
var seedLabel *tcod.Label

var mapmin float32 = 0.0
var mapmax float32 = 0.0
var oldmapmin float32 = 0.0
var oldmapmax float32 = 0.0

var params *tcod.ToolBar
var history *tcod.ToolBar
var colorMapGui *tcod.ToolBar

var voronoiCoef []float32 = []float32{-1.0, 1.0}

/* light 3x3 smoothing kernel :
1  2 1
2 20 2
1  2 1
*/
var smoothKernelSize int = 9

var smoothKernelDx []int = []int{-1, 0, 1, -1, 0, 1, -1, 0, 1}

var smoothKernelDy []int = []int{-1, -1, -1, 0, 0, 0, 1, 1, 1}

var smoothKernelWeight []float32 = []float32{1, 2, 1, 2, 20, 2, 1, 2, 1}

var mapGradient []tcod.Color = make([]tcod.Color, 256)

const MAX_COLOR_KEY = 10

// TCOD's land color map
var keyIndex []int = []int{0,
	int(sandHeight * 255),
	int(sandHeight*255) + 4,
	int(grassHeight * 255),
	int(grassHeight*255) + 10,
	int(snowHeight * 255),
	int(snowHeight*255) + 10,
	255}

var keyColor []tcod.Color = []tcod.Color{
	tcod.Color{0, 0, 50},      // deep water
	tcod.Color{30, 30, 170},   // water-sand transition
	tcod.Color{114, 150, 71},  // sand
	tcod.Color{80, 120, 10},   // sand-grass transition
	tcod.Color{17, 109, 7},    // grass
	tcod.Color{120, 220, 120}, // grass-snow transisiton
	tcod.Color{208, 208, 239}, // snow
	tcod.Color{255, 255, 255}}

var keyImages [MAX_COLOR_KEY]*tcod.ImageWidget

var nbColorKeys int = 8

var backupMap *tcod.HeightMap

var operations *OperationStatic = newOperationStatic()

//
// misc functions
//

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

func sin(f float32) float32 {
	return float32(math.Sin(float64(f)))
}

func cos(f float32) float32 {
	return float32(math.Cos(float64(f)))
}

func listRemove(l list.List, val interface{}) {
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value == val {
			l.Remove(e)
			break
		}
	}
}

func stringSliceContains(s []string, val string) bool {
	for _, e := range s {
		if e == val {
			return true
		}

	}
	return false
}

func peekOperation(l list.List) IOperation {
	if l.Back() != nil {
		return l.Back().Value.(IOperation)
	}
	return nil
}

//
//
// Operations
//

type OperationStatic struct {
	names       []string
	tips        []string
	list        list.List // the list of operation applied since the last clear
	needsRandom bool      // we need a random number generator
	needsNoise  bool      // we need a 2D noise
	currentOp   IOperation
	codebuf     string            // generated code buffer
	initCode    [NB_CODE][]string // list of global vars/functions to add to the generated code
}

func newOperationStatic() *OperationStatic {
	result := &OperationStatic{}
	result.list = list.List{}
	result.initCode = [NB_CODE][]string{
		[]string{},
		[]string{},
		[]string{}}
	// must match OpType enum
	result.names = []string{
		"norm",
		"+fbm",
		"*fbm",
		"hill",
		"\x18\x19 z",
		"smooth",
		"rain",
		"lerp fbm",
		"voronoi"}

	result.tips = []string{
		"Normalize heightmap so that values are between 0.0 and 1.0",
		"Add fbm noise to current heightmap",
		"Scale the heightmap by a fbm noise",
		"Add random hills on the heightmap",
		"[+/-] Raise or lower the whole heightmap",
		"Smooth the heightmap",
		"Simulate rain erosion on the heightmap",
		"Lerp between the heightmap and a fbm noise",
		"Add a voronoi diagram to the heightmap"}

	return result
}

// generate the code corresponding to the list of operations
func (self *OperationStatic) buildCode(codeType CodeType) string {
	self.codebuf = ""
	self.addCode(HEADER1[codeType])
	if self.needsRandom || self.needsNoise {
		switch codeType {
		case C:
			self.addCode(fmt.Sprintf("tcod.TCOD_random_t rnd=nil;\n", seed))
		case CPP:
			self.addCode(fmt.Sprintf("TCODRandom *rnd=new TCODRandom(%uU);\n", seed))
		case PY:
			self.addCode(fmt.Sprintf("rnd=libtcod.random_new_from_seed(%u)\n", seed))
		case GO:
			self.addCode(fmt.Sprintf("var rnd=tcod.NewRandomFromSeed(%d)\n", seed))
		default:
		}
	}
	if self.needsNoise {
		switch codeType {
		case C:
			self.addCode("tcod.TCOD_noise_t noise=nil;\n")
		case CPP:
			self.addCode("TCODNoise *noise=new TCODNoise(2,rnd);\n")
		case PY:
			self.addCode("noise=libtcod.noise_new(2,libtcod.NOISE_DEFAULT_HURST,libtcod.NOISE_DEFAULT_LACUNARITY,rnd)\n")
		case GO:
			self.addCode("var noise=tcod.NewNoise(2,rnd)\n")
		default:
		}
	}
	for _, s := range self.initCode[codeType] {
		self.addCode(s)
	}
	self.addCode(HEADER2[codeType])
	for e := self.list.Front(); e != nil; e = e.Next() {
		op := e.Value.(IOperation)
		code := op.getCode(codeType)
		self.addCode(code)
	}
	self.addCode(FOOTER1[codeType])
	if (self.needsRandom || self.needsNoise) && codeType == C {
		self.addCode(fmt.Sprintf("\trnd=tcod.TCOD_random_new_from_seed(%uU);\n", seed))
		if self.needsNoise {
			self.addCode("\tnoise=tcod.TCOD_noise_new(2,tcod.TCOD_NOISE_DEFAULT_HURST,tcod.TCOD_NOISE_DEFAULT_LACUNARITY,rnd);\n")
		}
	}
	self.addCode(FOOTER2[codeType])
	return self.codebuf
}

// remove all operation, clear the heightmap
func (self *OperationStatic) clear() {
	self.list = list.List{}
}

// cancel the last operation
func (self *OperationStatic) cancel() {
	if self.currentOp != nil {
		listRemove(self.list, self.currentOp)
		history.RemoveWidget(self.currentOp.getButton())
		self.currentOp.getButton().UnSelect()
		//delete self.currentOp
		self.currentOp = peekOperation(self.list)
		if self.currentOp != nil {
			self.currentOp.getButton().Select()
			self.currentOp.createParamUi()
		} else {
			params.Clear()
			params.SetVisible(false)
		}
		self.reseed() // replay the whole stack
	}
}

func (self *OperationStatic) reseed() {
	rnd = tcod.NewRandomFromSeed(seed)
	noise = tcod.NewNoise(2, rnd)
	addFbmDelta = 0
	scaleFbmDelta = 0
	hm.Clear()
	for e := self.list.Front(); e != nil; e = e.Next() {
		op := e.Value.(IOperation)
		op.runInternal()
	}

}

func (self *OperationStatic) run(op IOperation) {
	op.runInternal()
}

// add a global variable or a function to the generated code
func (self *OperationStatic) addInitCode(codeType CodeType, code string) {
	if !stringSliceContains(self.initCode[codeType], code) {
		self.initCode[codeType] = append(self.initCode[codeType], code)
	}
}

// add some code to the generated code
func (self *OperationStatic) addCode(code string) {
	self.codebuf += code
}

// run self operation and adds it in the list
func (self *OperationStatic) add(op IOperation) {
	backup()
	op.runInternal()
	if op.addInternal() {
		operations.list.PushBack(op)
		op.createParamUi()
		op.setButton(gui.NewRadioButton(operations.names[op.getOpType()], operations.tips[op.getOpType()], historyCbk, op))
		op.getButton().SetGroup(0)
		history.AddWidget(op.getButton())
		op.getButton().Select()
		operations.currentOp = op
	}
	//else delete self
}

func historyCbk(w tcod.IWidget, data interface{}) {
	op := data.(IOperation)
	op.createParamUi()
	op.getButton().Select()
	operations.currentOp = op
}

//
//
// Operation
//

type IOperation interface {
	runInternal()
	addInternal() bool
	createParamUi()
	getCode(codeType CodeType) string
	getButton() *tcod.RadioButton
	setButton(b *tcod.RadioButton)
	getOpType() OpType
}

type Operation struct {
	opType OpType
	button *tcod.RadioButton // button associated with self operation in history
}

func (self *Operation) initializeOperation(opType OpType) {
	self.opType = opType
}

func (self *Operation) createParamUi() {
	params.Clear()
	params.SetVisible(false)
}

func (self *Operation) historyCbk(w tcod.IWidget, data interface{}) {
	// abstract
}

// actually execute self operation
func (self *Operation) runInternal() {
	// abstract
}

// actually add self operation
func (self *Operation) addInternal() bool {
	return false
}

// the code corresponding to self operation
func (self *Operation) getCode(codeType CodeType) string {
	return ""
}

func (self *Operation) getButton() *tcod.RadioButton {
	return self.button
}

func (self *Operation) setButton(b *tcod.RadioButton) {
	self.button = b
}

func (self *Operation) getOpType() OpType {
	return self.opType
}

//
// Normalize operation
//

// normalize the heightmap
type NormalizeOperation struct {
	Operation
	min, max float32
}

func NewNormalizeOperation(min, max float32) *NormalizeOperation {
	result := &NormalizeOperation{}
	result.initializeNormalizeOperation(NORM, min, max)
	return result
}

func NewNormalizeOperationWithOptions(min, max float32) *NormalizeOperation {
	result := &NormalizeOperation{}
	result.initializeNormalizeOperation(NORM, min, max)
	return result
}

func (self *NormalizeOperation) initializeNormalizeOperation(opType OpType, min, max float32) {
	self.Operation.initializeOperation(opType)
	self.min, self.max = min, max
}

// Normalize
func (self *NormalizeOperation) getCode(codeType CodeType) string {
	switch codeType {
	case C:
		return fmt.Sprintf("\ttcod.TCOD_heightmap_normalize(hm,%g,%g);\n", self.min, self.max)
	case CPP:
		return fmt.Sprintf("\thm.normalize(%g,%g);\n", min, max)
	case PY:
		return fmt.Sprintf("    libtcod.heightmap_normalize(hm,%g,%g)\n", self.min, self.max)
	case GO:
		return fmt.Sprintf("    hm.NormalizeRange(%g,%g)\n", self.min, self.max)
	default:
	}
	return ""
}

func (self *NormalizeOperation) runInternal() {
	hm.NormalizeRange(self.min, self.max)
}

func (self *NormalizeOperation) addInternal() bool {
	prev := peekOperation(operations.list)
	if prev != nil && prev.getOpType() == NORM {
		return false
	}
	return true
}

func normalizeMinValueCbk(w tcod.IWidget, val string, data interface{}) {
	f, err := strconv.ParseFloat(val, 32)
	if err != nil {
		op := data.(*NormalizeOperation)
		if float32(f) < op.max {
			op.min = float32(f)
			if peekOperation(operations.list) == IOperation(op) {
				op.runInternal()
			} else {
				operations.reseed()
			}
		}
	}
}

func normalizeMaxValueCbk(w tcod.IWidget, val string, data interface{}) {
	f, err := strconv.ParseFloat(val, 32)
	if err != nil {
		op := data.(*NormalizeOperation)
		if float32(f) > op.min {
			op.max = float32(f)
			if peekOperation(operations.list) == IOperation(op) {
				op.runInternal()
			} else {
				operations.reseed()
			}
		}
	}
}

func (self *NormalizeOperation) createParamUi() {
	params.Clear()
	params.SetVisible(true)
	params.SetName(operations.names[NORM])
	tmp := fmt.Sprintf("%f", self.min)

	tbMin := gui.NewTextBoxWithTip(0, 0, 8, 10, "min", tmp, "Heightmap minimum value after the normalization")
	tbMin.SetCallback(normalizeMinValueCbk, self)
	params.AddWidget(tbMin)

	tmp = fmt.Sprintf("%f", self.max)
	tbMax := gui.NewTextBoxWithTip(0, 0, 8, 10, "max", tmp, "Heightmap maximum value after the normalization")
	tbMax.SetCallback(normalizeMaxValueCbk, self)
	params.AddWidget(tbMax)
}

//
// AddFbmOperation
// add noise to the heightmap
//
//
type AddFbmOperation struct {
	Operation
	zoom, offsetx, offsety, octaves, scale, offset float32
}

func NewAddFbmOperation(zoom, offsetx, offsety, octaves, scale, offset float32) *AddFbmOperation {
	result := &AddFbmOperation{}
	result.initializeAddFbmOperation(ADDFBM, zoom, offsetx, offsety, octaves, scale, offset)
	return result
}

func (self *AddFbmOperation) initializeAddFbmOperation(opType OpType, zoom, offsetx, offsety, octaves, scale, offset float32) {
	self.Operation.initializeOperation(opType)
	self.zoom = zoom
	self.offsetx = offsetx
	self.offsety = offsety
	self.octaves = octaves
	self.scale = scale
	self.offset = offset
}

// AddFbm
func (self *AddFbmOperation) getCode(codeType CodeType) string {
	switch codeType {
	case C:
		return fmt.Sprintf(
			"\ttcod.TCOD_heightmap_add_fbm(hm,noise,%g,%g,%g,%g,%g,%g,%g);\n",
			self.zoom, self.zoom, self.offsetx, self.offsety, self.octaves, self.offset, self.scale)
	case CPP:
		return fmt.Sprintf(
			"\thm.addFbm(noise,%g,%g,%g,%g,%g,%g,%g);\n",
			self.zoom, self.zoom, self.offsetx, self.offsety, self.octaves, self.offset, self.scale)
	case PY:
		return fmt.Sprintf(
			"    libtcod.heightmap_add_fbm(hm,noise,%g,%g,%g,%g,%g,%g,%g)\n",
			self.zoom, self.zoom, self.offsetx, self.offsety, self.octaves, self.offset, self.scale)
	case GO:
		return fmt.Sprintf(
			"    hm.AddFbm(noise,%g,%g,%g,%g,%g,%g,%g)\n",
			self.zoom, self.zoom, self.offsetx, self.offsety, self.octaves, self.offset, self.scale)
	default:
	}
	return ""
}

func (self *AddFbmOperation) runInternal() {
	hm.AddFbm(noise, self.zoom, self.zoom, self.offsetx, self.offsety, self.octaves, self.offset, self.scale)
}

func (self *AddFbmOperation) addInternal() bool {
	operations.needsNoise = true
	addFbmDelta += HM_WIDTH
	return true
}

func addFbmZoomValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*AddFbmOperation)
	op.zoom = val
	if peekOperation(operations.list) == IOperation(op) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func addFbmXOffsetValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*AddFbmOperation)
	op.offsetx = val
	if peekOperation(operations.list) == IOperation(op) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func addFbmYOffsetValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*AddFbmOperation)
	op.offsety = val
	if peekOperation(operations.list) == IOperation(op) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func addFbmOctavesValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*AddFbmOperation)
	op.octaves = val
	if peekOperation(operations.list) == IOperation(op) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func addFbmOffsetValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*AddFbmOperation)
	op.offset = val
	if peekOperation(operations.list) == IOperation(op) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func addFbmScaleValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*AddFbmOperation)
	op.scale = val
	if peekOperation(operations.list) == IOperation(op) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func (self *AddFbmOperation) createParamUi() {
	params.Clear()
	params.SetVisible(true)
	params.SetName(operations.names[ADDFBM])

	slider := gui.NewSlider(0, 0, 8, 0.1, 20.0, "zoom       ", "Noise zoom")
	slider.SetCallback(addFbmZoomValueCbk, self)
	params.AddWidget(slider)
	slider.SetValue(self.zoom)

	slider = gui.NewSlider(0, 0, 8, -100.0, 100.0, "xOffset    ", "Horizontal offset in the noise plan")
	slider.SetCallback(addFbmXOffsetValueCbk, self)
	params.AddWidget(slider)
	slider.SetValue(self.offsetx)

	slider = gui.NewSlider(0, 0, 8, -100.0, 100.0, "yOffset    ", "Vertical offset in the noise plan")
	slider.SetCallback(addFbmYOffsetValueCbk, self)
	params.AddWidget(slider)
	slider.SetValue(self.offsety)

	slider = gui.NewSlider(0, 0, 8, 1.0, 10.0, "octaves    ", "Number of octaves for the fractal sum")
	slider.SetCallback(addFbmOctavesValueCbk, self)
	params.AddWidget(slider)
	slider.SetValue(self.octaves)

	slider = gui.NewSlider(0, 0, 8, -1.0, 1.0, "noiseOffset", "Offset added to the noise value")
	slider.SetCallback(addFbmOffsetValueCbk, self)
	params.AddWidget(slider)
	slider.SetValue(self.offset)

	slider = gui.NewSlider(0, 0, 8, 0.01, 10.0, "scale      ", "The noise value is multiplied by self value")
	slider.SetCallback(addFbmScaleValueCbk, self)
	params.AddWidget(slider)
	slider.SetValue(self.scale)

}

//
// scale the heightmap by a noise function
//
//
type ScaleFbmOperation struct {
	AddFbmOperation
}

func NewScaleFbmOperation(zoom, offsetx, offsety, octaves, scale, offset float32) *ScaleFbmOperation {
	result := &ScaleFbmOperation{}
	result.initializeScaleFbmOperation(SCALEFBM, zoom, offsetx, offsety, octaves, scale, offset)
	return result
}

func (self *ScaleFbmOperation) initializeScaleFbmOperation(opType OpType, zoom, offsetx, offsety, octaves, scale, offset float32) {
	self.AddFbmOperation.initializeAddFbmOperation(opType, zoom, offsetx, offsety, octaves, scale, offset)
}

// ScaleFbm
func (self *ScaleFbmOperation) getCode(codeType CodeType) string {
	switch codeType {
	case C:
		return fmt.Sprintf(
			"\ttcod.TCOD_heightmap_scale_fbm(hm,noise,%g,%g,%g,%g,%g,%g,%g);\n"+
				"\tscaleFbmDelta += HM_WIDTH;\n",
			self.zoom, self.zoom, self.offsetx, self.offsety, self.octaves, self.offset, self.scale)
	case CPP:
		return fmt.Sprintf(
			"\thm.scaleFbm(noise,%g,%g,%g,%g,%g,%g,%g);\n"+
				"\tscaleFbmDelta += HM_WIDTH;\n",
			self.zoom, self.zoom, self.offsetx, self.offsety, self.octaves, self.offset, self.scale)
	case PY:
		return fmt.Sprintf(
			"    libtcod.heightmap_scale_fbm(hm,noise,%g,%g,%g,%g,%g,%g,%g)\n"+
				"    scaleFbmDelta += HM_WIDTH\n",
			self.zoom, self.zoom, self.offsetx, self.offsety, self.octaves, self.offset, self.scale)
	case GO:
		return fmt.Sprintf(
			"    hm.ScaleFbm(noise,%g,%g,%g,%g,%g,%g,%g)\n"+
				"    scaleFbmDelta += HM_WIDTH\n",
			self.zoom, self.zoom, self.offsetx, self.offsety, self.octaves, self.offset, self.scale)
	default:
	}
	return ""

}

func (self *ScaleFbmOperation) runInternal() {
	hm.ScaleFbm(noise, self.zoom, self.zoom, self.offsetx, self.offsety, self.octaves, self.offset, self.scale)
}

func (self *ScaleFbmOperation) addInternal() bool {
	operations.needsNoise = true
	scaleFbmDelta += HM_WIDTH
	return true
}

//
//
// Add a hill to the heightmap
//
type AddHillOperation struct {
	Operation
	nbHill                    int
	radius, radiusVar, height float32
}

func NewAddHillOperation(nbHill int, radius, radiusVar, height float32) *AddHillOperation {
	result := &AddHillOperation{}
	result.initializeAddHillOperation(ADDHILL, nbHill, radius, radiusVar, height)
	return result
}

func (self *AddHillOperation) initializeAddHillOperation(opType OpType, nbHill int, radius, radiusVar, height float32) {
	self.Operation.initializeOperation(opType)
	self.nbHill = nbHill
	self.radius = radius
	self.radiusVar = radiusVar
	self.height = height
}

// AddHill
func (self *AddHillOperation) getCode(codeType CodeType) string {
	switch codeType {
	case C:
		return fmt.Sprintf("\taddHill(hm,%d,%g,%g,%g);\n", self.nbHill, self.radius, self.radiusVar, self.height)
	case CPP:
		return fmt.Sprintf("\taddHill(hm,%d,%g,%g,%g);\n", self.nbHill, self.radius, self.radiusVar, self.height)
	case PY:
		return fmt.Sprintf("    addHill(hm,%d,%g,%g,%g)\n", self.nbHill, self.radius, self.radiusVar, self.height)
	case GO:
		return fmt.Sprintf("    addHill(hm,%d,%g,%g,%g)\n", self.nbHill, self.radius, self.radiusVar, self.height)
	default:
	}
	return ""
}

func (self *AddHillOperation) runInternal() {
	addHill(self.nbHill, self.radius, self.radiusVar, self.height)
}

func (self *AddHillOperation) addInternal() bool {
	operations.addInitCode(C,
		"#include <math.h>\n"+
			"void addHill(TCOD_heightmap_t *hm,int nbHill, float32 baseRadius, float32 radiusVar, float32 height)  {\n"+
			"\tint i;\n"+
			"\tfor (i=0; i<  nbHill; i++ ) {\n"+
			"\t\tfloat32 hillMinRadius=baseRadius*(1.0f-radiusVar);\n"+
			"\t\tfloat32 hillMaxRadius=baseRadius*(1.0f+radiusVar);\n"+
			"\t\tfloat32 radius = TCOD_random_get_float32(rnd,hillMinRadius, hillMaxRadius);\n"+
			"\t\tfloat32 theta = TCOD_random_get_float32(rnd,0.0f, 6.283185f); // between 0 and 2Pi\n"+
			"\t\tfloat32 dist = TCOD_random_get_float32(rnd,0.0f, (float32)MIN(HM_WIDTH,HM_HEIGHT)/2 - radius);\n"+
			"\t\tint xh = (int) (HM_WIDTH/2 + cos(theta) * dist);\n"+
			"\t\tint yh = (int) (HM_HEIGHT/2 + sin(theta) * dist);\n"+
			"\t\tTCOD_heightmap_add_hill(hm,(float32)xh,(float32)yh,radius,height);\n"+
			"\t}\n"+
			"}\n")
	operations.addInitCode(CPP,
		"#include <math.h>\n"+
			"void addHill(TCODHeightMap *hm,int nbHill, float32 baseRadius, float32 radiusVar, float32 height)  {\n"+
			"\tfor (int i=0; i<  nbHill; i++ ) {\n"+
			"\t\tfloat32 hillMinRadius=baseRadius*(1.0f-radiusVar);\n"+
			"\t\tfloat32 hillMaxRadius=baseRadius*(1.0f+radiusVar);\n"+
			"\t\tfloat32 radius = rnd->getfloat32(hillMinRadius, hillMaxRadius);\n"+
			"\t\tfloat32 theta = rnd->getfloat32(0.0f, 6.283185f); // between 0 and 2Pi\n"+
			"\t\tfloat32 dist = rnd->getfloat32(0.0f, (float32)MIN(HM_WIDTH,HM_HEIGHT)/2 - radius);\n"+
			"\t\tint xh = (int) (HM_WIDTH/2 + cos(theta) * dist);\n"+
			"\t\tint yh = (int) (HM_HEIGHT/2 + sin(theta) * dist);\n"+
			"\t\thm->addHill((float32)xh,(float32)yh,radius,height);\n"+
			"\t}\n"+
			"}\n")
	operations.addInitCode(PY,
		"def addHill(hm,nbHill,baseRadius,radiusVar,height) :\n"+
			"    for i in range(nbHill) :\n"+
			"        hillMinRadius=baseRadius*(1.0-radiusVar)\n"+
			"        hillMaxRadius=baseRadius*(1.0+radiusVar)\n"+
			"        radius = libtcod.random_get_float32(rnd,hillMinRadius, hillMaxRadius)\n"+
			"        theta = libtcod.random_get_float32(rnd,0.0, 6.283185) # between 0 and 2Pi\n"+
			"        dist = libtcod.random_get_float32(rnd,0.0, float32(min(HM_WIDTH,HM_HEIGHT))/2 - radius)\n"+
			"        xh = int(HM_WIDTH/2 + math.cos(theta) * dist)\n"+
			"        yh = int(HM_HEIGHT/2 + math.sin(theta) * dist)\n"+
			"        libtcod.heightmap_add_hill(hm,float32(xh),float32(yh),radius,height)\n")
	operations.addInitCode(GO,

		"func addHill(hm *tcod.HeightMap, nbHill int, baseRadius float32, radiusVar float32, height float32)  {\n"+
			"\tfor i:=0; i<  nbHill; i++ {\n"+
			"\t\thillMinRadius:=baseRadius*(1.0-radiusVar)\n"+
			"\t\thillMaxRadius:=baseRadius*(1.0+radiusVar)\n"+
			"\t\tradius := rnd.Getfloat32(hillMinRadius, hillMaxRadius)\n"+
			"\t\ttheta := rnd.Getfloat32(0.0, 6.283185) // between 0 and 2Pi\n"+
			"\t\tdist := rnd.Getfloat32(0.0, float32(min(HM_WIDTH,HM_HEIGHT)/2 - radius))\n"+
			"\t\txh := int(HM_WIDTH/2 + cos(theta) * dist)\n"+
			"\t\tyh := int(HM_HEIGHT/2 + sin(theta) * dist)\n"+
			"\t\thm.AddHill(float32(xh),float32(yh),radius,height)\n"+
			"\t}\n"+
			"}\n")

	return true
}

func addHillNbHillValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*AddHillOperation)
	op.nbHill = int(val)
	if peekOperation(operations.list) == IOperation(op) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func addHillRadiusValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*AddHillOperation)
	op.radius = val
	if peekOperation(operations.list) == IOperation(op) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func addHillRadiusVarValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*AddHillOperation)
	op.radiusVar = val
	if peekOperation(operations.list) == IOperation(op) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func addHillHeightValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*AddHillOperation)
	op.height = val
	if peekOperation(operations.list) == IOperation(op) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func (self *AddHillOperation) createParamUi() {

	params.Clear()
	params.SetVisible(true)
	params.SetName(operations.names[ADDHILL])

	slider := gui.NewSlider(0, 0, 8, 1.0, 50.0, "nbHill   ", "Number of hills")
	slider.SetCallback(addHillNbHillValueCbk, self)
	slider.SetFormat("%.0f")
	slider.SetSensitivity(2.0)
	params.AddWidget(slider)
	slider.SetValue(float32(self.nbHill))

	slider = gui.NewSlider(0, 0, 8, 1.0, 30.0, "radius   ", "Average radius of the hills")
	slider.SetCallback(addHillRadiusValueCbk, self)
	slider.SetFormat("%.1f")
	params.AddWidget(slider)
	slider.SetValue(self.radius)

	slider = gui.NewSlider(0, 0, 8, 0.0, 1.0, "radiusVar", "Variation of the radius of the hills")
	slider.SetCallback(addHillRadiusVarValueCbk, self)
	params.AddWidget(slider)
	slider.SetValue(self.radiusVar)

	slider = gui.NewSlider(0, 0, 8, 0.0, tcod.If(mapmax == mapmin, float32(1.0), float32((mapmax-mapmin)*0.5)).(float32), "height   ", "Height of the hills")
	slider.SetCallback(addHillHeightValueCbk, self)
	params.AddWidget(slider)
	slider.SetValue(self.height)
}

// add a scalar to the heightmap
//
type AddLevelOperation struct {
	Operation
	level float32
}

func NewAddLevelOperation(level float32) *AddLevelOperation {
	result := &AddLevelOperation{}
	result.initializeAddLevelOperation(ADDLEVEL, level)
	return result
}

func (self *AddLevelOperation) initializeAddLevelOperation(opType OpType, level float32) {
	self.Operation.initializeOperation(opType)
	self.level = level
}

// AddLevel
func (self *AddLevelOperation) getCode(codeType CodeType) string {
	switch codeType {
	case C:
		return fmt.Sprintf("\ttcod.TCOD_heightmap_add(hm,%g)\n\ttcod.TCOD_heightmap_Clamp(hm,0.0f,1.0f)\n", self.level)
	case CPP:
		return fmt.Sprintf("\thm.add(%g)\n\thm.Clamp(0.0f,1.0f)\n", self.level)
	case PY:
		return fmt.Sprintf("    libtcod.heightmap_add(hm,%g)\n    libtcod.heightmap_Clamp(hm,0.0,1.0)\n", self.level)
	case GO:
		return fmt.Sprintf("    hm.Add(%g)\n    hm.Clamp(0.0,1.0)\n", self.level)
	default:
	}
	return ""
}

func (self *AddLevelOperation) runInternal() {
	var min, max float32
	min, max = hm.GetMinMax()
	hm.Add(self.level)
	if min != max {
		hm.Clamp(min, max)
	}
}

func (self *AddLevelOperation) addInternal() bool {
	prev := peekOperation(operations.list)
	ret := true
	if prev != nil && prev.getOpType() == ADDLEVEL {
		// cumulated consecutive addLevel operation into a single call
		addOp := prev.(*AddLevelOperation)
		if addOp.level*self.level > 0 {
			addOp.level += self.level
			ret = false
		}
	}
	return ret
}

func raiseLowerValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*AddLevelOperation)
	op.level = val
	if IOperation(op) == peekOperation(operations.list) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func (self *AddLevelOperation) createParamUi() {
	params.Clear()
	params.SetName(operations.names[ADDLEVEL])
	params.SetVisible(true)

	slider := gui.NewSlider(0, 0, 8, -1.0, 1.0, "zOffset", "z value to add to the whole map")
	slider.SetCallback(raiseLowerValueCbk, self)
	params.AddWidget(slider)
	minLevel, maxLevel := hm.GetMinMax()
	if maxLevel == minLevel {
		slider.SetMinMax(-1.0, 1.0)
	} else {
		slider.SetMinMax(-(maxLevel - minLevel), (maxLevel - minLevel))
	}
	slider.SetValue(self.level)
}

// smooth a part of the heightmap
//
type SmoothOperation struct {
	Operation
	minLevel, maxLevel, radius float32
	count                      int
}

func NewSmoothOperation(minLevel, maxLevel float32, count int) *SmoothOperation {
	result := &SmoothOperation{}
	result.initializeSmoothOperation(SMOOTH, minLevel, maxLevel, count)
	return result
}

func (self *SmoothOperation) initializeSmoothOperation(opType OpType, minLevel, maxLevel float32, count int) {
	self.Operation.initializeOperation(opType)
	self.minLevel = minLevel
	self.maxLevel = maxLevel
	self.count = count
	self.radius = 0
}

// Smooth
func (self *SmoothOperation) getCode(codeType CodeType) string {
	switch codeType {
	case C:
		return fmt.Sprintf(
			"\tsmoothKernelWeight[4] = %g;\n"+
				"\t{\n"+
				"\t\tint i;\n"+
				"\t\tfor (i=%d; i>= 0; i--) {\n"+
				"\t\t\ttcod.TCOD_heightmap_kernel_transform(hm,smoothKernelSize,smoothKernelDx,smoothKernelDy,smoothKernelWeight,%g,%g);\n"+
				"\t\t}\n"+
				"\t}\n",
			20-self.radius*19, self.count, self.minLevel, self.maxLevel)
	case CPP:
		return fmt.Sprintf(
			"\tsmoothKernelWeight[4] = %g;\n"+
				"\tfor (int i=%d; i>= 0; i--) {\n"+
				"\t\thm.kernelTransform(smoothKernelSize,smoothKernelDx,smoothKernelDy,smoothKernelWeight,%g,%g);\n"+
				"\t}\n",
			20-self.radius*19, self.count, self.minLevel, self.maxLevel)
	case PY:
		return fmt.Sprintf(
			"    smoothKernelWeight[4] = %g\n"+
				"    for i in range(%d,-1,-1) :\n"+
				"        libtcod.heightmap_kernel_transform(hm,smoothKernelSize,smoothKernelDx,smoothKernelDy,smoothKernelWeight,%g,%g)\n",
			20-self.radius*19, self.count, self.minLevel, self.maxLevel)
	case GO:
		return fmt.Sprintf(
			"    smoothKernelWeight[4] = %g\n"+
				"    for i := %d; i >= 0; i-- {\n"+
				"        hm.KernelTransform(smoothKernelSize,smoothKernelDx,smoothKernelDy,smoothKernelWeight,%g,%g)\n"+
				"    }\n ",
			20-self.radius*19, self.count, self.minLevel, self.maxLevel)
	default:
	}
	return ""

}

func (self *SmoothOperation) runInternal() {
	smoothKernelWeight[4] = 20 - self.radius*19
	for i := self.count; i >= 0; i-- {
		hm.KernelTransform(smoothKernelSize, smoothKernelDx, smoothKernelDy, smoothKernelWeight, self.minLevel, self.maxLevel)
	}
}

func (self *SmoothOperation) addInternal() bool {
	operations.addInitCode(C,
		"// 3x3 kernel for smoothing operations\n"+
			"int smoothKernelSize=9;\n"+
			"int smoothKernelDx[9]={-1,0,1,-1,0,1,-1,0,1};\n"+
			"int smoothKernelDy[9]={-1,-1,-1,0,0,0,1,1,1};\n"+
			"float32 smoothKernelWeight[9]={1,2,1,2,20,2,1,2,1};\n")
	operations.addInitCode(CPP,
		"// 3x3 kernel for smoothing operations\n"+
			"int smoothKernelSize=9;\n"+
			"int smoothKernelDx[9]={-1,0,1,-1,0,1,-1,0,1};\n"+
			"int smoothKernelDy[9]={-1,-1,-1,0,0,0,1,1,1};\n"+
			"float32 smoothKernelWeight[9]={1,2,1,2,20,2,1,2,1};\n")
	operations.addInitCode(PY,
		"# 3x3 kernel for smoothing operations\n"+
			"smoothKernelSize=9\n"+
			"smoothKernelDx=[-1,0,1,-1,0,1,-1,0,1]\n"+
			"smoothKernelDy=[-1,-1,-1,0,0,0,1,1,1]\n"+
			"smoothKernelWeight=[1.0,2.0,1.0,2.0,20.0,2.0,1.0,2.0,1.0]\n")
	operations.addInitCode(GO,
		"// 3x3 kernel for smoothing operations\n"+
			"var smoothKernelSize int = 9\n"+
			"var smoothKernelDx [9]int =[9]int {-1,0,1,-1,0,1,-1,0,1}\n"+
			"var smoothKernelDy [9]int = [9]int {-1,-1,-1,0,0,0,1,1,1}\n"+
			"var smoothKernelWeight [9]float32 = [9]float32 {1,2,1,2,20,2,1,2,1}\n")

	return true
}

func smoothMinValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*SmoothOperation)
	op.minLevel = val
	if IOperation(op) == peekOperation(operations.list) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func smoothMaxValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*SmoothOperation)
	op.maxLevel = val
	if IOperation(op) == peekOperation(operations.list) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func smoothRadiusValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*SmoothOperation)
	op.radius = val
	if IOperation(op) == peekOperation(operations.list) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func smoothCountValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*SmoothOperation)
	op.count = int(val)
	if IOperation(op) == peekOperation(operations.list) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func (self *SmoothOperation) createParamUi() {
	params.Clear()
	params.SetName(operations.names[SMOOTH])
	params.SetVisible(true)

	slider := gui.NewSlider(0, 0, 8, minf(0.0, self.minLevel), maxf(1.0, self.maxLevel), "minLevel", "Land level above which the smooth operation is applied")
	slider.SetCallback(smoothMinValueCbk, self)
	params.AddWidget(slider)
	slider.SetValue(self.minLevel)

	slider = gui.NewSlider(0, 0, 8, minf(0.0, self.minLevel), maxf(1.0, self.maxLevel), "maxLevel", "Land level below which the smooth operation is applied")
	slider.SetCallback(smoothMaxValueCbk, self)
	params.AddWidget(slider)
	slider.SetValue(self.maxLevel)

	slider = gui.NewSlider(0, 0, 8, 1.0, 20.0, "amount", "Number of times the smoothing operation is applied")
	slider.SetCallback(smoothCountValueCbk, self)
	slider.SetFormat("%.0f")
	slider.SetSensitivity(4.0)
	params.AddWidget(slider)
	slider.SetValue(float32(self.count))

	slider = gui.NewSlider(0, 0, 8, 0.0, 1.0, "sharpness", "Radius of the blurring effect")
	slider.SetCallback(smoothRadiusValueCbk, self)
	params.AddWidget(slider)
	slider.SetValue(0.0)
}

//
//
// simulate rain erosion
//
//
type RainErosionOperation struct {
	Operation
	nbDroperations                 int
	erosionCoef, sedimentationCoef float32
}

func NewRainErosionOperation(nbDroperations int, erosionCoef, sedimentationCoef float32) *RainErosionOperation {
	result := &RainErosionOperation{}
	result.initializeRainErosionOperation(RAIN, nbDroperations, erosionCoef, sedimentationCoef)
	return result
}

func (self *RainErosionOperation) initializeRainErosionOperation(opType OpType, nbDroperations int, erosionCoef, sedimentationCoef float32) {
	self.Operation.initializeOperation(opType)
	self.nbDroperations = nbDroperations
	self.erosionCoef = erosionCoef
	self.sedimentationCoef = sedimentationCoef
}

// Rain
func (self *RainErosionOperation) getCode(codeType CodeType) string {
	switch codeType {
	case C:
		return fmt.Sprintf("\ttcod.TCOD_heightmap_rain_erosion(hm,%d,%g,%g,rnd);\n", self.nbDroperations, self.erosionCoef, self.sedimentationCoef)
	case CPP:
		return fmt.Sprintf("\thm.rainErosion(%d,%g,%g,rnd);\n", self.nbDroperations, self.erosionCoef, self.sedimentationCoef)
	case PY:
		return fmt.Sprintf("    libtcod.heightmap_rain_erosion(hm,%d,%g,%g,rnd)\n", self.nbDroperations, self.erosionCoef, self.sedimentationCoef)
	case GO:
		return fmt.Sprintf("    hm.RainErosion(%d,%g,%g,rnd)\n", self.nbDroperations, self.erosionCoef, self.sedimentationCoef)
	default:
	}
	return ""
}

func (self *RainErosionOperation) runInternal() {
	if !isNormalized {
		hm.Normalize()
	}
	hm.RainErosion(self.nbDroperations, self.erosionCoef, self.sedimentationCoef, rnd)
}

func (self *RainErosionOperation) addInternal() bool {
	operations.needsRandom = true
	return true
}

func rainErosionNbDroperationsValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*RainErosionOperation)
	op.nbDroperations = int(val)
	if IOperation(op) == peekOperation(operations.list) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func rainErosionErosionCoefValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*RainErosionOperation)
	op.erosionCoef = val
	if IOperation(op) == peekOperation(operations.list) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func rainErosionSedimentationCoefValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*RainErosionOperation)
	op.sedimentationCoef = val
	if IOperation(op) == peekOperation(operations.list) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func (self *RainErosionOperation) createParamUi() {
	params.Clear()
	params.SetName(operations.names[RAIN])
	params.SetVisible(true)

	slider := gui.NewSlider(0, 0, 8, 1000.0, 20000.0, "nbDroperations      ", "Number of rain droperations simulated")
	slider.SetCallback(rainErosionNbDroperationsValueCbk, self)
	params.AddWidget(slider)
	slider.SetFormat("%.0f")
	slider.SetValue(float32(self.nbDroperations))

	slider = gui.NewSlider(0, 0, 8, 0.01, 1.0, "erosion      ", "Erosion effect amount")
	slider.SetCallback(rainErosionErosionCoefValueCbk, self)
	params.AddWidget(slider)
	slider.SetValue(self.erosionCoef)

	slider = gui.NewSlider(0, 0, 8, 0.01, 1.0, "sedimentation", "Sedimentation effect amount")
	slider.SetCallback(rainErosionSedimentationCoefValueCbk, self)
	params.AddWidget(slider)
	slider.SetValue(self.sedimentationCoef)
}

//
//
// lerp current heightmap with a noise
//
type NoiseLerpOperation struct {
	AddFbmOperation
	coef float32
}

func NewNoiseLerpOperation(coef, zoom, offsetx, offsety, octaves, scale, offset float32) *NoiseLerpOperation {
	result := &NoiseLerpOperation{}
	result.initializeNoiseLerpOperation(NOISELERP, coef, zoom, offsetx, offsety, octaves, scale, offset)
	return result
}

func (self *NoiseLerpOperation) initializeNoiseLerpOperation(opType OpType, coef, zoom, offsetx, offsety, octaves, scale, offset float32) {
	self.AddFbmOperation.initializeAddFbmOperation(opType, zoom, offsetx, offsety, octaves, scale, offset)
	self.coef = coef
}

// NoiseLerp
func (self *NoiseLerpOperation) getCode(codeType CodeType) string {
	switch codeType {
	case C:
		return fmt.Sprintf(
			"\t{\n"+
				"\t\ttcod.TCOD_heightmap_t *tmp=tcod.TCOD_heightmap_new(HM_WIDTH,HM_HEIGHT);\n"+
				"\t\ttcod.TCOD_heightmap_add_fbm(tmp,noise,%g,%g,%g,%g,%g,%g,%g);\n"+
				"\t\ttcod.TCOD_heightmap_lerp(hm,tmp,hm,%g);\n"+
				"\t\ttcod.TCOD_heightmap_delete(tmp);\n"+
				"\t}\n",
			self.zoom, self.zoom, self.offsetx, self.offsety, self.octaves, self.offset, self.scale, self.coef)
	case CPP:
		return fmt.Sprintf(
			"\t{\n"+
				"\t\tTCODHeightMap tmp(HM_WIDTH,HM_HEIGHT);\n"+
				"\t\ttmp.addFbm(noise,%g,%g,%g,%g,%g,%g,%g);\n"+
				"\t\thm.lerp(hm,&tmp,%g);\n"+
				"\t}\n",
			self.zoom, self.zoom, self.offsetx, self.offsety, self.octaves, self.offset, self.scale, self.coef)
	case PY:
		return fmt.Sprintf(
			"    tmp=libtcod.heightmap_new(HM_WIDTH,HM_HEIGHT)\n"+
				"    libtcod.heightmap_add_fbm(tmp,noise,%g,%g,%g,%g,%g,%g,%g)\n"+
				"    libtcod.heightmap_lerp(hm,tmp,hm,%g)\n"+
				"    libtcod.heightmap_delete(tmp)\n",
			self.zoom, self.zoom, self.offsetx, self.offsety, self.octaves, self.offset, self.scale, self.coef)
	case GO:
		return fmt.Sprintf(
			"    tmp:=tcod.NewHeightMap(HM_WIDTH,HM_HEIGHT)\n"+
				"    tmp.AddFbm(noise,%g,%g,%g,%g,%g,%g,%g)\n"+
				"    hm.Lerp(hm,tmp,%g)\n",
			self.zoom, self.zoom, self.offsetx, self.offsety, self.octaves, self.offset, self.scale, self.coef)
	default:
	}
	return ""

}

func (self *NoiseLerpOperation) runInternal() {
	tmp := tcod.NewHeightMap(HM_WIDTH, HM_HEIGHT)
	tmp.AddFbm(noise, self.zoom, self.zoom, self.offsetx, self.offsety, self.octaves, self.offset, self.scale)
	hm.Lerp(hm, tmp, self.coef)
}

func (self *NoiseLerpOperation) addInternal() bool {
	operations.needsNoise = true
	addFbmDelta += HM_WIDTH
	return true
}

func noiseLerpValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*NoiseLerpOperation)
	op.coef = val
	if peekOperation(operations.list) == IOperation(op) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func (self *NoiseLerpOperation) createParamUi() {
	self.AddFbmOperation.createParamUi()
	params.SetName(operations.names[NOISELERP])

	slider := gui.NewSlider(0, 0, 8, -1.0, 1.0, "coef       ", "Coefficient of the lerp operation")
	slider.SetCallback(noiseLerpValueCbk, self)
	params.AddWidget(slider)
	slider.SetValue(self.coef)
}

// add a voronoi diagram
//
const MAX_VORONOI_COEF = 5

type VoronoiOperation struct {
	Operation
	nbPoints   int
	nbCoef     int
	coef       []float32
	coefSlider [MAX_VORONOI_COEF]*tcod.Slider
}

func NewVoronoiOperation(nbPoints, nbCoef int, coef []float32) *VoronoiOperation {
	result := &VoronoiOperation{}
	result.coef = make([]float32, MAX_VORONOI_COEF)
	result.initializeVoronoiOperation(VORONOI, nbPoints, nbCoef, coef)
	return result
}

func (self *VoronoiOperation) initializeVoronoiOperation(opType OpType, nbPoints, nbCoef int, coef []float32) {
	self.Operation.initializeOperation(opType)
	self.nbPoints = nbPoints
	self.nbCoef = nbCoef
	for i := 0; i < min(MAX_VORONOI_COEF, nbCoef); i++ {
		self.coef[i] = coef[i]
	}
	for i := min(MAX_VORONOI_COEF, nbCoef); i < MAX_VORONOI_COEF; i++ {
		self.coef[i] = 0.0
	}
	for i := 0; i < MAX_VORONOI_COEF; i++ {
		self.coefSlider[i] = nil
	}
}

func (self *VoronoiOperation) getCode(codeType CodeType) string {
	var coefstr string
	for i := 0; i < self.nbCoef; i++ {
		coefstr += fmt.Sprintf("%g,", self.coef[i])
	}
	switch codeType {
	case C:
		return fmt.Sprintf(
			"\t{\n"+
				"\t\tfloat32 coef[]={%s};\n"+
				"\t\ttcod.TCOD_heightmap_t *tmp =tcod.TCOD_heightmap_new(HM_WIDTH,HM_HEIGHT);\n"+
				"\t\ttcod.TCOD_heightmap_add_voronoi(tmp,%d,%d,coef,rnd);\n"+
				"\t\ttcod.TCOD_heightmap_normalize(tmp,0.0f,1.0f);\n"+
				"\t\ttcod.TCOD_heightmap_add(hm,tmp,hm);\n"+
				"\t\ttcod.TCOD_heightmap_delete(tmp);\n"+
				"\t}\n",
			coefstr, self.nbPoints, self.nbCoef)
	case CPP:
		return fmt.Sprintf(
			"\t{\n"+
				"\t\tfloat32 coef[]={%s};\n"+
				"\t\tTCODHeightMap tmp(HM_WIDTH,HM_HEIGHT);\n"+
				"\t\ttmp.addVoronoi(%d,%d,coef,rnd);\n"+
				"\t\ttmp.normalize();\n"+
				"\t\thm.add(hm,&tmp);\n"+
				"\t}\n",
			coefstr, self.nbPoints, self.nbCoef)
	case PY:
		return fmt.Sprintf(
			"    coef=[%s]\n"+
				"    tmp =libtcod.heightmap_new(HM_WIDTH,HM_HEIGHT)\n"+
				"    libtcod.heightmap_add_voronoi(tmp,%d,%d,coef,rnd)\n"+
				"    libtcod.heightmap_normalize(tmp)\n"+
				"    libtcod.heightmap_add(hm,tmp,hm)\n"+
				"    libtcod.heightmap_delete(tmp)\n",
			coefstr, self.nbPoints, self.nbCoef)
	case GO:
		return fmt.Sprintf(
			"    var coef []float32 = []float32 {%s}\n"+
				"    tmp := tcod.NewHeightMap(HM_WIDTH,HM_HEIGHT)\n"+
				"    tmp.AddVoronoi(%d,%d,coef,rnd)\n"+
				"    tmp.Normalize()\n"+
				"    hm.AddHm(hm,tmp)\n",
			coefstr, self.nbPoints, self.nbCoef)
	default:
	}
	return ""
}

func (self *VoronoiOperation) runInternal() {
	tmp := tcod.NewHeightMap(HM_WIDTH, HM_HEIGHT)
	tmp.AddVoronoi(self.nbPoints, self.nbCoef, self.coef, rnd)
	tmp.Normalize()
	hm.AddHm(hm, tmp)
}

func (self *VoronoiOperation) addInternal() bool {
	operations.needsRandom = true
	return true
}

func voronoiNbPointsValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*VoronoiOperation)
	op.nbPoints = int(val)
	if peekOperation(operations.list) == IOperation(op) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func voronoiNbCoefValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*VoronoiOperation)
	op.nbCoef = int(val)
	for i := 0; i < MAX_VORONOI_COEF; i++ {
		if i < op.nbCoef {
			op.coefSlider[i].SetVisible(true)
		} else {
			op.coefSlider[i].SetVisible(false)
		}
	}
	params.ComputeSize()
	if peekOperation(operations.list) == IOperation(op) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func voronoiCoefValueCbk(w tcod.IWidget, val float32, data interface{}) {
	op := data.(*VoronoiOperation)
	var coefnum int
	for coefnum = 0; coefnum < op.nbCoef; coefnum++ {
		if tcod.IWidget(op.coefSlider[coefnum]) == w {
			break
		}
	}
	op.coef[coefnum] = val
	if peekOperation(operations.list) == IOperation(op) {
		restore()
		op.runInternal()
	} else {
		operations.reseed()
	}
}

func (self *VoronoiOperation) createParamUi() {
	params.Clear()
	params.SetName(operations.names[VORONOI])
	params.SetVisible(true)

	slider := gui.NewSlider(0, 0, 8, 1.0, 50.0, "nbPoints", "Number of Voronoi points")
	slider.SetCallback(voronoiNbPointsValueCbk, self)
	params.AddWidget(slider)
	slider.SetFormat("%.0f")
	slider.SetSensitivity(2.0)
	slider.SetValue(float32(self.nbPoints))

	slider = gui.NewSlider(0, 0, 8, 1.0, float32(MAX_VORONOI_COEF-1), "nbCoef  ", "Number of Voronoi coefficients")
	slider.SetCallback(voronoiNbCoefValueCbk, self)
	params.AddWidget(slider)
	slider.SetSensitivity(4.0)
	slider.SetFormat("%.0f")
	slider.SetValue(float32(self.nbCoef))

	for i := 0; i < MAX_VORONOI_COEF; i++ {
		tmp := fmt.Sprintf("coef[%d] ", i)
		self.coefSlider[i] = gui.NewSlider(0, 0, 8, -5.0, 5.0, tmp, "Coefficient of Voronoi points")
		self.coefSlider[i].SetCallback(voronoiCoefValueCbk, self)
		params.AddWidget(self.coefSlider[i])
		self.coefSlider[i].SetValue(float32(self.coef[i]))
		if i >= self.nbCoef {
			self.coefSlider[i].SetVisible(false)
		}
	}
}

//
//
// Main program
//

func initColors() {
	tcod.ColorGenMap(mapGradient, nbColorKeys, keyColor, keyIndex)
}

func render() {
	isNormalized = true
	root.SetDefaultBackground(tcod.COLOR_BLACK)
	root.SetDefaultForeground(tcod.COLOR_WHITE)
	root.Clear()
	backupMap.Copy(hm)
	mapmin = 1e8
	mapmax = -1e8
	for i := 0; i < HM_WIDTH*HM_HEIGHT; i++ {
		v := hm.GetNthValue(i)
		if v < 0.0 || v > 1.0 {
			isNormalized = false
		}
		if v < mapmin {
			mapmin = v
		}
		if v > mapmax {
			mapmax = v
		}
	}

	if !isNormalized {
		backupMap.Normalize()
	}
	// render the TCODHeightMap
	for x := 0; x < HM_WIDTH; x++ {
		for y := 0; y < HM_HEIGHT; y++ {
			//z := backupMap.GetValue(x, y)
			z := hm.GetValue(x, y)
			val := uint8(z * 255)
			if slope {
				// render the slope map
				z = tcod.ClampF(0.0, 1.0, hm.GetSlope(x, y)*10.0)
				val = uint8(z * 255)
				root.SetCharBackground(x, y, tcod.Color{val, val, val}, tcod.BKGND_SET)
			} else if greyscale {
				// render the greyscale heightmap
				root.SetCharBackground(x, y, tcod.Color{val, val, val}, tcod.BKGND_SET)
			} else if normal {
				// render the normal map
				var n [3]float32
				hm.GetNormal(float32(x), float32(y), &n, mapmin)
				r := byte((n[0]*0.5 + 0.5) * 255)
				g := byte((n[1]*0.5 + 0.5) * 255)
				b := byte((n[2]*0.5 + 0.5) * 255)
				root.SetCharBackground(x, y, tcod.Color{r, g, b}, tcod.BKGND_SET)
			} else {
				// render the colored heightmap
				root.SetCharBackground(x, y, mapGradient[val], tcod.BKGND_SET)
			}
		}
	}

	minZLabel.SetValue(fmt.Sprintf("min z    : %.2f", mapmin))
	maxZLabel.SetValue(fmt.Sprintf("max z    : %.2f", mapmax))
	seedLabel.SetValue(fmt.Sprintf("seed     : %X", seed))
	landProportion := 100.0 - 100.0*backupMap.CountCells(0.0, sandHeight)/(hm.GetWidth()*hm.GetHeight())
	landMassLabel.SetValue(fmt.Sprintf("landMass : %d %%%%", int(landProportion)))
	if !isNormalized {
		root.PrintEx(HM_WIDTH/2, HM_HEIGHT-1, tcod.BKGND_NONE, tcod.CENTER, "the map is not normalized !")
	}
	// message
	msgDelay -= tcod.SysGetLastFrameLength()
	if msg != "" && msgDelay > 0.0 {
		h := root.PrintRectEx(HM_WIDTH/2, HM_HEIGHT/2+1, HM_WIDTH/2-2, 0, tcod.BKGND_NONE, tcod.CENTER, msg)
		root.SetDefaultBackground(tcod.COLOR_LIGHT_BLUE)
		if h > 0 {
			root.Rect(HM_WIDTH/4, HM_HEIGHT/2, HM_WIDTH/2, h+2, false, tcod.BKGND_SET)
		}
		root.SetDefaultBackground(tcod.COLOR_BLACK)
	}
}

func message(delay float32, fmts string, v ...interface{}) {
	msg = fmt.Sprintf(fmts, v...)
	msgDelay = delay
}

func backup() {
	// save the heightmap & RNG states
	for x := 0; x < HM_WIDTH; x++ {
		for y := 0; y < HM_HEIGHT; y++ {
			hmold.SetValue(x, y, hm.GetValue(x, y))
		}
	}
	fmt.Printf("Saving to backupRNG!\n")
	backupRnd = rnd.Save()
	oldNormalized = isNormalized
	oldmapmax = mapmax
	oldmapmin = mapmin
}

func restore() {
	// restore the previously saved heightmap & RNG states
	for x := 0; x < HM_WIDTH; x++ {
		for y := 0; y < HM_HEIGHT; y++ {
			hm.SetValue(x, y, hmold.GetValue(x, y))
		}
	}
	fmt.Printf("Restoring from backupRnd \n")
	if backupRnd != nil {
		rnd.Restore(backupRnd)
	}
	isNormalized = oldNormalized
	mapmax = oldmapmax
	mapmin = oldmapmin
}

func save() {
	// TODO
	message(2.0, "Saved.")
}

func load() {
	// TODO
}

func addHill(nbHill int, baseRadius, radiusVar, height float32) {
	for i := 0; i < nbHill; i++ {
		hillMinRadius := baseRadius * (1.0 - radiusVar)
		hillMaxRadius := baseRadius * (1.0 + radiusVar)
		radius := rnd.GetFloat(hillMinRadius, hillMaxRadius)
		theta := rnd.GetFloat(0.0, 6.283185) // between 0 and 2Pi
		dist := rnd.GetFloat(0.0, float32(min(HM_WIDTH, HM_HEIGHT)/2)-radius)
		xh := int(HM_WIDTH/2 + cos(theta)*dist)
		yh := int(HM_HEIGHT/2 + sin(theta)*dist)
		hm.AddHill(float32(xh), float32(yh), radius, height)
	}
}

func clearCbk(w tcod.IWidget, data interface{}) {
	hm.Clear()
	operations.clear()
	history.Clear()
	params.Clear()
	params.SetVisible(false)
}

func reseedCbk(w tcod.IWidget, data interface{}) {
	// //seed = int32(rnd.GetInt(0x7FFFFFFF, 0xFFFFFFFF))
	// TODO
	seed = uint32(rnd.GetInt(1, 9999999))
	operations.reseed()
	message(3.0, "Switching to seed %X", seed)
}

func cancelCbk(w tcod.IWidget, data interface{}) {
	operations.cancel()
}

// operations buttons callbacks
func normalizeCbk(w tcod.IWidget, data interface{}) {
	operations.add(NewNormalizeOperation(0.0, 1.0))
}

func addFbmCbk(w tcod.IWidget, data interface{}) {
	operations.add(NewAddFbmOperation(1.0, addFbmDelta, 0.0, 6.0, 1.0, 0.5))
}

func scaleFbmCbk(w tcod.IWidget, data interface{}) {
	operations.add(NewScaleFbmOperation(1.0, addFbmDelta, 0.0, 6.0, 1.0, 0.5))
}

func addHillCbk(w tcod.IWidget, data interface{}) {
	operations.add(NewAddHillOperation(25, 10.0, 0.5, tcod.If(mapmax == mapmin, float32(0.5), float32((mapmax-mapmin)*0.1)).(float32)))
}

func rainErosionCbk(w tcod.IWidget, data interface{}) {
	operations.add(NewRainErosionOperation(1000, 0.05, 0.05))
}

func smoothCbk(w tcod.IWidget, data interface{}) {
	operations.add(NewSmoothOperation(mapmin, mapmax, 2))
}

func voronoiCbk(w tcod.IWidget, data interface{}) {
	operations.add(NewVoronoiOperation(100, 2, voronoiCoef))
}

func noiseLerpCbk(w tcod.IWidget, data interface{}) {
	v := (mapmax - mapmin) * 0.5
	if v == 0.0 {
		v = 1.0
	}
	operations.add(NewNoiseLerpOperation(0.0, 1.0, addFbmDelta, 0.0, 6.0, v, v))
}

func raiseLowerCbk(w tcod.IWidget, data interface{}) {
	operations.add(NewAddLevelOperation(0.0))
}

// In/Out buttons callbacks

func exportCCbk(w tcod.IWidget, data interface{}) {
	code := operations.buildCode(C)
	ioutil.WriteFile("hm.c", ([]byte)(code), 0664)
	message(3.0, "The code has been exported to ./hm.c")
}

func exportPyCbk(w tcod.IWidget, data interface{}) {
	code := operations.buildCode(PY)
	ioutil.WriteFile("hm.py", ([]byte)(code), 0664)
	message(3.0, "The code has been exported to ./hm.py")
}

func exportGoCbk(w tcod.IWidget, data interface{}) {
	code := operations.buildCode(GO)
	ioutil.WriteFile("hm.go", ([]byte)(code), 0664)
	message(3.0, "The code has been exported to ./hm.go")
}

func exportCppCbk(w tcod.IWidget, data interface{}) {
	code := operations.buildCode(CPP)
	ioutil.WriteFile("hm.cpp", ([]byte)(code), 0664)
	message(3.0, "The code has been exported to ./hm.cpp")
}

func exportBmpCbk(w tcod.IWidget, data interface{}) {
	img := tcod.NewImage(HM_WIDTH, HM_HEIGHT)
	for x := 0; x < HM_WIDTH; x++ {
		for y := 0; y < HM_HEIGHT; y++ {
			z := hm.GetValue(x, y)
			val := (uint8)(z * 255)
			if slope {
				// render the slope map
				z = tcod.ClampF(0.0, 1.0, hm.GetSlope(x, y)*10.0)
				val = (uint8)(z * 255)
				img.PutPixel(x, y, tcod.Color{val, val, val})
			} else if greyscale {
				// render the greyscale heightmap
				img.PutPixel(x, y, tcod.Color{val, val, val})
			} else if normal {
				// render the normal map
				var n [3]float32
				hm.GetNormal(float32(x), float32(y), &n, mapmin)
				r := uint8((n[0]*0.5 + 0.5) * 255)
				g := uint8((n[1]*0.5 + 0.5) * 255)
				b := uint8((n[2]*0.5 + 0.5) * 255)
				img.PutPixel(x, y, tcod.Color{r, g, b})
			} else {
				// render the colored heightmap
				img.PutPixel(x, y, mapGradient[val])
			}
		}
	}
	img.Save("hm.bmp")
	message(3.0, "The bitmap has been exported to ./hm.bmp")
}

// Display buttons callbacks
func colorMapCbk(w tcod.IWidget, data interface{}) {
	slope = false
	greyscale = false
	normal = false
}

func greyscaleCbk(w tcod.IWidget, data interface{}) {
	slope = false
	greyscale = true
	normal = false
}

func slopeCbk(w tcod.IWidget, data interface{}) {
	slope = true
	greyscale = false
	normal = false
}

func normalCbk(w tcod.IWidget, data interface{}) {
	slope = false
	greyscale = false
	normal = true
}

func changeColorMapIdxCbk(w tcod.IWidget, val float32, data interface{}) {
	i := data.(int)
	keyIndex[i] = int(val)
	if i == 1 {
		sandHeight = float32(i) / 255.0
	}
	initColors()
}

func changeColorMapRedCbk(w tcod.IWidget, val float32, data interface{}) {
	i := data.(int)
	keyColor[i].R = byte(val)
	keyImages[i].SetDefaultBackground(keyColor[i], tcod.COLOR_BLACK)
	initColors()
}

func changeColorMapGreenCbk(w tcod.IWidget, val float32, data interface{}) {
	i := data.(int)
	keyColor[i].G = byte(val)
	keyImages[i].SetDefaultBackground(keyColor[i], tcod.COLOR_BLACK)
	initColors()
}

func changeColorMapBlueCbk(w tcod.IWidget, val float32, data interface{}) {
	i := data.(int)
	keyColor[i].B = byte(val)
	keyImages[i].SetDefaultBackground(keyColor[i], tcod.COLOR_BLACK)
	initColors()
}

func changeColorMapEndCbk(w tcod.IWidget, data interface{}) {
	colorMapGui.SetVisible(false)
}

func changeColorMapCbk(w tcod.IWidget, data interface{}) {
	colorMapGui.Move(w.GetX()+w.GetWidth()+2, w.GetY())
	colorMapGui.Clear()
	for i := 0; i < nbColorKeys; i++ {
		tmp := fmt.Sprintf("Color %d", i)
		colorMapGui.AddSeparator(tmp)
		hbox := gui.NewHBox(0, 0, 0)
		vbox := gui.NewVBox(0, 0, 0)
		colorMapGui.AddWidget(hbox)
		idxSlider := gui.NewSlider(0, 0, 3, 0.0, 255.0, "index", "Index of the key in the color map (0-255)")
		idxSlider.SetValue(float32(keyIndex[i]))
		idxSlider.SetFormat("%.0f")
		idxSlider.SetCallback(changeColorMapIdxCbk, i)
		vbox.AddWidget(idxSlider)
		keyImages[i] = gui.NewImageWidget(0, 0, 0, 2)
		keyImages[i].SetDefaultBackground(keyColor[i], keyColor[i])
		vbox.AddWidget(keyImages[i])
		hbox.AddWidget(vbox)

		vbox = gui.NewVBox(0, 0, 0)
		hbox.AddWidget(vbox)

		redSlider := gui.NewSlider(0, 0, 3, 0.0, 255.0, "r", "Red component of the color")
		redSlider.SetValue(float32(keyColor[i].R))
		redSlider.SetFormat("%.0f")
		redSlider.SetCallback(changeColorMapRedCbk, i)
		vbox.AddWidget(redSlider)
		greenSlider := gui.NewSlider(0, 0, 3, 0.0, 255.0, "g", "Green component of the color")
		greenSlider.SetValue(float32(keyColor[i].G))
		greenSlider.SetFormat("%.0f")
		greenSlider.SetCallback(changeColorMapGreenCbk, i)
		vbox.AddWidget(greenSlider)
		blueSlider := gui.NewSlider(0, 0, 3, 0.0, 255.0, "b", "Blue component of the color")
		blueSlider.SetValue(float32(keyColor[i].B))
		blueSlider.SetFormat("%.0f")
		blueSlider.SetCallback(changeColorMapBlueCbk, i)
		vbox.AddWidget(blueSlider)
	}
	colorMapGui.AddWidget(gui.NewButton("Ok", "", changeColorMapEndCbk, nil))
	colorMapGui.SetVisible(true)
}

func addOperationButton(tools *tcod.ToolBar, opType OpType, callback tcod.WidgetCallback) {
	tools.AddWidget(gui.NewButton(operations.names[opType], operations.tips[opType], callback, nil))
}

func buildGui() {
	// status bar
	gui.NewStatusBarDim(0, 0, HM_WIDTH, 1)

	vbox := gui.NewVBox(0, 2, 1)
	// stats
	stats := gui.NewToolBarWithWidth(0, 0, 21, "Stats", "Statistics about the current map")
	landMassLabel = gui.NewLabelWithTip(0, 0, "", "Ratio of land surface / total surface")
	minZLabel = gui.NewLabelWithTip(0, 0, "", "Minimum z value in the map")
	maxZLabel = gui.NewLabelWithTip(0, 0, "", "Maximum z value in the map")
	seedLabel = gui.NewLabelWithTip(0, 0, "", "Current random seed used to build the map")

	stats.AddWidget(landMassLabel)
	stats.AddWidget(minZLabel)
	stats.AddWidget(maxZLabel)
	stats.AddWidget(seedLabel)
	vbox.AddWidget(stats)

	// tools
	tools := gui.NewToolBarWithWidth(0, 0, 15, "Tools", "Tools to modify the heightmap")
	tools.AddWidget(gui.NewButton("cancel", "Delete the selected operation", cancelCbk, nil))
	tools.AddWidget(gui.NewButton("clear", "Remove all operations and reset all heightmap values to 0.0", clearCbk, nil))
	tools.AddWidget(gui.NewButton("reseed", "Replay all operations with a new random seed", reseedCbk, nil))

	// operations
	tools.AddSeparatorWithTip("Operations", "Apply a new operation to the map")
	addOperationButton(tools, NORM, normalizeCbk)
	addOperationButton(tools, ADDFBM, addFbmCbk)
	addOperationButton(tools, SCALEFBM, scaleFbmCbk)
	addOperationButton(tools, ADDHILL, addHillCbk)
	addOperationButton(tools, RAIN, rainErosionCbk)
	addOperationButton(tools, SMOOTH, smoothCbk)
	addOperationButton(tools, VORONOI, voronoiCbk)
	addOperationButton(tools, NOISELERP, noiseLerpCbk)
	addOperationButton(tools, ADDLEVEL, raiseLowerCbk)

	// display
	tools.AddSeparatorWithTip("Display", "Change the type of display")
	gui.SetDefaultRadioGroup(1)
	colormap := gui.NewRadioButton("colormap", "Enable colormap mode", colorMapCbk, nil)
	tools.AddWidget(colormap)

	colormap.Select()

	tools.AddWidget(gui.NewRadioButton("slope", "Enable slope mode", slopeCbk, nil))
	tools.AddWidget(gui.NewRadioButton("greyscale", "Enable greyscale mode", greyscaleCbk, nil))
	tools.AddWidget(gui.NewRadioButton("normal", "Enable normal map mode", normalCbk, nil))
	tools.AddWidget(gui.NewButton("change colormap", "Modify the colormap used by hmtool", changeColorMapCbk, nil))

	// change colormap gui
	colorMapGui = gui.NewToolBar(0, 0, "Colormap", "Select the color and position of the keys in the color map")
	colorMapGui.SetVisible(false)

	// in/out
	tools.AddSeparatorWithTip("In/Out", "Import/Export stuff")
	tools.AddWidget(gui.NewButton("export C", "Generate the C code for self heightmap in ./hm.c", exportCCbk, nil))
	tools.AddWidget(gui.NewButton("export CPP", "Generate the CPP code for self heightmap in ./hm.cpp", exportCppCbk, nil))
	tools.AddWidget(gui.NewButton("export PY", "Generate the python code for self heightmap in ./hm.py", exportPyCbk, nil))
	tools.AddWidget(gui.NewButton("export GO", "Generate the Go code for self heightmap in ./hm.go", exportGoCbk, nil))
	tools.AddWidget(gui.NewButton("export bmp", "Save self heightmap as a bitmap in ./hm.bmp", exportBmpCbk, nil))

	vbox.AddWidget(tools)

	// params box
	params = gui.NewToolBar(0, 0, "Params", "Parameters of the current tool")
	vbox.AddWidget(params)
	params.SetVisible(false)

	// history
	history = gui.NewToolBarWithWidth(0, tools.GetY()+1+tools.GetHeight(), 15, "History", "History of operations")
	vbox.AddWidget(history)
}

func main() {
	// change dir to program dir
	program := os.Args[0]
	dir, _ := path.Split(program)
	os.Chdir(dir)

	root = tcod.NewRootConsole(HM_WIDTH, HM_HEIGHT, "height map tool", false)
	guicon = tcod.NewConsole(HM_WIDTH, HM_HEIGHT)
	gui = tcod.NewGui(guicon)
	guicon.SetKeyColor(tcod.Color{255, 0, 255})
	backupMap = tcod.NewHeightMap(HM_WIDTH, HM_HEIGHT)

	tcod.SysSetFps(25)
	initColors()
	buildGui()
	hm = tcod.NewHeightMap(HM_WIDTH, HM_HEIGHT)
	hmold = tcod.NewHeightMap(HM_WIDTH, HM_HEIGHT)
	rnd = tcod.NewRandomFromSeed(seed)
	noise = tcod.NewNoise(2, rnd)
	var fade byte = 50
	var creditsEnd bool = false

	for !root.IsWindowClosed() {
		render()
		guicon.SetDefaultBackground(tcod.Color{255, 0, 255})
		guicon.Clear()
		gui.RenderWidgets()
		if !creditsEnd {
			creditsEnd = root.RenderCredits(HM_WIDTH-20, HM_HEIGHT-7, true)
		}
		if gui.GetFocusedWidget() != nil {
			if fade < 200 {
				fade += 20
			}
		} else {
			if fade > 80 {
				fade -= 20
			}

		}
		guicon.Blit(0, 0, HM_WIDTH, HM_HEIGHT, root, 0, 0, float32(fade)/255.0, float32(fade)/255.0)
		root.Flush()
		key := root.CheckForKeypress(tcod.KEY_PRESSED)
		gui.UpdateWidgets(key)
		switch key.C {
		case '+':
			(NewAddLevelOperation((mapmax - mapmin) / 50)).runInternal()
		case '-':
			(NewAddLevelOperation(-(mapmax - mapmin) / 50)).runInternal()
		default:
		}
		switch key.Vk {
		case tcod.K_PRINTSCREEN:
			tcod.SysSaveScreenshot()
		default:
		}
	}

}
