package tcod

import (
	"fmt"
	"container/vector"
	"strings"
	"strconv"
)

//
// Misc functions
//
func min(a, b int) int {
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

func absf(v float) float {
	if v < 0 {
		return -v
	}
	return v
}

func assertEqual(a interface{}, b interface{}) {
	if a != b {
		panic(fmt.Sprintf("Assertion error: a != b %v %v", a, b))
	}
}

// Inserts byte c in string a given position
// and returns new string
func insert(s string, c int, pos int) string {
	if pos < len(s) {
		return s[0:pos] + string(c) + s[pos:len(s)]
	}
	return s + string(c)
}

func replace(s string, c int, pos int) string {
	// if pos is beyond string length, then extend s
	if len(s) <= pos {
		return s + strings.Repeat(" ", pos-len(s)) + string(c)
	}
	return s[0:pos] + string(c) + s[(pos+1):len(s)]
}


// deletes string char at given position
// and returns new string
func delete(s string, pos int) string {
	if pos < len(s) {
		return s[0:pos] + s[pos+1:len(s)]
	}
	return s
}

func padRight(s string, length int, c int) string {
	if len(s) > length {
		return s[0:length]
	}
	return s + strings.Repeat(string(c), length-len(s))
}

//
// Generic widget interface
//

type IWidget interface {
	IsVisible() bool
	ComputeSize()
	Update(w IWidget, k Key)
	Render(w IWidget)
	GetX() int
	SetX(x int)
	GetY() int
	SetY(y int)
	GetWidth() int
	SetWidth(w int)
	GetHeight() int
	SetHeight(h int)
	GetUserData() interface{}
	SetUserData(data interface{})
	GetTip() string
	SetTip(tip string)
	GetMouseIn() bool
	SetMouseIn(mouseIn bool)
	GetMouseL() bool
	SetMouseL(mouseL bool)
	GetVisible() bool
	SetVisible(visible bool)
	SetGui(*Gui)
	GetGui() *Gui
	SetBackgroundColor(col, colFocus Color)
	SetForegroundColor(col, colFocus Color)
	GetBackgroundColor() (col, colFocus Color)
	GetForegroundColor() (col, colFocus Color)
	GetCurrentColors() (fore, back Color)
	onMouseIn()
	onMouseOut()
	onButtonPress()
	onButtonRelease()
	onButtonClick()
	expand(x, y int)
}


//
// WidgetVector collection
//
type WidgetVector struct {
	d vector.Vector
}

func NewWidgetVector() *WidgetVector {
	return &WidgetVector{vector.Vector{}}
}

func (self *WidgetVector) Push(w IWidget) {
	self.d.Push(w)
}

func (self *WidgetVector) Clear() {
	self.d = vector.Vector{}
}

func (self *WidgetVector) At(index int) IWidget {
	return self.d.At(index).(IWidget)
}

func (self *WidgetVector) Len() int {
	return self.d.Len()
}

func (self *WidgetVector) Remove(w IWidget) {
	for i, x := range self.d {
		if x.(IWidget) == w {
			self.d.Delete(i)
			break
		}
	}
}

// Iterate over all elements; driver for range
func (self *WidgetVector) iterate(c chan<- IWidget) {
	for _, v := range self.d {
		c <- v.(IWidget)
	}
	close(c)
}


// Channel iterator for range.
func (self *WidgetVector) Iter() <-chan IWidget {
	c := make(chan IWidget)
	go self.iterate(c)
	return c
}


//
// GUI info
//

type Gui struct {
	focus         IWidget // focused widget
	keyboardFocus IWidget // keyboard focused widget
	mouse         Mouse
	elapsed       float
	con           IConsole
	widgetVector  *WidgetVector
	rbs           *RadioButtonStatic
	tbs           *TextBoxStatic
}


func NewGui(console IConsole) *Gui {
	return &Gui{
		con:          console,
		widgetVector: NewWidgetVector(),
		tbs:          NewTextBoxStatic(),
		rbs:          NewRadioButtonStatic()}
}

func (self *Gui) Register(w IWidget) {
	w.SetGui(self)
	self.widgetVector.Push(w)
}

func (self *Gui) Unregister(w IWidget) {
	if self.focus == w {
		self.focus = nil
	}
	if self.keyboardFocus == w {
		self.keyboardFocus = nil
	}
	self.widgetVector.Remove(w)
}

func (self *Gui) updateWidgetsIntern(k Key) {
	self.elapsed = SysGetLastFrameLength()
	for e := range self.widgetVector.Iter() {
		w := e.(IWidget)
		if w.IsVisible() {
			w.ComputeSize()
			w.Update(w, k)
		}
	}
}

func (self *Gui) SetConsole(console IConsole) {
	self.con = console
}

func (self *Gui) UpdateWidgets(k Key) {
	self.mouse = MouseGetStatus()
	self.updateWidgetsIntern(k)
}

func (self *Gui) RenderWidgets() {
	for e := range self.widgetVector.Iter() {
		w := e.(IWidget)
		if w.IsVisible() {
			fore, back := self.con.GetForegroundColor(), self.con.GetBackgroundColor()
			w.Render(w)
			self.con.SetForegroundColor(fore)
			self.con.SetBackgroundColor(back)
		}
	}
}

func (self *Gui) IsFocused(w IWidget) bool {
	return self.focus == w
}

func (self *Gui) IsKeyboardFocused(w IWidget) bool {
	return self.keyboardFocus == w
}

func (self *Gui) GetFocusedWidget() IWidget {
	return self.focus
}

func (self *Gui) GetFocusedKeyboardWidget() IWidget {
	return self.keyboardFocus
}


func (self *Gui) UnSelectRadioGroup(group int) {
	self.rbs.UnSelectGroup(group)
}


func (self *Gui) SetDefaultRadioGroup(group int) {
	self.rbs.SetDefaultGroup(group)
}

//
// Widget root class
//

type Widget struct {
	x, y, w, h int
	userData   interface{}
	tip        string
	mouseIn    bool
	mouseL     bool
	visible    bool
	back       Color
	fore       Color
	backFocus  Color
	foreFocus  Color
	gui        *Gui
}


type WidgetCallback func(w IWidget, userData interface{})

func (self *Gui) newWidget() *Widget {
	result := &Widget{}
	self.Register(result)
	return result
}

func (self *Gui) NewWidget() *Widget {
	result := self.newWidget()
	result.initializeWidget(0, 0, 0, 0)
	return result
}
//

func (self *Gui) NewWidgetAt(x, y int) *Widget {
	result := self.newWidget()
	result.initializeWidget(x, y, 0, 0)
	return result
}

func (self *Gui) NewWidgetDim(x, y, w, h int) *Widget {
	result := self.newWidget()
	result.initializeWidget(x, y, w, h)
	return result
}


// Multiple dispatch: self
func (self *Widget) initializeWidget(x, y, w, h int) {
	//assertEqual(self, iself)
	self.x = x
	self.y = y
	self.w = w
	self.h = h
	self.tip = ""
	self.mouseIn = false
	self.mouseL = false
	self.visible = true
	self.back = Color{40, 40, 120}
	self.fore = Color{220, 220, 180}
	self.backFocus = Color{70, 70, 130}
	self.foreFocus = Color{255, 255, 255}
}


func (self *Widget) Delete() {
	self.gui.Unregister(self)
}

func (self *Widget) GetGui() *Gui {
	return self.gui
}

func (self *Widget) SetGui(gui *Gui) {
	self.gui = gui
}

func (self *Widget) GetX() int {
	return self.x
}

func (self *Widget) SetX(x int) {
	self.x = x
}

func (self *Widget) GetY() int {
	return self.y
}

func (self *Widget) SetY(y int) {
	self.y = y
}

func (self *Widget) GetWidth() int {
	return self.w
}

func (self *Widget) SetWidth(w int) {
	self.w = w
}

func (self *Widget) GetHeight() int {
	return self.h
}

func (self *Widget) SetHeight(h int) {
	self.h = h
}

func (self *Widget) GetUserData() interface{} {
	return self.userData
}

func (self *Widget) SetUserData(data interface{}) {
	self.userData = data
}

func (self *Widget) GetTip() string {
	return self.tip
}

func (self *Widget) SetTip(tip string) {
	self.tip = tip
}

func (self *Widget) GetMouseIn() bool {
	return self.mouseIn
}

func (self *Widget) SetMouseIn(mouseIn bool) {
	self.mouseIn = mouseIn
}

func (self *Widget) GetMouseL() bool {
	return self.mouseL
}

func (self *Widget) SetMouseL(mouseL bool) {
	self.mouseL = mouseL
}

func (self *Widget) GetVisible() bool {
	return self.visible
}
func (self *Widget) SetVisible(visible bool) {
	self.visible = visible
}

func (self *Widget) SetBackgroundColor(col, colFocus Color) {
	self.back = col
	self.backFocus = colFocus
}

func (self *Widget) SetForegroundColor(col, colFocus Color) {
	self.fore = col
	self.foreFocus = colFocus
}

func (self *Widget) GetBackgroundColor() (col, colFocus Color) {
	return self.back, self.backFocus
}

func (self *Widget) GetForegroundColor() (col, colFocus Color) {
	return self.fore, self.foreFocus
}

func (self *Widget) GetCurrentColors() (fore, back Color) {
	return If(self.mouseIn, self.foreFocus, self.fore).(Color), 
			 If(self.mouseIn, self.backFocus, self.back).(Color)
}

// both self and iself denote the same object
// here we emulate inheritance in go
// function receives self as receiver and as first param in interface
func (self *Widget) Update(iself IWidget, k Key) {
	//assertEqual(self, iself)
	curs := MouseIsCursorVisible()
	g := self.gui

	if curs {
		if g.mouse.Cx >= iself.GetX() && g.mouse.Cx < iself.GetX()+iself.GetWidth() &&
			g.mouse.Cy >= iself.GetY() && g.mouse.Cy < iself.GetY()+iself.GetHeight() {
			if !iself.GetMouseIn() {
				iself.SetMouseIn(true)
				iself.onMouseIn()
			}
			if g.focus != iself {
				g.focus = iself
			}
		} else {
			if iself.GetMouseIn() {
				iself.SetMouseIn(false)
				iself.onMouseOut()
			}
			iself.SetMouseL(false)
			if iself == g.focus {
				g.focus = nil
			}
		}
	}
	if iself.GetMouseIn() || (!curs && iself == g.focus) {
		if g.mouse.LButton && !iself.GetMouseL() {
			iself.SetMouseL(true)
			iself.onButtonPress()
		} else if !g.mouse.LButton && iself.GetMouseL() {
			iself.onButtonRelease()
			g.keyboardFocus = nil
			if iself.GetMouseL() {
				iself.onButtonClick()
			}
			iself.SetMouseL(false)
		} else if g.mouse.LButtonPressed {
			g.keyboardFocus = nil
			iself.onButtonClick()
		}
	}
}


func (self *Widget) Move(x, y int) {
	self.x = x
	self.y = y
}

func (self *Widget) ComputeSize() {
	// abstract
}

func (self *Widget) Render(iself IWidget) {
	// abstract
}


func (self *Widget) IsVisible() bool {
	return self.visible
}

func (self *Widget) onMouseIn() {
	// abstract
}

func (self *Widget) onMouseOut() {
	// abstract
}

func (self *Widget) onButtonPress() {
	// abstract
}

func (self *Widget) onButtonRelease() {
	// abstract
}

func (self *Widget) onButtonClick() {
	// abstract
}

func (self *Widget) expand(x, y int) {
	// abstract
}


//
// Button
//

type Button struct {
	Widget
	pressed  bool
	label    string
	callback WidgetCallback
}


func (self *Gui) newButton() *Button {
	result := &Button{}
	self.Register(result)
	return result
}

func (self *Gui) NewButton(label string, tip string, callback WidgetCallback, userData interface{}) *Button {
	result := self.newButton()
	result.initializeButton(0, 0, 0, 0, label, tip, callback, userData)
	return result
}

func (self *Gui) NewButtonDim(x, y, width, height int, label string, tip string, callback WidgetCallback, userData interface{}) *Button {
	result := self.newButton()
	result.initializeButton(x, y, width, height, label, tip, callback, userData)
	return result
}

func (self *Button) initializeButton(x, y, width, height int, label string, tip string, callback WidgetCallback, userData interface{}) {
	self.Widget.initializeWidget(x, y, width, height)
	self.label = label
	self.tip = tip
	self.userData = userData
	self.callback = callback
	self.x = x
	self.y = y
	self.w = width
	self.h = height
}


func (self *Button) SetLabel(newLabel string) {
	self.label = newLabel
}

func (self *Button) IsPressed() bool {
	return self.pressed
}

func (self *Button) ComputeSize() {
	self.w = len(self.label) + 2
	self.h = 1
}

func (self *Button) Render(iself IWidget) {
	con := self.gui.con
	fore, back := iself.GetCurrentColors()
	con.SetForegroundColor(fore)
	con.SetBackgroundColor(back)
	con.PrintCenter(self.x+self.w/2, self.y, BKGND_NONE, self.label)
	if self.w > 0 && self.h > 0 {
		con.Rect(self.x, self.y, self.w, self.h, true, BKGND_SET)
	}
	if self.label != "" {
		if self.pressed && self.mouseIn {
			//con.PrintCenter(self.x+self.w/2, self.y, BKGND_NONE, "-%s-", self.label)
			con.PrintCenter(self.x+self.w/2, self.y, BKGND_NONE, "%s", self.label)
			//con.PrintLeft(self.x + 1, self.y, BKGND_NONE, self.label)
		} else {
			con.PrintCenter(self.x+self.w/2, self.y, BKGND_NONE, self.label)
			//con.PrintLeft(self.x + 1, self.y, BKGND_NONE, self.label)
		}
	}
}

func (self *Button) onButtonPress() {
	self.pressed = true
}

func (self *Button) onButtonRelease() {
	self.pressed = false
}

func (self *Button) onButtonClick() {
	if self.callback != nil {
		self.callback(self, self.userData)
	}
}

func (self *Button) expand(width, height int) {
	if self.w < width {
		self.w = width
	}
}


//
// Status bar
//
//
type StatusBar struct {
	Widget
}

func (self *Gui) newStatusBar() *StatusBar {
	result := &StatusBar{}
	self.Register(result)
	return result
}

func (self *Gui) NewStatusBar() *StatusBar {
	result := self.newStatusBar()
	result.initializeStatusBar(0, 0, 0, 0)
	return result
}

func (self *Gui) NewStatusBarDim(x, y, w, h int) *StatusBar {
	result := self.newStatusBar()
	result.initializeStatusBar(x, y, w, h)
	return result

}

func (self *StatusBar) initializeStatusBar(x, y, w, h int) {
	self.initializeWidget(x, y, w, h)
}

func (self *StatusBar) Render(iself IWidget) {
	con := self.gui.con
	focus := self.gui.focus
	con.SetBackgroundColor(self.back)
	con.Rect(self.x, self.y, self.w, self.h, true, BKGND_SET)
	if focus != nil && focus.GetTip() != "" {
		con.SetForegroundColor(self.fore)
		con.PrintLeftRect(self.x+1, self.y, self.w, self.h, BKGND_NONE, focus.GetTip())
	}
}


//
//
//
//
// Image
//
//
type ImageWidget struct {
	Widget
	back Color
}


func (self *Gui) newImageWidget() *ImageWidget {
	result := &ImageWidget{}
	self.Register(result)
	return result
}

func (self *Gui) NewImageWidget(x, y, w, h int) *ImageWidget {
	result := self.newImageWidget()
	result.initializeImageWidget(x, y, w, h, "")
	return result
}

func (self *Gui) NewImageWidgetWithTip(x, y, w, h int, tip string) *ImageWidget {
	result := self.newImageWidget()
	result.initializeImageWidget(x, y, w, h, tip)
	return result
}

func (self *ImageWidget) initializeImageWidget(x, y, w, h int, tip string) {
	self.Widget.initializeWidget(x, y, w, h)
	self.tip = tip
	self.back = COLOR_BLACK
}


func (self *ImageWidget) Render(iself IWidget) {
	con := self.gui.con
	con.SetBackgroundColor(self.back)
	con.Rect(self.x, self.y, self.w, self.h, true, BKGND_SET)

}


func (self *ImageWidget) expand(width, height int) {
	if width > self.w {
		self.w = width
	}
	if height > self.h {
		self.h = height
	}
}


//
//
// Container
//
//
//

type Container struct {
	Widget
	content *WidgetVector
}


func (self *Gui) newContainer() *Container {
	result := &Container{}
	self.Register(result)
	return result
}

func (self *Gui) NewContainer(x, y, w, h int) *Container {
	result := self.newContainer()
	result.initializeContainer(x, y, w, h)
	return result
}


func (self *Container) initializeContainer(x, y, w, h int) {
	self.Widget.initializeWidget(x, y, w, h)
	self.content = &WidgetVector{}
}


func (self *Container) AddWidget(w IWidget) {
	self.content.Push(w)
	self.gui.Unregister(w)
}

func (self *Container) RemoveWidget(w IWidget) {
	self.content.Remove(w)
}

func (self *Container) Render(iself IWidget) {
	for w := range self.content.Iter() {
		if w.IsVisible() {
			w.Render(w)
		}
	}
}

func (self *Container) Clear() {
	self.content.Clear()
}

func (self *Container) Update(iself IWidget, k Key) {
	self.Widget.Update(iself, k)

	for w := range self.content.Iter() {
		if w.IsVisible() {
			w.Update(w, k)
		}
	}
}


//
//
// VBox
//
type VBox struct {
	Container
	padding int
}


func (self *Gui) newVBox() *VBox {
	result := &VBox{}
	self.Register(result)
	return result
}


func (self *Gui) NewVBox(x, y, padding int) *VBox {
	result := self.newVBox()
	result.initializeVBox(x, y, padding)
	return result

}

func (self *VBox) initializeVBox(x, y, padding int) {
	self.Container.initializeContainer(x, y, 0, 0)
	self.padding = padding
}


func (self *VBox) ComputeSize() {
	cury := self.y
	self.w = 0
	for w := range self.content.Iter() {
		if w.IsVisible() {
			w.SetX(self.x)
			w.SetY(cury)
			w.ComputeSize()
			if w.GetWidth() > self.w {
				self.w = w.GetWidth()
			}
			cury += w.GetHeight() + self.padding
		}
	}
	self.h = cury - self.y

	for w := range self.content.Iter() {
		if w.IsVisible() {
			w.expand(self.w, w.GetHeight())
		}
	}
}


//
//
// HBox
//

type HBox struct {
	VBox
}


func (self *Gui) newHBox() *HBox {
	result := &HBox{}
	self.Register(result)
	return result
}

func (self *Gui) NewHBox(x, y, padding int) *HBox {
	result := self.newHBox()
	result.initializeHBox(x, y, padding)
	return result
}

func (self *HBox) initializeHBox(x, y, padding int) {
	self.VBox.initializeVBox(x, y, padding)
}


func (self *HBox) ComputeSize() {
	curx := self.x
	self.h = 0
	for w := range self.content.Iter() {
		if w.IsVisible() {
			w.SetY(self.y)
			w.SetX(curx)
			w.ComputeSize()
			if w.GetHeight() > self.h {
				self.h = w.GetHeight()
			}
			curx += w.GetWidth() + self.padding
		}
	}

	self.w = curx - self.x
	for w := range self.content.Iter() {
		if w.IsVisible() {
			w.expand(w.GetWidth(), self.h)
		}
	}
}


//
//
// Toolbar
//

//
// Separator
//
//
type Separator struct {
	Widget
	txt string
}

func (self *Gui) newSeparator() *Separator {
	result := &Separator{}
	self.Register(result)
	return result
}

func (self *Gui) NewSeparator(txt string) *Separator {
	result := self.newSeparator()
	result.initializeSeparator(txt, "")
	return result
}

func (self *Gui) NewSeparatorWithTip(txt, tip string) *Separator {
	result := self.newSeparator()
	result.initializeSeparator(txt, tip)
	return result

}

func (self *Separator) initializeSeparator(txt, tip string) {
	self.Widget.initializeWidget(0, 0, 0, 1)
	self.txt = txt

}


func (self *Separator) ComputeSize() {
	self.w = If(self.txt != "", len(self.txt)+2, 0).(int)
}

func (self *Separator) expand(width, height int) {
	if self.w < width {
		self.w = width
	}
}

func (self *Separator) Render(iself IWidget) {
	con := self.gui.con
	con.SetBackgroundColor(self.back)
	con.SetForegroundColor(self.fore)
	con.Hline(self.x, self.y, self.w, BKGND_SET)
	con.SetChar(self.x-1, self.y, CHAR_TEEE)
	con.SetChar(self.x+self.w, self.y, CHAR_TEEW)
	con.SetBackgroundColor(self.fore)
	con.SetForegroundColor(self.back)
	con.PrintCenter(self.x+self.w/2, self.y, BKGND_SET, " %s ", self.txt)
}


type ToolBar struct {
	Container
	name       string
	fixedWidth int
	shouldPrintFrame bool
}


func (self *Gui) newToolBar() *ToolBar {
	result := &ToolBar{}
	self.Register(result)
	return result
}


func (self *Gui) NewToolBarWithWidth(x, y, w int, name, tip string) *ToolBar {
	result := self.newToolBar()
	result.initializeToolBar(x, y, w, name, tip)
	return result

}

func (self *Gui) NewToolBar(x, y int, name, tip string) *ToolBar {
	result := self.newToolBar()
	result.initializeToolBar(x, y, 0, name, tip)
	return result

}


func (self *ToolBar) initializeToolBar(x, y, w int, name, tip string) {
	self.Container.initializeContainer(x, y, w, 2)
	self.name = name
	self.tip = tip
	if w == 0 {
		self.w = len(name) + 4
		self.fixedWidth = 0
	} else {
		self.w = max(len(name)+4, w)
		self.fixedWidth = max(len(name)+4, w)
	}
	self.shouldPrintFrame = true

}

func (self *ToolBar) SetShouldPrintFrame(value bool) {
	self.shouldPrintFrame = value
}

func (self *ToolBar) GetShouldPrintFrame() bool {
	return self.shouldPrintFrame
}


func (self *ToolBar) Render(iself IWidget) {
	con := self.gui.con
   fore, back := iself.GetCurrentColors()
	con.SetForegroundColor(fore)
	con.SetBackgroundColor(back)
	if self.shouldPrintFrame {
		con.PrintFrame(self.x, self.y, self.w, self.h, true, BKGND_SET, self.name)
	}
	self.Container.Render(iself)
}

func (self *ToolBar) SetName(name string) {
	self.name = name
	self.fixedWidth = max(len(name)+4, self.fixedWidth)

}

func (self *ToolBar) AddSeparator(txt string) {
	self.AddWidget(self.gui.NewSeparator(txt))
}

func (self *ToolBar) AddSeparatorWithTip(txt string, tip string) {
	self.AddWidget(self.gui.NewSeparatorWithTip(txt, tip))
}

func (self *ToolBar) ComputeSize() {
	cury := self.y + 1
	self.w = If(self.name != "", len(self.name)+4, 2).(int)
	for w := range self.content.Iter() {
		if w.IsVisible() {
			w.SetX(self.x + 1)
			w.SetY(cury)
			w.ComputeSize()
			if w.GetWidth()+2 > self.w {
				self.w = w.GetWidth() + 2
			}
			cury += w.GetHeight()
		}
	}
	if self.w < self.fixedWidth {
		self.w = self.fixedWidth
	}
	self.h = cury - self.y + 1
	for w := range self.content.Iter() {
		if w.IsVisible() {
			w.expand(self.w-2, w.GetHeight())
		}
	}
}



//
// ToggleButton
//


type ToggleButton struct {
	Button
	pressed bool
}

func (self *Gui) newToggleButton() *ToggleButton {
	result := &ToggleButton{}
	self.Register(result)
	return result
}

func (self *Gui) NewToggleButton(label, tip string, callback WidgetCallback, userData interface{}) *ToggleButton {
	result := self.newToggleButton()
	result.initializeToggleButton(0, 0, 0, 0, label, tip, callback, userData)
	return result
}

func (self *Gui) NewToggleButtonWithTip(x, y, width, height int, label, tip string, callback WidgetCallback, userData interface{}) *ToggleButton {
	result := self.newToggleButton()
	result.initializeToggleButton(x, y, width, height, label, tip, callback, userData)
	return result
}

func (self *ToggleButton) initializeToggleButton(x, y, width, height int, label string, tip string, callback WidgetCallback, userData interface{}) {
	self.Button.initializeButton(x, y, width, height, label, tip, callback, userData)
}

func (self *ToggleButton) IsPressed() bool {
	return self.pressed
}

func (self *ToggleButton) SetPressed(val bool) {
	self.pressed = val
}


func (self *ToggleButton) onButtonClick() {
	self.pressed = !self.pressed
	if self.callback != nil {
		self.callback(self, self.userData)
	}
}

func (self *ToggleButton) Render(iself IWidget) {
	con := self.gui.con

	fore, back := iself.GetCurrentColors()
	con.SetBackgroundColor(back)
	con.SetForegroundColor(fore)
	con.Rect(self.x, self.y, self.w, self.h, true, BKGND_SET)
	if self.label != "" {
		con.PrintLeft(self.x, self.y, BKGND_NONE, "%c %s",
			If(self.pressed, CHAR_CHECKBOX_SET, CHAR_CHECKBOX_UNSET).(int), self.label)
	} else {
		con.PrintLeft(self.x, self.y, BKGND_NONE, "%c",
			If(self.pressed, CHAR_CHECKBOX_SET, CHAR_CHECKBOX_UNSET).(int), self.label)
	}
}


//
//
// Label
//

type Label struct {
	Widget
	label string
}


func (self *Gui) newLabel() *Label {
	result := &Label{}
	self.Register(result)
	return result
}


func (self *Gui) NewLabel(x, y int, label string) *Label {
	result := self.newLabel()
	result.initializeLabel(x, y, label, "")
	return result
}

func (self *Gui) NewLabelWithTip(x, y int, label string, tip string) *Label {
	result := self.newLabel()
	result.initializeLabel(x, y, label, tip)
	return result
}


func (self *Label) initializeLabel(x, y int, label string, tip string) {
	self.Widget.initializeWidget(x, y, 0, 1)
	self.x = x
	self.y = y
	self.label = label
	self.tip = tip
}


func (self *Label) Render(iself IWidget) {
	con := self.gui.con
	con.SetBackgroundColor(self.back)
	con.SetForegroundColor(self.fore)
	con.PrintLeft(self.x, self.y, BKGND_NONE, self.label)
}

func (self *Label) ComputeSize() {
	self.w = len(self.label)
}

func (self *Label) SetValue(label string) {
	self.label = label
}

func (self *Label) expand(width, height int) {
	if self.w < width {
		self.w = width
	}
}


//
//
// TextBox
//

type TextBoxCallback func(w IWidget, val string, data interface{})

type TextBoxStatic struct {
	blinkingDelay float
}


func NewTextBoxStatic() *TextBoxStatic {
	return &TextBoxStatic{
		blinkingDelay: 0.5}
}


type TextBox struct {
	Widget
	label            string
	txt              string
	blink            float
	pos, offset      int
	boxx, boxw, maxw int
	insert           bool
	callback         TextBoxCallback
	data             interface{}
}


func (self *Gui) newTextBox() *TextBox {
	result := &TextBox{}
	self.Register(result)
	return result
}

func (self *Gui) NewTextBox(x, y, w, maxw int, label, value string) *TextBox {
	result := self.newTextBox()
	result.initializeTextBox(x, y, w, maxw, label, value, "")
	return result
}

func (self *Gui) NewTextBoxWithTip(x, y, w, maxw int, label, value, tip string) *TextBox {
	result := self.newTextBox()
	result.initializeTextBox(x, y, w, maxw, label, value, tip)
	return result
}

func (self *TextBox) initializeTextBox(x, y, w, maxw int, label, value, tip string) {
	self.Widget.initializeWidget(x, y, w, 0)
	self.x = x
	self.y = y
	self.w = w
	self.h = 1
	self.maxw = maxw
	self.label = label
	if len(value) > maxw {
		self.txt = value[0:maxw]
	} else {
		self.txt = value
	}
	self.tip = tip
	self.boxw = w
	if label != "" {
		self.boxx = len(label) + 1
		self.w += self.boxx
	}
}


func (self *TextBox) Render(iself IWidget) {
	// save colors
	con := self.gui.con
	g := self.gui

	con.SetBackgroundColor(self.back)
	con.SetForegroundColor(self.fore)
	con.Rect(self.x, self.y, self.w, self.h, true, BKGND_SET)
	if self.label != "" {
		con.PrintLeft(self.x, self.y, BKGND_NONE, self.label)
	}

	con.SetBackgroundColor(If(g.IsKeyboardFocused(self), self.foreFocus, self.fore).(Color))
	con.SetForegroundColor(If(g.IsKeyboardFocused(self), self.backFocus, self.back).(Color))
	con.Rect(self.x+self.boxx, self.y, self.boxw, self.h, false, BKGND_SET)
	length := len(self.txt) - self.offset
	if length > self.boxw {
		length = self.boxw
	}
	if self.txt != "" {
		con.PrintLeft(self.x+self.boxx, self.y, BKGND_NONE, padRight(self.txt[self.offset:], length, ' '))
	}
	if g.IsKeyboardFocused(self) && self.blink > 0.0 {
		if self.insert {
			con.SetBack(self.x+self.boxx+self.pos-self.offset, self.y, self.fore, BKGND_SET)
			con.SetFore(self.x+self.boxx+self.pos-self.offset, self.y, self.back)
		} else {
			con.SetBack(self.x+self.boxx+self.pos-self.offset, self.y, self.back, BKGND_SET)
			con.SetFore(self.x+self.boxx+self.pos-self.offset, self.y, self.fore)
		}
	}
}


func (self *TextBox) Update(iself IWidget, k Key) {
	g := self.gui
	tbs := self.gui.tbs
	if g.keyboardFocus == IWidget(self) {
		self.blink -= g.elapsed
		if self.blink < -tbs.blinkingDelay {
			self.blink += 2 * tbs.blinkingDelay
		}
		if k.Vk == K_SPACE || k.Vk == K_CHAR ||
			(k.Vk >= K_0 && k.Vk <= K_9) ||
			(k.Vk >= K_KP0 && k.Vk <= K_KP9) {
			if !self.insert || len(self.txt) < self.maxw {
				if self.insert && self.pos < len(self.txt) {
					self.txt = insert(self.txt, int(k.C), self.pos)
				} else {
					self.txt = replace(self.txt, int(k.C), self.pos)
				}
				if self.pos < self.maxw {
					self.pos++
				}
				if self.pos >= self.boxw {
					self.offset = self.pos - self.boxw + 1
				}
				if self.callback != nil {
					self.callback(self, self.txt, self.data)
				}
			}
			self.blink = tbs.blinkingDelay
		}
		switch k.Vk {
		case K_LEFT:
			if self.pos > 0 {
				self.pos--
			}
			if self.pos < self.offset {
				self.offset = self.pos
			}
			self.blink = tbs.blinkingDelay
		case K_RIGHT:
			if self.pos < len(self.txt) {
				self.pos++
			}
			if self.pos >= self.boxw {
				self.offset = self.pos - self.boxw + 1
			}
			self.blink = tbs.blinkingDelay
		case K_HOME:
			self.pos, self.offset = 0, 0
			self.blink = tbs.blinkingDelay
		case K_BACKSPACE:
			if self.pos > 0 {
				self.pos--
				self.txt = delete(self.txt, self.pos)
				if self.callback != nil {
					self.callback(self, self.txt, self.data)
				}
				if self.pos < self.offset {
					self.offset = self.pos
				}
			}
			self.blink = tbs.blinkingDelay
		case K_DELETE:
			if self.pos < len(self.txt) {
				self.txt = delete(self.txt, self.pos)
				if self.callback != nil {
					self.callback(self, self.txt, self.data)
				}
			}
			self.blink = tbs.blinkingDelay
		case K_END:
			self.pos = len(self.txt)
			if self.pos >= self.boxw {
				self.offset = self.pos - self.boxw + 1
			}
			self.blink = tbs.blinkingDelay
		default:
		}
	}
	self.Widget.Update(iself, k)
}


func (self *TextBox) SetBlinkingDelay(delay float) {
	self.gui.tbs.blinkingDelay = delay
}

func (self *TextBox) GetBlinkingDelay() float {
	return self.gui.tbs.blinkingDelay
}

func (self *TextBox) SetText(txt string) {
	if self.maxw < len(txt) {
		self.txt = txt[0:self.maxw]
	} else {
		self.txt = txt
	}
}

func (self *TextBox) GetText() string {
	return self.txt
}

func (self *TextBox) SetCallback(callback TextBoxCallback, data interface{}) {
	self.callback = callback
	self.data = data
}


func (self *TextBox) onButtonClick() {
	g := self.gui
	if g.mouse.Cx >= self.x+self.boxx && g.mouse.Cx < self.x+self.boxx+self.boxw {
		g.keyboardFocus = self
	}
}


//
// RadioButton
//
//

type RadioButtonStatic struct {
	defaultGroup int
	groupSelect  [512]*RadioButton
	init         bool
}

func NewRadioButtonStatic() *RadioButtonStatic {
	return &RadioButtonStatic{
		defaultGroup: 0,
		init:         false,
		groupSelect:  [512]*RadioButton{}}
}


func (self *RadioButtonStatic) UnSelectGroup(group int) {
	self.groupSelect[group] = nil
}


func (self *RadioButtonStatic) SetDefaultGroup(group int) {
	self.defaultGroup = group
}


type RadioButton struct {
	Button
	foreSelection, backSelection Color 
	useSelectionColor bool
	group int
}


func (self *Gui) newRadioButton() *RadioButton {
	result := &RadioButton{}
	self.Register(result)
	return result
}


func (self *Gui) NewRadioButton(label string, tip string, callback WidgetCallback, userData interface{}) *RadioButton {
	result := self.newRadioButton()
	result.initializeRadioButton(0, 0, 0, 0, label, tip, callback, userData)
	return result
}

func (self *Gui) NewRadioButtonWithTip(x, y, width, height int, label string, tip string, callback WidgetCallback, userData interface{}) *RadioButton {
	result := self.newRadioButton()
	result.initializeRadioButton(x, y, width, height, label, tip, callback, userData)
	return result
}


func (self *RadioButton) initializeRadioButton(x, y, width, height int, label string, tip string, callback WidgetCallback, userData interface{}) {
	self.Button.initializeButton(x, y, width, height, label, tip, callback, userData)
}

func (self *RadioButton) SetGroup(group int) {
	self.group = group
}


func (self *RadioButton) SetUseSelectionColor(use bool) {
	self.useSelectionColor = use
}

func (self *RadioButton) GetUseSelectionColor() bool {
	return self.useSelectionColor
}

func (self *RadioButton) SetSelectionColor(fore, back Color) {
	self.foreSelection, self.backSelection = fore, back
}

func (self *RadioButton) GetSelectionColor() (fore, back Color) {
	return self.foreSelection, self.backSelection
}

func (self *RadioButton) GetCurrentColors() (fore, back Color) {
	fore, back = self.Button.GetCurrentColors()
	if self.useSelectionColor && self.IsSelected() {
		fore, back = self.foreSelection, self.backSelection
	}
	return
}

func (self *RadioButton) Render(iself IWidget) {
	con := self.gui.con
	fore, back := iself.GetCurrentColors()
	con.SetForegroundColor(fore)
	con.SetBackgroundColor(back)
	self.Button.Render(iself)
	if self.IsSelected() && !self.GetUseSelectionColor() {
		con.PutCharEx(self.x, self.y, '>', fore, back)
	}
}

func (self *RadioButton) IsSelected() bool {
	rbs := self.gui.rbs
	return rbs.groupSelect[self.group] == self
}

func (self *RadioButton) Select() {
	rbs := self.gui.rbs
	rbs.groupSelect[self.group] = self
}

func (self *RadioButton) UnSelect() {
	rbs := self.gui.rbs
	rbs.groupSelect[self.group] = nil
}


func (self *RadioButton) onButtonClick() {
	self.Select()
	self.Button.onButtonClick()
}


//
//
// Slider
//

type SliderCallback func(w IWidget, val float, data interface{})

type Slider struct {
	TextBox
	min, max     float
	value        float
	sensitivity  float
	onArrows     bool
	drag         bool
	dragx, dragy int
	dragValue    float
	fmt          string
	callback     SliderCallback
	data         interface{}
}


func (self *Gui) newSlider() *Slider {
	result := &Slider{}
	self.Register(result)
	return result
}

func (self *Gui) NewSlider(x, y, w int, min, max float, label string, tip string) *Slider {
	result := self.newSlider()
	result.initializeSlider(x, y, w, min, max, label, tip)
	return result
}


func (self *Slider) initializeSlider(x, y, w int, min, max float, label string, tip string) {
	self.TextBox.initializeTextBox(x, y, w, 10, label, "", tip)
	self.min = min
	self.max = max
	self.value = (min + max) * 0.5
	self.sensitivity = 1.0
	self.onArrows = false
	self.drag = false
	self.fmt = ""
	self.callback = nil
	self.data = nil
	self.valueToText()
	self.w += 2

}

func (self *Slider) GetCurrentColors() (fore, back Color) {
	fore = If(self.onArrows || self.drag, self.foreFocus, self.fore).(Color)
	back = If(self.onArrows || self.drag, self.backFocus, self.back).(Color)
	return 
}

func (self *Slider) Render(iself IWidget) {
	con := self.gui.con
	fore, back := iself.GetCurrentColors()
	con.SetBackgroundColor(back)
	con.SetForegroundColor(fore)
	self.w -= 2
	self.TextBox.Render(iself)
	self.w += 2
	con.Rect(self.x+self.w-2, self.y, 2, 1, true, BKGND_SET)

	con.PutCharEx(self.x+self.w-2, self.y, CHAR_ARROW_W, fore, back)
	con.PutCharEx(self.x+self.w-1, self.y, CHAR_ARROW_E, fore, back)
}

func (self *Slider) Update(iself IWidget, k Key) {
	con := self.gui.con
	mouse := self.gui.mouse
	oldValue := self.value
	self.TextBox.Update(iself, k)
	self.textToValue()

	if mouse.Cx >= self.x+self.w-2 && mouse.Cx < self.x+self.w && mouse.Cy == self.y {
		self.onArrows = true
	} else {
		self.onArrows = false
	}
	if self.drag {
		if self.dragy == -1 {
			self.dragx = mouse.X
			self.dragy = mouse.Y
		} else {
			mdx := (float(mouse.X-self.dragx) * self.sensitivity) / float(con.GetWidth()*8)
			mdy := (float(mouse.Y-self.dragy) * self.sensitivity) / float(con.GetHeight()*8)
			oldValue := self.value
			if absf(mdy) > absf(mdx) {
				mdx = -mdy
			}
			self.value = self.dragValue + (self.max-self.min)*mdx
			self.value = ClampF(self.min, self.max, self.value)
			if self.value != oldValue {
				self.valueToText()
				self.textToValue()
			}
		}
	}
	if self.value != oldValue && self.callback != nil {
		self.callback(self, self.value, self.data)
	}
}


func (self *Slider) SetMinMax(min, max float) {
	self.min = min
	self.max = max
}

func (self *Slider) SetCallback(callback SliderCallback, data interface{}) {
	self.callback = callback
	self.data = data
}

func (self *Slider) SetFormat(fmt string) {
	self.fmt = fmt
}

func (self *Slider) SetValue(value float) {
	self.value = ClampF(self.min, self.max, value)
	self.valueToText()
}

func (self *Slider) SetSensitivity(sensitivity float) {
	self.sensitivity = sensitivity
}


func (self *Slider) valueToText() {
	self.txt = fmt.Sprintf(If(self.fmt != "", self.fmt, "%.2f").(string), self.value)
}


func (self *Slider) textToValue() {
	f, err := strconv.Atof32(self.txt)
	if err != nil {
		self.value = 0
	} else {
		self.value = float(f)
	}
}

func (self *Slider) onButtonPress() {
	if self.onArrows {
		self.drag = true
		self.dragy = -1
		self.dragValue = self.value
		MouseShowCursor(false)
	}
}

func (self *Slider) onButtonRelease() {
	if self.drag {
		self.drag = false
		MouseMove((self.x+self.w-2)*8, self.y*8)
		MouseShowCursor(true)
	}
}
