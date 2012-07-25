package tcod

//
// This is go-bindings package for libtcod-go
// Most API is wrapped except:
// - custom containers - Go has it's own containers
// - threads, mutexes and semaphores - they are replaced by goroutines and channels
// - SDL renderer - currently Go has very cumbersome C callback mechanism
//

/*
 #include <stdio.h>
 #include <stdlib.h>
 #include <string.h>
 #include "libtcod.h"
 #include "libtcod_int.h"

 // This is a workaround for cgo disability to process varargs
 // These functions are copied verbatim from console_c and replaced ... with simple string
 // Formatting will be done in Go functions

 void _TCOD_console_print(TCOD_console_t con,int x, int y, char *s) {
 	TCOD_console_data_t *dat;
 	if (! con ) con=TCOD_ctx.root;
 	dat=(TCOD_console_data_t *)con;
	char *t = strdup(s);
 	TCOD_console_print(con,x,y,"%s",t);
	free(t);
 }


 void  _TCOD_console_print_ex(TCOD_console_t con,int x, int y, TCOD_bkgnd_flag_t flag, TCOD_alignment_t alignment, const
 char *s) {
 	TCOD_console_data_t *dat;
 	if (! con ) con=TCOD_ctx.root;
 	dat=(TCOD_console_data_t *)con;
	char *t = strdup(s);
 	TCOD_console_print_ex(con,x,y,flag,alignment,"%s",t);
	free(t);
 }

 int _TCOD_console_print_rect(TCOD_console_t con,int x, int y, int w, int h, char *s) {
 	return TCOD_console_print_rect(con,x,y,w,h,"%s",s);
 }


 int _TCOD_console_print_rect_ex(TCOD_console_t con,int x, int y, int w, int h, TCOD_bkgnd_flag_t flag, TCOD_alignment_t alignment, const char *s) {
	return TCOD_console_print_rect_ex(con, x, y, w, h, flag, alignment, "%s", s);
 }

 int _TCOD_console_height_rect(TCOD_console_t con,int x, int y, int w, int h, char *s) {
 	int ret;
 	ret = TCOD_console_get_height_rect(con,x,y,w,h, "%s", s);
 	return ret;
 }


 void _TCOD_console_print_frame(TCOD_console_t con,int x,int y,int w,int h, bool empty, TCOD_bkgnd_flag_t flag, char *s) {
 	TCOD_console_data_t *dat;
 	if (! con ) con=TCOD_ctx.root;
 	dat=(TCOD_console_data_t *)con;
 	TCOD_console_put_char(con,x,y,TCOD_CHAR_NW,flag);
 	TCOD_console_put_char(con,x+w-1,y,TCOD_CHAR_NE,flag);
 	TCOD_console_put_char(con,x,y+h-1,TCOD_CHAR_SW,flag);
 	TCOD_console_put_char(con,x+w-1,y+h-1,TCOD_CHAR_SE,flag);
 	TCOD_console_hline(con,x+1,y,w-2,flag);
 	TCOD_console_hline(con,x+1,y+h-1,w-2,flag);
 	if ( h > 2 ) {
 		TCOD_console_vline(con,x,y+1,h-2,flag);
 		TCOD_console_vline(con,x+w-1,y+1,h-2,flag);
 		if ( empty ) {
 			TCOD_console_rect(con,x+1,y+1,w-2,h-2,true,flag);
 		}
 	}
 	if (s) {
 		int xs;
 		TCOD_color_t tmp;
 		xs = x + (w-strlen(s)-2)/2;
 		tmp=dat->back; // swap colors
 		dat->back=dat->fore;
 		dat->fore=tmp;
 		TCOD_console_print(con,xs,y," %s ",s);
 		tmp=dat->back; // swap colors
 		dat->back=dat->fore;
 		dat->fore=tmp;
 	}
   }


  float _TCOD_heightmap_get_nth_value(const TCOD_heightmap_t *hm, int nth) {
  	return hm->values[nth];
  }

  void _TCOD_heightmap_set_nth_value(const TCOD_heightmap_t *hm, int nth, float val) {
  	hm->values[nth] = val;
  }

  // These functions are a workaround for passing structs from Go to C
  // I must pass them as pointers otherwise they are passed incorrectly
  // when more than one structure is passed

  bool _TCOD_color_equals(TCOD_color_t *c1, TCOD_color_t *c2) {
  	return TCOD_color_equals(*c1, *c2);
  }

  TCOD_color_t _TCOD_color_lerp(TCOD_color_t *c1, TCOD_color_t *c2, float coef) {
  	return TCOD_color_lerp(*c1, *c2, coef);
  }

  TCOD_color_t _TCOD_color_add(TCOD_color_t *c1, TCOD_color_t *c2) {
    return TCOD_color_add(*c1, *c2);
  }

  TCOD_color_t _TCOD_color_subtract(TCOD_color_t *c1, TCOD_color_t *c2) {
    return TCOD_color_subtract(*c1, *c2);
  }

  TCOD_color_t _TCOD_color_multiply(TCOD_color_t *c1, TCOD_color_t *c2) {
    return TCOD_color_multiply(*c1, *c2);
  }

  TCOD_color_t _TCOD_color_multiply_scalar(TCOD_color_t *c, float value) {
	return TCOD_color_multiply_scalar(*c, value);
  }

  void _TCOD_console_set_color_control(TCOD_colctrl_t ctrl, TCOD_color_t *fore, TCOD_color_t *back) {
	TCOD_console_set_color_control(ctrl, *fore, *back);
  }

  void _TCOD_console_set_key_color(TCOD_console_t console, TCOD_color_t *color) {
	TCOD_console_set_key_color(console, *color);
  }


  void _TCOD_console_set_char_background(TCOD_console_t console, int x, int y, TCOD_color_t *color, TCOD_bkgnd_flag_t flag) {
	TCOD_console_set_char_background(console, x, y, *color, flag);
  }

  void _TCOD_console_set_char_foreground(TCOD_console_t console, int x, int y, TCOD_color_t *color) {
	TCOD_console_set_char_foreground(console, x, y, *color);
  }

  void _TCOD_console_set_fade(uint8 val, TCOD_color_t *color) {
	TCOD_console_set_fade(val, *color);
  }

  void _TCOD_console_put_char_ex(TCOD_console_t console, int x, int y, int c, TCOD_color_t *fore, TCOD_color_t *back) {
	TCOD_console_put_char_ex(console, x, y, c, *fore, *back);
  }

  void _TCOD_text_set_colors(TCOD_text_t text, TCOD_color_t *fore, TCOD_color_t *back, float transparency) {
	TCOD_text_set_colors(text, *fore, *back, transparency);
  }

  typedef struct {
  	char *name;
  	TCOD_value_type_t value_type;
  	TCOD_value_t value;
  } _prop_t;

  void  _TCOD_struct_add_structure(TCOD_parser_struct_t *s, TCOD_parser_struct_t *subs) {
    TCOD_struct_add_structure(*s, *subs);
  }

  void _TCOD_image_put_pixel(TCOD_image_t img, int x, int y, TCOD_color_t *color) {
	 TCOD_image_put_pixel(img, x, y, *color);

  }

  void _TCOD_image_clear(TCOD_image_t img, TCOD_color_t *color) {
	TCOD_image_clear(img, *color);
  }

  void _TCOD_image_set_key_color(TCOD_image_t img, TCOD_color_t *color) {
	TCOD_image_set_key_color(img, *color);
  }

  void _TCOD_zip_put_color(TCOD_zip_t zip, TCOD_color_t *val) {
	 TCOD_zip_put_color(zip, *val);
  }

  bool _TCOD_sys_file_exists(const char * filename) {
	return TCOD_sys_file_exists(filename);
  }

*/
import "C"

import (
	"container/vector"
	"fmt"
	"runtime"
	"unsafe"
)

type void unsafe.Pointer

//
// Misc functions
//

type BkgndFlag C.TCOD_bkgnd_flag_t

func BkgndAlpha(alpha float32) BkgndFlag {
	return BkgndFlag(BKGND_ALPH | (((uint8)(alpha * 255)) << 8))
}

func BkgndAddAlpha(alpha float32) BkgndFlag {
	return BkgndFlag(BKGND_ADDA | (((uint8)(alpha * 255)) << 8))
}

func If(condition bool, tv, fv interface{}) interface{} {
	if condition {
		return tv
	} else {
		return fv
	}
	return nil
}

func Clamp(a, b, x int) int {
	return If(x < a, a, If(x > b, b, x).(int)).(int)
}

func ClampF(a, b, x float32) float32 {
	return If(x < a, a, If(x > b, b, x).(float32)).(float32)
}

//
// TODO should free those strings?
func toStringSlice(l C.TCOD_list_t, free bool) (result []string) {
	size := C.TCOD_list_size(l)

	result = make([]string, int(size))
	for i := 0; i < int(size); i++ {
		c := (*C.char)(C.TCOD_list_get(l, C.int(i)))
		result[i] = C.GoString(c)
		if free {
			C.free(unsafe.Pointer(c))
		}
	}
	if free {
		C.TCOD_list_delete(l)
	}
	return
}

func vectorShift(v *vector.Vector) (result interface{}) {
	result = v.At(0)
	v.Delete(0)
	return

}

func vectorRemove(v *vector.Vector, el interface{}) {
	for i := 0; i < v.Len(); i++ {
		if el == v.At(i) {
			v.Delete(i)
			break
		}
	}
}

//
//
// Key handling
//
type KeyCode C.TCOD_keycode_t

type Key struct {
	Vk      KeyCode
	C       byte
	Pressed bool
	LAlt    bool
	LCtrl   bool
	RAlt    bool
	RCtrl   bool
	Shift   bool
}

func toKey(k C.TCOD_key_t) (result Key) {
	result.Vk = KeyCode(k.vk)
	result.C = byte(k.c)
	result.Pressed = toBool(k.pressed)
	result.LAlt = toBool(k.lalt)
	result.LCtrl = toBool(k.lctrl)
	result.RAlt = toBool(k.ralt)
	result.RCtrl = toBool(k.rctrl)
	result.Shift = toBool(k.shift)
	return
}

func fromKey(k Key) (result C.TCOD_key_t) {
	result.vk = C.TCOD_keycode_t(k.Vk)
	result.c = C.char(k.C)
	result.pressed = fromBool(k.Pressed)
	result.lalt = fromBool(k.LAlt)
	result.lctrl = fromBool(k.LCtrl)
	result.ralt = fromBool(k.RAlt)
	result.rctrl = fromBool(k.RCtrl)
	result.shift = fromBool(k.Shift)
	return
}

//
//
// Bool handling
//
func toBool(b C.bool) bool {
	if int(b) == 1 {
		return true
	} else {
		return false
	}
	return false
}

func fromBool(b bool) C.bool {
	if b {
		return C.bool(1)
	} else {
		return C.bool(0)
	}
	return C.bool(0)
}

//
// Color handling
//

var COLOR_BLACK Color = Color{0, 0, 0}
var COLOR_DARKER_GREY Color = Color{31, 31, 31}
var COLOR_DARK_GREY Color = Color{63, 63, 63}
var COLOR_GREY Color = Color{128, 128, 128}
var COLOR_LIGHT_GREY Color = Color{191, 191, 191}
var COLOR_WHITE Color = Color{255, 255, 255}
var COLOR_RED Color = Color{255, 0, 0}
var COLOR_ORANGE Color = Color{255, 127, 0}
var COLOR_YELLOW Color = Color{255, 255, 0}
var COLOR_CHARTREUSE Color = Color{127, 255, 0}
var COLOR_GREEN Color = Color{0, 255, 0}
var COLOR_SEA Color = Color{0, 255, 127}
var COLOR_CYAN Color = Color{0, 255, 255}
var COLOR_SKY Color = Color{0, 127, 255}
var COLOR_BLUE Color = Color{0, 0, 255}
var COLOR_VIOLET Color = Color{127, 0, 255}
var COLOR_MAGENTA Color = Color{255, 0, 255}
var COLOR_PINK Color = Color{255, 0, 127}
var COLOR_DARK_RED Color = Color{127, 0, 0}
var COLOR_DARK_ORANGE Color = Color{127, 63, 0}
var COLOR_DARK_YELLOW Color = Color{127, 127, 0}
var COLOR_DARK_CHARTREUSE Color = Color{63, 127, 0}
var COLOR_DARK_GREEN Color = Color{0, 127, 0}
var COLOR_DARK_SEA Color = Color{0, 127, 63}
var COLOR_DARK_CYAN Color = Color{0, 127, 127}
var COLOR_DARK_SKY Color = Color{0, 63, 127}
var COLOR_DARK_BLUE Color = Color{0, 0, 127}
var COLOR_DARK_VIOLET Color = Color{63, 0, 127}
var COLOR_DARK_MAGENTA Color = Color{127, 0, 127}
var COLOR_DARK_PINK Color = Color{127, 0, 63}
var COLOR_DARKER_RED Color = Color{63, 0, 0}
var COLOR_DARKER_ORANGE Color = Color{63, 31, 0}
var COLOR_DARKER_YELLOW Color = Color{63, 63, 0}
var COLOR_DARKER_CHARTREUSE Color = Color{31, 63, 0}
var COLOR_DARKER_GREEN Color = Color{0, 63, 0}
var COLOR_DARKER_SEA Color = Color{0, 63, 31}
var COLOR_DARKER_CYAN Color = Color{0, 63, 63}
var COLOR_DARKER_SKY Color = Color{0, 31, 63}
var COLOR_DARKER_BLUE Color = Color{0, 0, 63}
var COLOR_DARKER_VIOLET Color = Color{31, 0, 63}
var COLOR_DARKER_MAGENTA Color = Color{63, 0, 63}
var COLOR_DARKER_PINK Color = Color{63, 0, 31}
var COLOR_LIGHT_RED Color = Color{255, 127, 127}
var COLOR_LIGHT_ORANGE Color = Color{255, 191, 127}
var COLOR_LIGHT_YELLOW Color = Color{255, 255, 127}
var COLOR_LIGHT_CHARTREUSE Color = Color{191, 255, 127}
var COLOR_LIGHT_GREEN Color = Color{127, 255, 127}
var COLOR_LIGHT_SEA Color = Color{127, 255, 191}
var COLOR_LIGHT_CYAN Color = Color{127, 255, 255}
var COLOR_LIGHT_SKY Color = Color{127, 191, 255}
var COLOR_LIGHT_BLUE Color = Color{127, 127, 255}
var COLOR_LIGHT_VIOLET Color = Color{191, 127, 255}
var COLOR_LIGHT_MAGENTA Color = Color{255, 127, 255}
var COLOR_LIGHT_PINK Color = Color{255, 127, 191}
var COLOR_DESATURATED_RED Color = Color{127, 63, 63}
var COLOR_DESATURATED_ORANGE Color = Color{127, 95, 63}
var COLOR_DESATURATED_YELLOW Color = Color{127, 127, 63}
var COLOR_DESATURATED_CHARTREUSE Color = Color{95, 127, 63}
var COLOR_DESATURATED_GREEN Color = Color{63, 127, 63}
var COLOR_DESATURATED_SEA Color = Color{63, 127, 95}
var COLOR_DESATURATED_CYAN Color = Color{63, 127, 127}
var COLOR_DESATURATED_SKY Color = Color{63, 95, 127}
var COLOR_DESATURATED_BLUE Color = Color{63, 63, 127}
var COLOR_DESATURATED_VIOLET Color = Color{95, 63, 127}
var COLOR_DESATURATED_MAGENTA Color = Color{127, 63, 127}
var COLOR_DESATURATED_PINK Color = Color{127, 63, 95}
var COLOR_SILVER Color = Color{203, 203, 203}
var COLOR_GOLD Color = Color{255, 255, 102}

type Color struct {
	R uint8
	G uint8
	B uint8
}

type ColCtrl C.TCOD_colctrl_t

func fromColor(c Color) (result C.TCOD_color_t) {
	result.r = C.uint8(c.R)
	result.g = C.uint8(c.G)
	result.b = C.uint8(c.B)
	return
}

func toColor(c C.TCOD_color_t) (result Color) {
	result.R = uint8(c.r)
	result.G = uint8(c.g)
	result.B = uint8(c.b)
	return
}

func NewColorRGB(r, g, b uint8) Color {
	return Color{R: r, G: g, B: b}
}

func NewColorHSV(h, s, v float32) Color {
	return toColor(C.TCOD_color_HSV(C.float(h), C.float(s), C.float(v)))
}

// basic operations 
func (self Color) Equals(c2 Color) bool {
	cc1 := fromColor(self)
	cc2 := fromColor(c2)
	return toBool(C._TCOD_color_equals((*C.TCOD_color_t)(&cc1), (*C.TCOD_color_t)(&cc2)))
}

func (self Color) Add(c2 Color) Color {
	cc1 := fromColor(self)
	cc2 := fromColor(c2)
	return toColor(C._TCOD_color_add((*C.TCOD_color_t)(&cc1), (*C.TCOD_color_t)(&cc2)))
}

func (self Color) Subtract(c2 Color) Color {
	cc1 := fromColor(self)
	cc2 := fromColor(c2)
	return toColor(C._TCOD_color_subtract((*C.TCOD_color_t)(&cc1), (*C.TCOD_color_t)(&cc2)))
}

func (self Color) Multiply(c2 Color) Color {
	cc1 := fromColor(self)
	cc2 := fromColor(c2)
	return toColor(C._TCOD_color_multiply((*C.TCOD_color_t)(&cc1), (*C.TCOD_color_t)(&cc2)))
}

func (self Color) MultiplyScalar(value float32) Color {
	c := fromColor(self)
	return toColor(C._TCOD_color_multiply_scalar((*C.TCOD_color_t)(&c), C.float(value)))
}

func (self Color) Lerp(c2 Color, coef float32) Color {
	cc1 := fromColor(self)
	cc2 := fromColor(c2)
	return toColor(C._TCOD_color_lerp((*C.TCOD_color_t)(&cc1), (*C.TCOD_color_t)(&cc2), C.float(coef)))
}

// HSV transformations

func (self Color) Lighten(ratio float32) Color {
	return self.Lerp(COLOR_WHITE, ratio)
}

func (self Color) Darken(ratio float32) Color {
	return self.Lerp(COLOR_BLACK, ratio)
}

func (self Color) SetHSV(h float32, s float32, v float32) Color {
	c := C.TCOD_color_t{}
	C.TCOD_color_set_HSV(&c, C.float(h), C.float(s), C.float(v))
	return toColor(c)
}

func (self Color) GetHue() float32 {
	return float32(C.TCOD_color_get_hue(fromColor(self)))
}

func (self Color) SetHue(h float32) Color {
	c := C.TCOD_color_t{}
	C.TCOD_color_set_hue(&c, C.float(h))
	return toColor(c)
}

func (self Color) GetSaturation() float32 {
	return float32(C.TCOD_color_get_saturation(fromColor(self)))
}

func (self Color) SetSaturation(h float32) Color {
	c := C.TCOD_color_t{}
	C.TCOD_color_set_saturation(&c, C.float(h))
	return toColor(c)
}

func (self Color) GetValue() float32 {
	return float32(C.TCOD_color_get_value(fromColor(self)))
}

func (self Color) SetValue(h float32) Color {
	c := C.TCOD_color_t{}
	C.TCOD_color_set_value(&c, C.float(h))
	return toColor(c)
}

func (self Color) ShiftHue(hshift float32) Color {
	c := C.TCOD_color_t{}
	C.TCOD_color_shift_hue(&c, C.float(hshift))
	return toColor(c)
}

func (self Color) ScaleHSV(scoef, vcoef float32) Color {
	c := C.TCOD_color_t{}
	C.TCOD_color_scale_HSV(&c, C.float(scoef), C.float(vcoef))
	return toColor(c)
}

func (self Color) GetHSV() (h, s, v float32) {
	var ch, cs, sv C.float
	C.TCOD_color_get_HSV(fromColor(self), &ch, &cs, &sv)
	h = float32(ch)
	s = float32(cs)
	v = float32(sv)
	return
}

func ColorGenMap(cmap []Color, nbKey int, keyColor []Color, keyIndex []int) {
	for segment := 0; segment < nbKey-1; segment++ {
		idxStart := keyIndex[segment]
		idxEnd := keyIndex[segment+1]
		for idx := idxStart; idx <= idxEnd; idx++ {
			cmap[idx] = keyColor[segment].Lerp(keyColor[segment+1], float32(idx-idxStart)/float32(idxEnd-idxStart))
		}
	}
}

//
// Mouse
//
//
type Mouse struct {
	X, Y           int
	Dx, Dy         int
	Cx, Cy         int
	Dcx, Dcy       int
	LButton        bool
	RButton        bool
	MButton        bool
	LButtonPressed bool
	RButtonPressed bool
	MButtonPressed bool
	WheelUp        bool
	WheelDown      bool
}

func fromMouse(m Mouse) (result C.TCOD_mouse_t) {
	result.x = C.int(m.X)
	result.y = C.int(m.Y)
	result.dx = C.int(m.Dx)
	result.dy = C.int(m.Dy)
	result.cx = C.int(m.Cx)
	result.cy = C.int(m.Cy)
	result.dcx = C.int(m.Dcx)
	result.dcy = C.int(m.Dcy)
	result.lbutton = fromBool(m.LButton)
	result.rbutton = fromBool(m.RButton)
	result.mbutton = fromBool(m.MButton)
	result.lbutton_pressed = fromBool(m.LButtonPressed)
	result.rbutton_pressed = fromBool(m.RButtonPressed)
	result.mbutton_pressed = fromBool(m.MButtonPressed)
	result.wheel_up = fromBool(m.WheelUp)
	result.wheel_down = fromBool(m.WheelDown)
	return
}

func toMouse(m C.TCOD_mouse_t) (result Mouse) {
	result.X = int(m.x)
	result.Y = int(m.y)
	result.Dx = int(m.dx)
	result.Dy = int(m.dy)
	result.Cx = int(m.cx)
	result.Cy = int(m.cy)
	result.Dcx = int(m.dcx)
	result.Dcy = int(m.dcy)
	result.LButton = toBool(m.lbutton)
	result.RButton = toBool(m.rbutton)
	result.MButton = toBool(m.mbutton)
	result.LButtonPressed = toBool(m.lbutton_pressed)
	result.RButtonPressed = toBool(m.rbutton_pressed)
	result.MButtonPressed = toBool(m.mbutton_pressed)
	result.WheelUp = toBool(m.wheel_up)
	result.WheelDown = toBool(m.wheel_down)
	return
}

func MouseGetStatus() Mouse {
	return toMouse(C.TCOD_mouse_get_status())
}

func MouseShowCursor(visible bool) {
	C.TCOD_mouse_show_cursor(fromBool(visible))
}

func MouseIsCursorVisible() bool {
	return toBool(C.TCOD_mouse_is_cursor_visible())
}

func MouseMove(x, y int) {
	C.TCOD_mouse_move(C.int(x), C.int(y))
}

//
//
// Console
//
//
type Alignment C.TCOD_alignment_t
type Renderer C.TCOD_renderer_t

type IConsole interface {
	GetData() C.TCOD_console_t
	GetDefaultBackground() Color
	GetDefaultForeground() Color
	SetDefaultForeground(color Color)
	SetDefaultBackground(color Color)
	Clear()
	GetCharBackground(x, y int) Color
	GetCharForeground(x, y int) Color
	SetCharBackground(x, y int, color Color, flag BkgndFlag)
	SetCharForeground(x, y int, color Color)
	SetChar(x, y int, c int)
	PutChar(x, y, c int, flag BkgndFlag)
	PutCharEx(x, y, c int, fore, back Color)
	Print(x, y int, fmts string, v ...interface{})
	PrintEx(x, y int, flag BkgndFlag, alignment Alignment, fmts string, v ...interface{})
	PrintRect(x, y, w, h int, fmts string, v ...interface{}) int
	PrintRectEx(x, y, w, h int, flag BkgndFlag, alignment Alignment, fmts string, v ...interface{}) int
	HeightRect(x, y, w, h int, fmts string, v ...interface{}) int
	SetBackgroundFlag(flag BkgndFlag)
	GetBackgroundFlag() BkgndFlag
	SetAlignment(alignment Alignment)
	GetAlignment() Alignment
	Rect(x, y, w, h int, clear bool, flag BkgndFlag)
	Hline(x, y, l int, flag BkgndFlag)
	Vline(x, y, l int, flag BkgndFlag)
	PrintFrame(x, y, w, h int, empty bool, flag BkgndFlag, fmts string, v ...interface{})
	GetChar(x, y int) int
	GetWidth() int
	GetHeight() int
	SetKeyColor(color Color)
	Blit(xSrc, ySrc, wSrc, hSrc int, dst IConsole, xDst, yDst int, foregroundAlpha, backgroundAlpha float32)
}

// Console 

type Console struct {
	Data C.TCOD_console_t
}

func deleteConsole(c *Console) {
	C.TCOD_console_delete(c.Data)
}

func NewConsole(w, h int) *Console {
	result := &Console{C.TCOD_console_new(C.int(w), C.int(h))}
	runtime.SetFinalizer(result, deleteConsole)
	return result
}

func (self *Console) GetData() C.TCOD_console_t {
	return self.Data
}

func (self *Console) SetDefaultBackground(color Color) {
	C.TCOD_console_set_default_background(self.Data, fromColor(color))
}

func (self *Console) SetDefaultForeground(color Color) {
	C.TCOD_console_set_default_foreground(self.Data, fromColor(color))
}

func (self *Console) Clear() {
	C.TCOD_console_clear(self.Data)
}

func (self *Console) SetCharBackground(x, y int, color Color, flag BkgndFlag) {
	ccolor := fromColor(color)
	C._TCOD_console_set_char_background(self.Data, C.int(x), C.int(y), (*C.TCOD_color_t)(&ccolor), C.TCOD_bkgnd_flag_t(flag))
}

func (self *Console) SetCharForeground(x, y int, color Color) {
	ccolor := fromColor(color)
	C._TCOD_console_set_char_foreground(self.Data, C.int(x), C.int(y), (*C.TCOD_color_t)(&ccolor))
}

func (self *Console) SetChar(x, y int, c int) {
	C.TCOD_console_set_char(self.Data, C.int(x), C.int(y), C.int(c))
}

func (self *Console) PutChar(x, y, c int, flag BkgndFlag) {
	C.TCOD_console_put_char(self.Data, C.int(x), C.int(y), C.int(c), C.TCOD_bkgnd_flag_t(flag))
}

func (self *Console) PutCharEx(x, y, c int, fore, back Color) {
	forec := fromColor(fore)
	backc := fromColor(back)
	C._TCOD_console_put_char_ex(self.Data, C.int(x), C.int(y), C.int(c),
		(*C.TCOD_color_t)(&forec), (*C.TCOD_color_t)(&backc))
}

func (self *Console) Print(x, y int, fmts string, v ...interface{}) {
	s := fmt.Sprintf(fmts, v...)
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C._TCOD_console_print(self.Data, C.int(x), C.int(y), cs)
}

func (self *Console) PrintEx(x, y int, flag BkgndFlag, alignment Alignment, fmts string, v ...interface{}) {
	s := fmt.Sprintf(fmts, v...)
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C._TCOD_console_print_ex(self.Data, C.int(x), C.int(y), C.TCOD_bkgnd_flag_t(flag), C.TCOD_alignment_t(alignment), cs)
}

func (self *Console) PrintRect(x, y, w, h int, fmts string, v ...interface{}) int {
	s := fmt.Sprintf(fmts, v...)
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return int(C._TCOD_console_print_rect(self.Data, C.int(x), C.int(y), C.int(w), C.int(h), cs))
}

func (self *Console) PrintRectEx(x, y, w, h int, flag BkgndFlag, alignment Alignment, fmts string, v ...interface{}) int {
	s := fmt.Sprintf(fmts, v...)
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return int(C._TCOD_console_print_rect_ex(self.Data, C.int(x), C.int(y), C.int(w), C.int(h), C.TCOD_bkgnd_flag_t(flag),
		C.TCOD_alignment_t(alignment), cs))
}

func (self *Console) SetBackgroundFlag(flag BkgndFlag) {
	C.TCOD_console_set_background_flag(self.Data, C.TCOD_bkgnd_flag_t(flag))
}

func (self *Console) GetBackgroundFlag() BkgndFlag {
	return BkgndFlag(C.TCOD_console_get_background_flag(self.Data))
}

func (self *Console) SetAlignment(alignment Alignment) {
	C.TCOD_console_set_alignment(self.Data, C.TCOD_alignment_t(alignment))
}

func (self *Console) GetAlignment() Alignment {
	return Alignment(C.TCOD_console_get_alignment(self.Data))
}

func (self *Console) HeightRect(x, y, w, h int, fmts string, v ...interface{}) int {
	s := fmt.Sprintf(fmts, v...)
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return int(C._TCOD_console_height_rect(self.Data, C.int(x), C.int(y), C.int(w), C.int(h), cs))
}

func (self *Console) Rect(x, y, w, h int, clear bool, flag BkgndFlag) {
	C.TCOD_console_rect(self.Data, C.int(x), C.int(y), C.int(w), C.int(h), fromBool(clear), C.TCOD_bkgnd_flag_t(flag))
}

func (self *Console) Hline(x, y, l int, flag BkgndFlag) {
	C.TCOD_console_hline(self.Data, C.int(x), C.int(y), C.int(l), C.TCOD_bkgnd_flag_t(flag))
}

func (self *Console) Vline(x, y, l int, flag BkgndFlag) {
	C.TCOD_console_hline(self.Data, C.int(x), C.int(y), C.int(l), C.TCOD_bkgnd_flag_t(flag))
}

func (self *Console) PrintFrame(x, y, w, h int, empty bool, flag BkgndFlag, fmts string, v ...interface{}) {
	s := fmt.Sprintf(fmts, v...)
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C._TCOD_console_print_frame(self.Data, C.int(x), C.int(y), C.int(w), C.int(h),
		fromBool(empty), C.TCOD_bkgnd_flag_t(flag), cs)

}

// TODO check unicode support
//TCODLIB_API void TCOD_console_map_string_to_font_utf(const wchar_t *s, int fontCharX, int fontCharY);
//TCODLIB_API void TCOD_console_print_left_utf(TCOD_console_t con,int x, int y, TCOD_bkgnd_flag_t flag, const wchar_t *fmt, ...);
//TCODLIB_API void TCOD_console_print_right_utf(TCOD_console_t con,int x, int y, TCOD_bkgnd_flag_t flag, const wchar_t *fmt, ...);
//TCODLIB_API void TCOD_console_print_center_utf(TCOD_console_t con,int x, int y, TCOD_bkgnd_flag_t flag, const wchar_t *fmt, ...);
//TCODLIB_API int TCOD_console_print_left_rect_utf(TCOD_console_t con,int x, int y, int w, int h, TCOD_bkgnd_flag_t flag, const wchar_t *fmt, ...);
//TCODLIB_API int TCOD_console_print_right_rect_utf(TCOD_console_t con,int x, int y, int w, int h, TCOD_bkgnd_flag_t flag, const wchar_t *fmt, ...);
//TCODLIB_API int TCOD_console_print_center_rect_utf(TCOD_console_t con,int x, int y, int w, int h, TCOD_bkgnd_flag_t flag, const wchar_t *fmt, ...);
//TCODLIB_API int TCOD_console_height_left_rect_utf(TCOD_console_t con,int x, int y, int w, int h, const wchar_t *fmt, ...);
//TCODLIB_API int TCOD_console_height_right_rect_utf(TCOD_console_t con,int x, int y, int w, int h, const wchar_t *fmt, ...);
//TCODLIB_API int TCOD_console_height_center_rect_utf(TCOD_console_t con,int x, int y, int w, int h, const wchar_t *fmt, ...);
//#endif

func (self *Console) GetDefaultBackground() Color {
	return toColor(C.TCOD_console_get_default_background(self.Data))
}

func (self *Console) GetDefaultForeground() Color {
	return toColor(C.TCOD_console_get_default_foreground(self.Data))
}

func (self *Console) GetCharBackground(x, y int) Color {
	return toColor(C.TCOD_console_get_char_background(self.Data, C.int(x), C.int(y)))
}

func (self *Console) GetCharForeground(x, y int) Color {
	return toColor(C.TCOD_console_get_char_foreground(self.Data, C.int(x), C.int(y)))
}

func (self *Console) GetChar(x, y int) int {
	return int(C.TCOD_console_get_char(self.Data, C.int(x), C.int(y)))
}

func (self *Console) GetWidth() int {
	return int(C.TCOD_console_get_width(self.Data))
}

func (self *Console) GetHeight() int {
	return int(C.TCOD_console_get_height(self.Data))
}

func (self *Console) SetKeyColor(color Color) {
	ccolor := fromColor(color)
	C._TCOD_console_set_key_color(self.Data, (*C.TCOD_color_t)(&ccolor))
}

func (self *Console) Blit(xSrc, ySrc, wSrc, hSrc int, dst IConsole, xDst, yDst int, foregroundAlpha, backgroundAlpha float32) {
	C.TCOD_console_blit(self.Data, C.int(xSrc), C.int(ySrc), C.int(wSrc), C.int(hSrc),
		dst.GetData(), C.int(xDst), C.int(yDst), C.float(foregroundAlpha), C.float(backgroundAlpha))
}

// RootConsole

type RootConsole struct {
	Console
}

func NewRootConsole(w, h int, title string, fullscreen bool) *RootConsole {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))
	C.TCOD_console_init_root(C.int(w), C.int(h), ctitle, fromBool(fullscreen), C.TCOD_renderer_t(RENDERER_SDL))
	// in root console, Data field is nil
	return &RootConsole{}
}

func NewRootConsoleWithFont(w, h int, title string, fullscreen bool, fontFile string, fontFlags, nbCharHoriz,
	nbCharVertic int, renderer Renderer) *RootConsole {
	cfontFile := C.CString(fontFile)
	defer C.free(unsafe.Pointer(cfontFile))
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))
	C.TCOD_console_set_custom_font(cfontFile, C.int(fontFlags), C.int(nbCharHoriz), C.int(nbCharVertic))
	C.TCOD_console_init_root(C.int(w), C.int(h), ctitle, fromBool(fullscreen), C.TCOD_renderer_t(renderer))
	// in root console, Data field is nil
	return &RootConsole{}
}

func (self *RootConsole) SetWindowTitle(title string) {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))
	C.TCOD_console_set_window_title(ctitle)

}

func (self *RootConsole) SetFullscreen(fullscreen bool) {
	C.TCOD_console_set_fullscreen(fromBool(fullscreen))
}

func (self *RootConsole) IsFullscreen() bool {
	return toBool(C.TCOD_console_is_fullscreen())
}

func (self *RootConsole) IsWindowClosed() bool {
	return toBool(C.TCOD_console_is_window_closed())
}

func (self *RootConsole) SetCustomFont(fontFile string, flags int, nbCharHoriz int, nbCharVertic int) {
	cfontFile := C.CString(fontFile)
	defer C.free(unsafe.Pointer(cfontFile))
	C.TCOD_console_set_custom_font(cfontFile, C.int(flags), C.int(nbCharHoriz), C.int(nbCharVertic))
}

func (self *RootConsole) MapAsciiCodeToFont(asciiCode, fontCharX, fontCharY int) {
	C.TCOD_console_map_ascii_code_to_font(C.int(asciiCode), C.int(fontCharX), C.int(fontCharY))
}

func (self *RootConsole) MapAsciiCodesToFont(asciiCode, fontCharX, fontCharY int) {
	C.TCOD_console_map_ascii_code_to_font(C.int(asciiCode), C.int(fontCharX), C.int(fontCharY))
}

func (self *RootConsole) MapStringToFont(s string, fontCharX, fontCharY int) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.TCOD_console_map_string_to_font(cs, C.int(fontCharX), C.int(fontCharY))
}

func (self *RootConsole) SetDirty(x, y, w, h int) {
	C.TCOD_console_set_dirty(C.int(x), C.int(y), C.int(w), C.int(h))
}

func (self *RootConsole) SetFade(val uint8, fade Color) {
	ccolor := fromColor(fade)
	C._TCOD_console_set_fade(C.uint8(val), (*C.TCOD_color_t)(&ccolor))
}

func (self *RootConsole) GetFade() uint8 {
	return uint8(C.TCOD_console_get_fade())
}

func (self *RootConsole) GetFadingColor() Color {
	return toColor(C.TCOD_console_get_fading_color())
}

func (self *RootConsole) Flush() {
	C.TCOD_console_flush()
}

func (self *RootConsole) SetColorControl(ctrl ColCtrl, fore, back Color) {
	forec := fromColor(fore)
	backc := fromColor(back)
	C._TCOD_console_set_color_control(C.TCOD_colctrl_t(ctrl),
		(*C.TCOD_color_t)(&forec), (*C.TCOD_color_t)(&backc))
}

func (self *RootConsole) CheckForKeypress(flags int) Key {
	return toKey(C.TCOD_console_check_for_keypress(C.int(flags)))
}

func (self *RootConsole) WaitForKeypress(flush bool) Key {
	return toKey(C.TCOD_console_wait_for_keypress(fromBool(flush)))
}

func (self *RootConsole) SetKeyboardRepeat(initialDelay, interval int) {
	C.TCOD_console_set_keyboard_repeat(C.int(initialDelay), C.int(interval))
}

func (self *RootConsole) DisableKeyboardRepeat() {
	C.TCOD_console_disable_keyboard_repeat()
}

func (self *RootConsole) IsKeyPressed(keyCode KeyCode) bool {
	return toBool(C.TCOD_console_is_key_pressed(C.TCOD_keycode_t(keyCode)))
}

func (self *RootConsole) Credits() {
	C.TCOD_console_credits()
}

func (self *RootConsole) ResetCredits() {
	C.TCOD_console_credits_reset()
}

func (self *RootConsole) RenderCredits(x, y int, alpha bool) bool {
	return toBool(C.TCOD_console_credits_render(C.int(x), C.int(y), fromBool(alpha)))
}

//
// Bresenham line algorithm
// Fully ported to Go for easier callbacks
//
//

type LineListener func(x, y int, userData interface{}) bool

type Point struct {
	x, y int
}

// thread-safe versions
type BresenhamData struct {
	stepx  int
	stepy  int
	e      int
	deltax int
	deltay int
	origx  int
	origy  int
	destx  int
	desty  int
}

var bresenhamData BresenhamData

func lineInitMt(xFrom, yFrom, xTo, yTo int, data *BresenhamData) {
	data.origx = xFrom
	data.origy = yFrom
	data.destx = xTo
	data.desty = yTo
	data.deltax = xTo - xFrom
	data.deltay = yTo - yFrom
	if data.deltax > 0 {
		data.stepx = 1
	} else if data.deltax < 0 {
		data.stepx = -1
	} else {
		data.stepx = 0
	}
	if data.deltay > 0 {
		data.stepy = 1
	} else if data.deltay < 0 {
		data.stepy = -1
	} else {
		data.stepy = 0
	}
	if data.stepx*data.deltax > data.stepy*data.deltay {
		data.e = data.stepx * data.deltax
		data.deltax *= 2
		data.deltay *= 2
	} else {
		data.e = data.stepy * data.deltay
		data.deltax *= 2
		data.deltay *= 2
	}
}

func lineStepMt(xCur, yCur *int, data *BresenhamData) bool {
	if data.stepx*data.deltax > data.stepy*data.deltay {
		if data.origx == data.destx {
			return true
		}
		data.origx += data.stepx
		data.e -= data.stepy * data.deltay
		if data.e < 0 {
			data.origy += data.stepy
			data.e += data.stepx * data.deltax
		}
	} else {
		if data.origy == data.desty {
			return true
		}
		data.origy += data.stepy
		data.e -= data.stepx * data.deltax
		if data.e < 0 {
			data.origx += data.stepx
			data.e += data.stepy * data.deltay
		}
	}
	*xCur = data.origx
	*yCur = data.origy
	return false
}

func lineInit(xFrom, yFrom, xTo, yTo int) {
	lineInitMt(xFrom, yFrom, xTo, yTo, &bresenhamData)
}

func lineStep(xCur, yCur *int) bool {
	return lineStepMt(xCur, yCur, &bresenhamData)
}

func LineMt(xo, yo, xd, yd int, listener LineListener, userData interface{}, data *BresenhamData) bool {
	lineInitMt(xo, yo, xd, yd, data)
	if !listener(xo, yo, userData) {
		return false
	}
	for !lineStepMt(&xo, &yo, data) {
		if !listener(xo, yo, userData) {
			return false
		}
	}
	return true
}

func Line(xo, yo, xd, yd int, userData interface{}, listener LineListener) bool {
	return LineMt(xo, yo, xd, yd, listener, userData, &bresenhamData)
}

// returns vector with Points where the line was drawn
func LinePoints(xo, yo, xd, yd int) vector.Vector {
	result := vector.Vector{}
	Line(xo, yo, xd, yd, nil, func(x, y int, data interface{}) bool {
		result.Push(Point{x, y})
		return true
	})
	return result
}

//
//
// Name generator
//

//
func NamegenParse(filename string, random *Random) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	C.TCOD_namegen_parse(cfilename, random.Data)
}

// generate a name
func NamegenGenerate(name string) string {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	c := C.TCOD_namegen_generate(cname, fromBool(true))
	defer C.free(unsafe.Pointer(c))
	return C.GoString(c)
}

// generate a name using a custom generation rule
func NamegenGenerateCustom(name, rule string) string {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	crule := C.CString(rule)
	defer C.free(unsafe.Pointer(crule))
	c := C.TCOD_namegen_generate_custom(cname, crule, fromBool(true))
	defer C.free(unsafe.Pointer(c))
	return C.GoString(c)
}

// retrieve the list of all available syllable set names
func NamegenGetSets() []string {
	return toStringSlice(C.TCOD_namegen_get_sets(), false)
}

// delete a generator
func NamegenDestroy() {
	C.TCOD_namegen_destroy()
}

//
//
// Text field
// TODO this is available only in debug version? (1.5.0)
//
//
type Text struct {
	Data C.TCOD_text_t
}

func deleteText(t *Text) {
	C.TCOD_text_delete(t.Data)
}

func NewText(x, y, w, h, maxChars int) *Text {
	result := &Text{C.TCOD_text_init(C.int(x), C.int(y), C.int(w), C.int(h), C.int(maxChars))}
	runtime.SetFinalizer(result, deleteText)
	return result
}

func (self *Text) SetProperties(cursorChar int, blinkInterval int, prompt string, tabSize int) {
	cprompt := C.CString(prompt)
	defer C.free(unsafe.Pointer(cprompt))
	C.TCOD_text_set_properties(self.Data, C.int(cursorChar), C.int(blinkInterval), cprompt, C.int(tabSize))
}

func (self *Text) SetColors(fore, back Color, backTransparency float32) {
	forec := fromColor(fore)
	backc := fromColor(back)
	C._TCOD_text_set_colors(self.Data,
		(*C.TCOD_color_t)(&forec),
		(*C.TCOD_color_t)(&backc),
		C.float(backTransparency))
}

func (self *Text) Update(key Key) {
	C.TCOD_text_update(self.Data, fromKey(key))
}

func (self *Text) Render(console IConsole) {
	C.TCOD_text_render(self.Data, console.GetData())
}

func (self *Console) RenderText(text *Text) {
	C.TCOD_text_render(text.Data, self.Data)
}

func (self *Text) Get() string {
	t := C.TCOD_text_get(self.Data)
	return C.GoString(t)

}

func (self *Text) Reset() {
	C.TCOD_text_reset(self.Data)
}

func SysElapsedMilliseconds() uint32 {
	return uint32(C.TCOD_sys_elapsed_milli())
}

func SysElapsedSeconds() float32 {
	return float32(C.TCOD_sys_elapsed_seconds())
}

func SysSleepMilliseconds(val uint32) {
	C.TCOD_sys_sleep_milli(C.uint32(val))
}

func SysSaveScreenshot() {
	C.TCOD_sys_save_screenshot(nil)
}

func SysSaveScreenshotToFile(filename string) {
	if filename == "" {
		C.TCOD_sys_save_screenshot(nil)
	} else {
		cfilename := C.CString(filename)
		defer C.free(unsafe.Pointer(cfilename))
		C.TCOD_sys_save_screenshot(cfilename)
	}
}

func SysForceFullscreenResolution(width, height int) {
	C.TCOD_sys_force_fullscreen_resolution(C.int(width), C.int(height))
}

func SysSetFps(val int) {
	C.TCOD_sys_set_fps(C.int(val))
}

func SysGetFps() int {
	return int(C.TCOD_sys_get_fps())
}

func SysGetLastFrameLength() float32 {
	return float32(C.TCOD_sys_get_last_frame_length())
}

func SysGetCurrentResolution() (w, h int) {
	var cw, ch C.int
	C.TCOD_sys_get_current_resolution(&cw, &ch)
	w, h = int(cw), int(ch)
	return
}

func SysGetFullscreenOffsets() (offx, offy int) {
	var coffx, coffy C.int
	C.TCOD_sys_get_fullscreen_offsets(&coffx, &coffy)
	offx, offy = int(coffx), int(coffy)
	return
}

func SysUpdateChar(asciiCode, fontx, fonty int, img Image, x, y int) {
	C.TCOD_sys_update_char(C.int(asciiCode), C.int(fontx), C.int(fonty), img.Data, C.int(x), C.int(y))
}

func SysGetCharSize() (w, h int) {
	var cw, ch C.int
	C.TCOD_sys_get_char_size(&cw, &ch)
	w, h = int(cw), int(ch)
	return
}

// filesystem stuff
func SysCreateDirectory(path string) bool {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return toBool(C.TCOD_sys_create_directory(cpath))
}

func SysDeleteFile(path string) bool {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return toBool(C.TCOD_sys_delete_file(cpath))
}

func SysDeleteDirectory(path string) bool {
	return toBool(C.TCOD_sys_delete_directory(C.CString(path)))
}

func SysIsDirectory(path string) bool {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return toBool(C.TCOD_sys_is_directory(cpath))
}

func SysGetDirectoryContent(path, pattern string) []string {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	cpattern := C.CString(pattern)
	defer C.free(unsafe.Pointer(cpattern))
	return toStringSlice(
		C.TCOD_sys_get_directory_content(
			cpath, cpattern),
		true)
}

func SysFileExists(filename string) bool {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	return toBool(C._TCOD_sys_file_exists(cfilename))
}

func SysGetNumCores() int {
	return int(C.TCOD_sys_get_num_cores())
}

// Clipboard 

func SysClipboardSet(value string) {
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	C.TCOD_sys_clipboard_set(cvalue)
}

func SysClipboardGet() string {
	return C.GoString(C.TCOD_sys_clipboard_get())
}

//
// Field Of View Map
//

type Map struct {
	Data C.TCOD_map_t
}

type FovAlgorithm C.TCOD_fov_algorithm_t

// destroy a map
func deleteMap(m *Map) {
	C.TCOD_map_delete(m.Data)
}

func NewMap(width, height int) *Map {
	result := &Map{C.TCOD_map_new(C.int(width), C.int(height))}
	runtime.SetFinalizer(result, deleteMap)
	return result
}

// set all cells as solid rock (cannot see through nor walk)
func (self *Map) Clear(isTransparent bool, isWalkable bool) {
	C.TCOD_map_clear(self.Data, fromBool(isTransparent), fromBool(isWalkable))
}

// copy a map to another, reallocating it when needed
func (self *Map) Copy(dest Map) {
	C.TCOD_map_copy(self.Data, dest.Data)
}

// change a cell properties
func (self *Map) SetProperties(x, y int, isTransparent bool, isWalkable bool) {
	C.TCOD_map_set_properties(self.Data, C.int(x), C.int(y), fromBool(isTransparent), fromBool(isWalkable))
}

// calculate the field of view (potentially visible cells from player_x,player_y)
func (self *Map) ComputeFov(playerX, playerY, maxRadius int, lightWalls bool, algo FovAlgorithm) {
	C.TCOD_map_compute_fov(self.Data, C.int(playerX), C.int(playerY),
		C.int(maxRadius), fromBool(lightWalls),
		C.TCOD_fov_algorithm_t(algo))
}

// check if a cell is in the last computed field of view
func (self *Map) IsInFov(x, y int) bool {
	return toBool(C.TCOD_map_is_in_fov(self.Data, C.int(x), C.int(y)))
}

func (self *Map) SetInFov(x, y int, fov bool) {
	C.TCOD_map_set_in_fov(self.Data, C.int(x), C.int(y), fromBool(fov))
}

// retrieve properties from the map

func (self *Map) IsTransparent(x, y int) bool {
	return toBool(C.TCOD_map_is_transparent(self.Data, C.int(x), C.int(y)))
}

func (self *Map) IsWalkable(x, y int) bool {
	return toBool(C.TCOD_map_is_walkable(self.Data, C.int(x), C.int(y)))
}

func (self *Map) GetWidth() int {
	return int(C.TCOD_map_get_width(self.Data))
}

func (self *Map) GetHeight() int {
	return int(C.TCOD_map_get_height(self.Data))
}

func (self *Map) GetNbCells() int {
	return int(C.TCOD_map_get_nb_cells(self.Data))
}

//
// BSP Dungeon generation
//
//
type Bsp struct {
	X, Y, W, H         int   // node position & size
	Position           int   // position of splitting
	Level              uint8 // level in the tree
	Horizontal         bool  // horizontal splitting ?
	next, father, sons *Bsp  // BSP tree hierarchy structuring
}

type BspListener func(node *Bsp, userData interface{}) bool

func (self *Bsp) AddSon(son *Bsp) {
	lastson := self.sons
	son.father = self
	for lastson != nil && lastson.next != nil {
		lastson = lastson.next
	}
	if lastson != nil {
		lastson.next = son
	} else {
		self.sons = son
	}
}

func NewBspWithSize(x, y, w, h int) (result *Bsp) {
	result = new(Bsp)
	*result = Bsp{X: x, Y: y, W: w, H: h}
	return
}

func (self *Bsp) Left() *Bsp {
	return self.sons
}

func (self *Bsp) Right() *Bsp {
	if self.sons != nil {
		return self.sons.next
	} else {
		return nil
	}
	return nil
}

func (self *Bsp) Father() *Bsp {
	return self.father
}

func (self *Bsp) IsLeaf() bool {
	return self.sons == nil
}

func NewBspIntern(father *Bsp, left bool) *Bsp {
	bsp := new(Bsp)
	if father.Horizontal {
		bsp.X = father.X
		bsp.W = father.W
		if left {
			bsp.Y = father.Y
		} else {
			bsp.Y = father.Position
		}
		if left {
			bsp.H = father.Position - bsp.Y
		} else {
			bsp.H = father.Y + father.H - father.Position
		}
	} else {
		bsp.Y = father.Y
		bsp.H = father.H
		if left {
			bsp.X = father.X
		} else {
			bsp.X = father.Position
		}
		if left {
			bsp.W = father.Position - bsp.X
		} else {
			bsp.W = father.X + father.W - father.Position
		}
	}
	bsp.Level = father.Level + 1
	return bsp
}

func (self *Bsp) TraversePreOrder(listener BspListener, userData interface{}) bool {
	if !listener(self, userData) {
		return false
	}
	if self.Left() != nil && !self.Left().TraversePreOrder(listener, userData) {
		return false
	}
	if self.Right() != nil && !self.Right().TraversePreOrder(listener, userData) {
		return false
	}
	return true
}

func (self *Bsp) TraverseInOrder(listener BspListener, userData interface{}) bool {
	if self.Left() != nil && !self.Left().TraverseInOrder(listener, userData) {
		return false
	}
	if !listener(self, userData) {
		return false
	}
	if self.Right() != nil && !self.Right().TraverseInOrder(listener, userData) {
		return false
	}
	return true
}

func (self *Bsp) TraversePostOrder(listener BspListener, userData interface{}) bool {
	if self.Left() != nil && !self.Left().TraversePostOrder(listener, userData) {
		return false
	}
	if self.Right() != nil && !self.Right().TraversePostOrder(listener, userData) {
		return false
	}
	if !listener(self, userData) {
		return false
	}
	return true
}

func (self *Bsp) TraverseLevelOrder(listener BspListener, userData interface{}) bool {
	stack := vector.Vector{}
	stack.Push(self)
	for len(stack) > 0 {
		node := vectorShift(&stack).(*Bsp)
		if node.Left() != nil {
			stack.Push(node.Left())
		}
		if node.Right() != nil {
			stack.Push(node.Right())
		}
		if !listener(node, userData) {
			return false
		}
	}
	return true
}

// TODO can it store Go values in list structure??
// maybe replace it with record
func (self *Bsp) TraverseInvertedLevelOrder(listener BspListener, userData interface{}) bool {

	stack1 := vector.Vector{}
	stack2 := vector.Vector{}
	stack1.Push(self)
	for stack1.Len() > 0 {
		node := vectorShift(&stack1).(*Bsp)
		stack2.Push(node)
		if node.Left() != nil {
			stack1.Push(node.Left())
		}
		if node.Right() != nil {
			stack1.Push(node.Right())
		}
	}
	for stack2.Len() > 0 {
		node := stack2.Pop().(*Bsp)
		if !listener(node, userData) {
			return false
		}
	}
	return true
}

func (self *Bsp) RemoveSons() {
	node := self.sons
	var nextNode *Bsp
	for node != nil {
		nextNode = node.next
		node.RemoveSons()
		node = nextNode
	}
	self.sons = nil
}

func (self *Bsp) SplitOnce(horizontal bool, position int) {
	self.Horizontal = horizontal
	self.Position = position
	self.AddSon(NewBspIntern(self, true))
	self.AddSon(NewBspIntern(self, false))
}

func (self *Bsp) SplitRecursive(randomizer *Random, nb int, minHSize int, minVSize int, maxHRatio float32, maxVRatio float32) {
	var horiz bool
	var position int
	if nb == 0 || (self.W < 2*minHSize && self.H < 2*minVSize) {
		return
	}
	// promote square rooms
	if self.H < 2*minVSize || float32(self.W) > float32(self.H)*maxHRatio {
		horiz = false
	} else if self.W < 2*minHSize || float32(self.H) > float32(self.W)*maxVRatio {
		horiz = true
	} else {
		horiz = (randomizer.GetInt(0, 1) == 0)
	}
	if horiz {
		position = randomizer.GetInt(self.Y+minVSize, self.Y+self.H-minVSize)
	} else {
		position = randomizer.GetInt(self.X+minHSize, self.X+self.W-minHSize)
	}
	self.SplitOnce(horiz, position)
	if self.Left() != nil {
		self.Left().SplitRecursive(randomizer, nb-1, minHSize, minVSize, maxHRatio, maxVRatio)
	}
	if self.Right() != nil {
		self.Right().SplitRecursive(randomizer, nb-1, minHSize, minVSize, maxHRatio, maxVRatio)
	}
}

func (self *Bsp) Resize(x, y, w, h int) {
	self.X, self.Y, self.W, self.H = x, y, w, h
	if self.Left() != nil {
		if self.Horizontal {
			self.Left().Resize(x, y, w, self.Position-y)
			if self.Right() != nil {
				self.Right().Resize(x, self.Position, w, y+h-self.Position)
			}
		} else {
			self.Left().Resize(x, y, self.Position-x, h)
			if self.Right() != nil {
				self.Right().Resize(self.Position, y, x+w-self.Position, h)
			}
		}
	}
}

func (self *Bsp) Contains(x, y int) bool {
	return x >= self.X && y >= self.Y && x < self.X+self.W && y < self.Y+self.H
}

func (self *Bsp) FindNode(x, y int) *Bsp {
	if !self.Contains(x, y) {
		return nil
	}
	if !self.IsLeaf() {
		var left, right *Bsp
		left = self.Left()
		if left.Contains(x, y) {
			return left.FindNode(x, y)
		}
		right = self.Right()
		if right.Contains(x, y) {
			return right.FindNode(x, y)
		}
	}
	return self
}

//
// HeightMap
//

type HeightMap struct {
	Data *C.TCOD_heightmap_t
}

func deleteHeightmap(h *HeightMap) {
	C.TCOD_heightmap_delete(h.Data)
}

func NewHeightMap(w, h int) *HeightMap {
	result := &HeightMap{C.TCOD_heightmap_new(C.int(w), C.int(h))}
	runtime.SetFinalizer(result, deleteHeightmap)
	return result
}

func (self *HeightMap) GetValue(x, y int) float32 {
	return float32(C.TCOD_heightmap_get_value(self.Data, C.int(x), C.int(y)))
}

func (self *HeightMap) GetWidth() int {
	return int(self.Data.w)
}

func (self *HeightMap) GetHeight() int {
	return int(self.Data.h)
}

func (self *HeightMap) GetInterpolatedValue(x, y float32) float32 {
	return float32(C.TCOD_heightmap_get_interpolated_value(self.Data, C.float(x), C.float(y)))
}

func (self *HeightMap) SetValue(x, y int, value float32) {
	C.TCOD_heightmap_set_value(self.Data, C.int(x), C.int(y), C.float(value))
}

func (self *HeightMap) GetNthValue(nth int) float32 {
	return float32(C._TCOD_heightmap_get_nth_value(self.Data, C.int(nth)))
}

func (self *HeightMap) SetNthValue(nth int, value float32) {
	C._TCOD_heightmap_set_nth_value(self.Data, C.int(nth), C.float(value))
}

func (self *HeightMap) GetSlope(x, y int) float32 {
	return float32(C.TCOD_heightmap_get_slope(self.Data, C.int(x), C.int(y)))
}

func (self *HeightMap) GetNormal(x, y float32, n *[3]float32, waterLevel float32) {
	C.TCOD_heightmap_get_normal(self.Data, C.float(x), C.float(y),
		(*C.float)(unsafe.Pointer(&n[0])),
		C.float(waterLevel))
}

func (self *HeightMap) CountCells(min, max float32) int {
	return int(C.TCOD_heightmap_count_cells(self.Data, C.float(min), C.float(max)))
}

func (self *HeightMap) HasLandOnBorder(waterLevel float32) bool {
	return toBool(C.TCOD_heightmap_has_land_on_border(self.Data, C.float(waterLevel)))
}

func (self *HeightMap) GetMinMax() (min, max float32) {
	var cmin, cmax C.float
	C.TCOD_heightmap_get_minmax(self.Data, &cmin, &cmax)
	min, max = float32(cmin), float32(cmax)
	return
}

func (self *HeightMap) Copy(source *HeightMap) {
	C.TCOD_heightmap_copy(source.Data, self.Data)
}

func (self *HeightMap) Add(value float32) {
	C.TCOD_heightmap_add(self.Data, C.float(value))
}

func (self *HeightMap) Scale(value float32) {
	C.TCOD_heightmap_scale(self.Data, C.float(value))
}

func (self *HeightMap) Clamp(min, max float32) {
	C.TCOD_heightmap_clamp(self.Data, C.float(min), C.float(max))
}

func (self *HeightMap) Normalize() {
	self.NormalizeRange(0, 1)
}

func (self *HeightMap) NormalizeRange(min, max float32) {
	C.TCOD_heightmap_normalize(self.Data, C.float(min), C.float(max))
}

func (self *HeightMap) Clear() {
	C.TCOD_heightmap_clear(self.Data)
}

func (self *HeightMap) Lerp(hm1 *HeightMap, hm2 *HeightMap, coef float32) {
	C.TCOD_heightmap_lerp_hm(hm1.Data, hm2.Data, self.Data, C.float(coef))
}

func (self *HeightMap) AddHm(hm1 *HeightMap, hm2 *HeightMap) {
	C.TCOD_heightmap_add_hm(hm1.Data, hm2.Data, self.Data)
}

func (self *HeightMap) Multiply(hm1 *HeightMap, hm2 *HeightMap) {
	C.TCOD_heightmap_multiply_hm(hm1.Data, hm2.Data, self.Data)
}

func (self *HeightMap) AddHill(hx, hy, hradius, hheight float32) {
	C.TCOD_heightmap_add_hill(self.Data, C.float(hx), C.float(hy), C.float(hradius), C.float(hheight))
}

func (self *HeightMap) DigHill(hx, hy, hradius, hheight float32) {
	C.TCOD_heightmap_dig_hill(self.Data, C.float(hx), C.float(hy), C.float(hradius), C.float(hheight))
}

func (self *HeightMap) DigBezier(px, py *[4]int, startRadius, startDepth, endRadius, endDepth float32) {
	C.TCOD_heightmap_dig_bezier(self.Data,
		(*C.int)(unsafe.Pointer(&px[0])),
		(*C.int)(unsafe.Pointer(&py[0])),
		C.float(startRadius), C.float(startDepth), C.float(endRadius), C.float(endDepth))
}

func (self *HeightMap) RainErosion(nbDrops int, erosionCoef, sedimentationCoef float32, rnd *Random) {
	C.TCOD_heightmap_rain_erosion(self.Data, C.int(nbDrops), C.float(erosionCoef), C.float(sedimentationCoef), rnd.Data)
}

func (self *HeightMap) KernelTransform(kernelsize int, dx, dy []int, weight []float32, minLevel, maxLevel float32) {
	C.TCOD_heightmap_kernel_transform(self.Data, C.int(kernelsize),
		(*C.int)(unsafe.Pointer(&dx[0])),
		(*C.int)(unsafe.Pointer(&dy[0])),
		(*C.float)(unsafe.Pointer(&weight[0])),
		C.float(minLevel),
		C.float(maxLevel))
}

func (self *HeightMap) AddVoronoi(nbPoints, nbCoef int, coef []float32, rnd *Random) {
	C.TCOD_heightmap_add_voronoi(self.Data, C.int(nbPoints), C.int(nbCoef), (*C.float)(unsafe.Pointer(&coef[0])), rnd.Data)
}

func (self *HeightMap) AddFbm(noise *Noise, mulx, muly, addx, addy, octaves, delta, scale float32) {
	C.TCOD_heightmap_add_fbm(self.Data, noise.Data, C.float(mulx),
		C.float(muly), C.float(addx), C.float(addy), C.float(octaves), C.float(delta), C.float(scale))
}

func (self *HeightMap) ScaleFbm(noise *Noise, mulx, muly, addx, addy, octaves, delta, scale float32) {
	C.TCOD_heightmap_scale_fbm(self.Data, noise.Data, C.float(mulx),
		C.float(muly), C.float(addx), C.float(addy), C.float(octaves), C.float(delta), C.float(scale))
}

func (self *HeightMap) Islandify(seaLevel float32, random *Random) {
	C.TCOD_heightmap_islandify(self.Data, C.float(seaLevel), random.Data)
}

//
// Image
//

type Image struct {
	Data C.TCOD_image_t
}

func deleteImage(img *Image) {
	C.TCOD_image_delete(img.Data)
}

func newImage(data C.TCOD_image_t) *Image {
	result := &Image{data}
	runtime.SetFinalizer(result, deleteImage)
	return result
}

func NewImage(width, height int) *Image {
	return newImage(C.TCOD_image_new(C.int(width), C.int(height)))
}

func NewImageFromConsole(console *Console) *Image {
	return newImage(C.TCOD_image_from_console(console.Data))
}

func (self *Image) RefreshConsole(console *Console) {
	C.TCOD_image_refresh_console(self.Data, console.Data)
}

func LoadImage(filename string) *Image {
	return newImage(C.TCOD_image_load(C.CString(filename)))
}

func (self *Image) Clear(color Color) {
	ccolor := fromColor(color)
	C._TCOD_image_clear(self.Data, (*C.TCOD_color_t)(&ccolor))
}

func (self *Image) Invert() {
	C.TCOD_image_invert(self.Data)
}

func (self *Image) Hflip() {
	C.TCOD_image_hflip(self.Data)
}

func (self *Image) Rotate90(numRotations int) {
	C.TCOD_image_rotate90(self.Data, C.int(numRotations))
}

func (self *Image) Vflip() {
	C.TCOD_image_vflip(self.Data)
}

func (self *Image) Scale(neww, newh int) {
	C.TCOD_image_scale(self.Data, C.int(neww), C.int(newh))
}

func (self *Image) Save(filename string) {
	C.TCOD_image_save(self.Data, C.CString(filename))
}

func (self *Image) GetSize(w, h *int) {
	var cw, ch C.int
	C.TCOD_image_get_size(self.Data, &cw, &ch)
	*w = int(cw)
	*h = int(ch)
}

func (self *Image) GetPixel(x, y int) Color {
	return toColor(C.TCOD_image_get_pixel(self.Data, C.int(x), C.int(y)))
}

func (self *Image) GetAlpha(x, y int) int {
	return int(C.TCOD_image_get_alpha(self.Data, C.int(x), C.int(y)))
}

func (self *Image) GetMipmapPixel(x0, y0, x1, y1 float32) Color {
	return toColor(C.TCOD_image_get_mipmap_pixel(self.Data, C.float(x0), C.float(y0),
		C.float(x1), C.float(y1)))
}

func (self *Image) PutPixel(x, y int, color Color) {
	ccolor := fromColor(color)
	C._TCOD_image_put_pixel(self.Data, C.int(x), C.int(y), (*C.TCOD_color_t)(&ccolor))
}

func (self *Image) Blit(console *Console, x, y float32, bkgndFlag BkgndFlag, scalex, scaley, angle float32) {
	C.TCOD_image_blit(self.Data, console.Data, C.float(x), C.float(y),
		C.TCOD_bkgnd_flag_t(bkgndFlag), C.float(scalex), C.float(scaley), C.float(angle))
}

func (self *Image) BlitRect(console *Console, x, y, w, h int, flag BkgndFlag) {
	C.TCOD_image_blit_rect(self.Data, console.Data, C.int(x), C.int(y), C.int(w), C.int(h), C.TCOD_bkgnd_flag_t(flag))
}

func (self *Image) Blit2x(dest *Console, dx, dy, sx, sy, w, h int) {
	C.TCOD_image_blit_2x(self.Data, dest.Data, C.int(dx), C.int(dy), C.int(sx), C.int(sy), C.int(w), C.int(h))
}

func (self *Image) SetKeyColor(keyColor Color) {
	ckeyColor := fromColor(keyColor)
	C._TCOD_image_set_key_color(self.Data, (*C.TCOD_color_t)(&ckeyColor))
}

func (self *Image) IsPixelTransparent(x, y int) bool {
	return toBool(C.TCOD_image_is_pixel_transparent(self.Data, C.int(x), C.int(y)))
}

//
// Path
//
//
type Path struct {
	Data C.TCOD_path_t
}

func deletePath(path *Path) {
	C.TCOD_path_delete(path.Data)
}

func NewPathUsingMap(m *Map, diagonalCost float32) *Path {
	result := &Path{C.TCOD_path_new_using_map(m.Data, C.float(diagonalCost))}
	runtime.SetFinalizer(result, deletePath)
	return result
}

// Not implemented - go not supporting callbacks
//func PathNewUsingFunction() {
//	//TCODLIB_API TCOD_path_t
//  TCOD_path_new_using_function(int map_width, int map_height, TCOD_path_func_t func, void *user_Data, float diagonalCost);
//}

func (self *Path) Compute(ox, oy, dx, dy int) bool {
	return toBool(C.TCOD_path_compute(self.Data, C.int(ox), C.int(oy), C.int(dx), C.int(dy)))
}

func (self *Path) Walk(recalcWhenNeeded bool) (x, y int) {
	var cx, cy C.int
	C.TCOD_path_walk(self.Data, &cx, &cy, fromBool(recalcWhenNeeded))
	x, y = int(cx), int(cy)
	return
}

func (self *Path) IsEmpty() bool {
	return toBool(C.TCOD_path_is_empty(self.Data))
}

func (self *Path) Size() int {
	return int(C.TCOD_path_size(self.Data))
}

func (self *Path) Get(index int) (x, y int) {
	var cx, cy C.int
	C.TCOD_path_get(self.Data, C.int(index), &cx, &cy)
	x, y = int(cx), int(cy)
	return
}

func (self *Path) GetOrigin() (x, y int) {
	var cx, cy C.int
	C.TCOD_path_get_origin(self.Data, &cx, &cy)
	x, y = int(cx), int(cy)
	return
}

func (self *Path) GetDestination() (x, y int) {
	var cx, cy C.int
	C.TCOD_path_get_destination(self.Data, &cx, &cy)
	x, y = int(cx), int(cy)
	return
}

//
// Dijkstra path
//

type Dijkstra struct {
	Data C.TCOD_dijkstra_t
}

func deleteDijkstra(d *Dijkstra) {
	C.TCOD_dijkstra_delete(d.Data)
}

func NewDijkstraUsingMap(m *Map, diagonalCost float32) *Dijkstra {
	result := &Dijkstra{C.TCOD_dijkstra_new(m.Data, C.float(diagonalCost))}
	runtime.SetFinalizer(result, deleteDijkstra)
	return result
}

// Not implemented - go not supporting callbacks
//func DijkstraNewUsingFunction() {
//	//TCODLIB_API TCOD_Dijkstra_t
//   TCOD_Dijkstra_new_using_function(int map_width, int map_height, TCOD_Dijkstra_func_t func, void *user_Data, float diagonalCost);
//}

func (self *Dijkstra) Compute(rootX, rootY int) {
	C.TCOD_dijkstra_compute(self.Data, C.int(rootX), C.int(rootY))
}

func (self *Dijkstra) GetDistance(x, y int) float32 {
	return float32(C.TCOD_dijkstra_get_distance(self.Data, C.int(x), C.int(y)))
}

func (self *Dijkstra) PathSet(x, y int) bool {
	return toBool(C.TCOD_dijkstra_path_set(self.Data, C.int(x), C.int(y)))
}

func (self *Dijkstra) IsEmpty() bool {
	return toBool(C.TCOD_dijkstra_is_empty(self.Data))
}

func (self *Dijkstra) Size() int {
	return int(C.TCOD_dijkstra_size(self.Data))
}

func (self *Dijkstra) Get(index int) (x, y int) {
	var cx, cy C.int
	C.TCOD_dijkstra_get(self.Data, C.int(index), &cx, &cy)
	x, y = int(cx), int(cy)
	return
}

func (self *Dijkstra) PathWalk() (x, y int) {
	var cx, cy C.int
	C.TCOD_dijkstra_path_walk(self.Data, &cx, &cy)
	x, y = int(cx), int(cy)
	return
}

//
// Mersenne Random generator
//

type RandomAlgo C.TCOD_random_algo_t

type Distribution C.TCOD_distribution_t

type Random struct {
	Data C.TCOD_random_t
}

type Dice struct {
	Data C.TCOD_dice_t
}

func fromDice(d Dice) C.TCOD_dice_t {
	return d.Data
}

func toDice(d C.TCOD_dice_t) Dice {
	return Dice{d}
}

func deleteRandom(r *Random) {
	C.TCOD_random_delete(r.Data)
}

func newRandom(data C.TCOD_random_t) *Random {
	result := &Random{data}
	runtime.SetFinalizer(result, deleteRandom)
	return result
}

func GetRandomInstance() *Random {
	return newRandom(C.TCOD_random_get_instance())
}

func NewRandom() *Random {
	return newRandom(C.TCOD_random_new(C.TCOD_random_algo_t(RNG_MT)))
}

func NewRandomWithAlgo(algo RandomAlgo) *Random {
	return newRandom(C.TCOD_random_new(C.TCOD_random_algo_t(algo)))
}

func NewRandomFromSeedWithAlgo(seed uint32, algo RandomAlgo) *Random {
	return newRandom(C.TCOD_random_new_from_seed(C.TCOD_random_algo_t(algo), C.uint32(seed)))
}

func NewRandomFromSeed(seed uint32) *Random {
	return newRandom(
		C.TCOD_random_new_from_seed(
			C.TCOD_random_algo_t(RNG_MT),
			C.uint32(seed)))
}

func (self *Random) Save() *Random {
	result := newRandom(C.TCOD_random_save(self.Data))
	return result
}

func (self *Random) Restore(backup *Random) {
	C.TCOD_random_restore(self.Data, backup.Data)
}

func (self *Random) SetDistribution(distribution Distribution) {
	C.TCOD_random_set_distribution(self.Data, C.TCOD_distribution_t(distribution))
}

func (self *Random) GetInt(min, max int) int {
	return int(C.TCOD_random_get_int(self.Data, C.int(min), C.int(max)))
}

func (self *Random) GetFloat(min, max float32) float32 {
	return float32(C.TCOD_random_get_float(self.Data, C.float(min), C.float(max)))
}

func (self *Random) GetDouble(min, max float64) float64 {
	return float64(C.TCOD_random_get_double(self.Data, C.double(min), C.double(max)))
}

func (self *Random) GetIntMean(min, max, mean int) int {
	return int(C.TCOD_random_get_int_mean(self.Data, C.int(min), C.int(max), C.int(mean)))
}

func (self *Random) GetFloatMean(min, max, mean float32) float32 {
	return float32(C.TCOD_random_get_float_mean(self.Data, C.float(min), C.float(max), C.float(mean)))
}

func (self *Random) GetDoubleMean(min, max, mean float64) float64 {
	return float64(C.TCOD_random_get_double_mean(self.Data, C.double(min), C.double(max), C.double(mean)))
}

func NewDice(s string) *Dice {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	result := &Dice{Data: C.TCOD_random_dice_new(cs)}
	return result
}

func (self *Dice) Roll(random *Random) int {
	return int(C.TCOD_random_dice_roll(random.Data, self.Data))
}

func RollDice(random *Random, s string) int {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return int(C.TCOD_random_dice_roll_s(random.Data, cs))
}

//
// Parser library
//

type ParserValueType C.TCOD_value_type_t

type ParserStruct struct {
	Data C.TCOD_parser_struct_t
}

type Parser struct {
	Data C.TCOD_parser_t
}

type ParserProperty struct {
	Name      string
	ValueType ParserValueType
	Value     interface{}
}

func (self ParserStruct) GetName() string {
	return C.GoString(C.TCOD_struct_get_name(self.Data))
}

func (self ParserStruct) AddProperty(name string, valueType ParserValueType, mandatory bool) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	C.TCOD_struct_add_property(self.Data, cname, C.TCOD_value_type_t(valueType), fromBool(mandatory))
}

func (self ParserStruct) AddListProperty(name string, valueType ParserValueType, mandatory bool) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.TCOD_struct_add_list_property(self.Data, cname, C.TCOD_value_type_t(valueType), fromBool(mandatory))
}

func (self ParserStruct) AddValueList(name string, valueList []string, mandatory bool) {
	cvalueList := make([]*C.char, len(valueList))
	for i := range valueList {
		cvalueList[i] = C.CString(valueList[i])
	}
	C.TCOD_struct_add_value_list_sized(self.Data, C.CString(name),
		(**C.char)(unsafe.Pointer(&cvalueList[0])), C.int(len(valueList)), fromBool(mandatory))

	for i := range cvalueList {
		C.free(unsafe.Pointer(cvalueList[i]))
	}

}

func (self ParserStruct) AddFlag(propname string) {
	cpropname := C.CString(propname)
	defer C.free(unsafe.Pointer(cpropname))
	C.TCOD_struct_add_flag(self.Data, cpropname)
}

func (self ParserStruct) AddStructure(substruct ParserStruct) {
	// TODO is this necessary ??
	//	struct1 := self.Data
	//	substruct2 := struct_.Data
	C._TCOD_struct_add_structure(&self.Data, &substruct.Data)
}

func (self *ParserStruct) IsMandatory(propname string) bool {
	cpropname := C.CString(propname)
	defer C.free(unsafe.Pointer(cpropname))
	return toBool(C.TCOD_struct_is_mandatory(self.Data, cpropname))
}

func (self *ParserStruct) GetType(propname string) ParserValueType {
	cpropname := C.CString(propname)
	defer C.free(unsafe.Pointer(cpropname))
	return ParserValueType(C.TCOD_struct_get_type(self.Data, cpropname))
}

func deleteParser(p *Parser) {
	C.TCOD_parser_delete(p.Data)
}

func NewParser() *Parser {
	result := &Parser{C.TCOD_parser_new()}
	runtime.SetFinalizer(result, deleteParser)
	return result
}

func (self *Parser) RegisterStruct(name string) ParserStruct {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return ParserStruct{C.TCOD_parser_new_struct(self.Data, cname)}
}

// TODO custom parsers are not supported
// TCODLIB_API TCOD_value_type_t TCOD_parser_new_custom_type(TCOD_parser_t parser,TCOD_parser_custom_t custom_type_parser);

// TODO listeners are not supported
// Running parser return list of parsed properties
func (self *Parser) Run(filename string) []ParserProperty {
	// run parser with default listeners
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	C.TCOD_parser_run(self.Data, cfilename, nil)

	// extract properties to Go structures
	var cprop *C._prop_t
	var prop ParserProperty
	var l C.TCOD_list_t = C.TCOD_list_t(((*C.TCOD_parser_int_t)(self.Data)).props)
	result := make([]ParserProperty, C.TCOD_list_size(l))

	for i := 0; i < int(C.TCOD_list_size(l)); i++ {

		cprop = (*C._prop_t)(unsafe.Pointer(C.TCOD_list_get(l, C.int(i))))

		prop.Name = C.GoString(cprop.name)
		prop.ValueType = ParserValueType(cprop.value_type)
		if cprop.value_type == TYPE_STRING ||
			(cprop.value_type >= TYPE_VALUELIST00 && cprop.value_type <= TYPE_VALUELIST15) {
			prop.Value = C.GoString(*((**C.char)(unsafe.Pointer(&cprop.value))))
		} else if cprop.value_type == TYPE_INT {
			prop.Value = int(*((*C.int)(unsafe.Pointer(&cprop.value))))
		} else if cprop.value_type == TYPE_FLOAT {
			prop.Value = float32(*((*C.float)(unsafe.Pointer(&cprop.value))))
		} else if cprop.value_type == TYPE_BOOL {
			prop.Value = toBool(*((*C.bool)(unsafe.Pointer(&cprop.value))))
		} else if cprop.value_type == TYPE_COLOR {
			prop.Value = toColor(*((*C.TCOD_color_t)(unsafe.Pointer(&cprop.value))))
		} else if cprop.value_type == TYPE_DICE {
			prop.Value = toDice(*((*C.TCOD_dice_t)(unsafe.Pointer(&cprop.value))))
		} else if cprop.value_type >= TYPE_LIST {
			elType := cprop.value_type - TYPE_LIST
			elList := C.TCOD_list_t(*(*C.TCOD_list_t)(unsafe.Pointer(&cprop.value)))
			elListSize := int(C.TCOD_list_size(elList))
			if elType == TYPE_STRING {
				prop.Value = make([]string, elListSize)
				for j := 0; j < elListSize; j++ {
					elValue := (*C.char)(unsafe.Pointer(C.TCOD_list_get(elList, C.int(j))))
					prop.Value.([]string)[j] = C.GoString(elValue)
				}
			} else if elType == TYPE_INT {
				prop.Value = make([]int, elListSize)
				for j := 0; j < elListSize; j++ {
					elValue := C.TCOD_list_get(elList, C.int(j))
					prop.Value.([]int)[j] = int(*(*C.int)(unsafe.Pointer(&elValue)))
				}
			} else if elType == TYPE_FLOAT {
				prop.Value = make([]float32, elListSize)
				for j := 0; j < elListSize; j++ {
					elValue := C.TCOD_list_get(elList, C.int(j))
					prop.Value.([]float32)[j] = float32(*(*C.float)(unsafe.Pointer(&elValue)))
				}
			} else if elType == TYPE_BOOL {
				prop.Value = make([]bool, elListSize)
				for j := 0; j < elListSize; j++ {
					elValue := C.TCOD_list_get(elList, C.int(j))
					prop.Value.([]bool)[j] = toBool(*(*C.bool)(unsafe.Pointer(&elValue)))
				}
			} else if elType == TYPE_DICE {
				prop.Value = make([]Dice, elListSize)
				for j := 0; j < elListSize; j++ {
					elValue := *(*C.TCOD_dice_t)(unsafe.Pointer(C.TCOD_list_get(elList, C.int(j))))
					prop.Value.([]Dice)[j] = toDice(elValue)
				}
			} else if elType == TYPE_COLOR {
				prop.Value = make([]Color, elListSize)
				for j := 0; j < elListSize; j++ {
					elValue := *(*C.TCOD_color_t)(unsafe.Pointer(C.TCOD_list_get(elList, C.int(j))))
					prop.Value.([]Color)[j] = toColor(elValue)
				}
			}
		}
		result[i] = prop
	}
	return result
}

//
// Perlin noise
//

// Noise NEW 

const NOISE_MAX_OCTAVES = 128
const NOISE_MAX_DIMENSIONS = 4
const NOISE_DEFAULT_HURST = 0.5
const NOISE_DEFAULT_LACUNARITY = 2.0

type NoiseType C.TCOD_noise_type_t

type Noise struct {
	Data C.TCOD_noise_t
}

type FloatArray []float32

func deleteNoise(noise *Noise) {
	C.TCOD_noise_delete(noise.Data)
}

func newNoise(d C.TCOD_noise_t) *Noise {
	result := &Noise{d}
	runtime.SetFinalizer(result, deleteNoise)
	return result
}

func NewNoise(dimensions int, random *Random) *Noise {
	return newNoise(C.TCOD_noise_new(C.int(dimensions), C.float(NOISE_DEFAULT_HURST),
		C.float(NOISE_DEFAULT_LACUNARITY), random.Data))
}

func NewNoiseWithOptions(dimensions int, hurst float32, lacunarity float32, random *Random) *Noise {
	return newNoise(C.TCOD_noise_new(C.int(dimensions), C.float(hurst), C.float(lacunarity), random.Data))
}

func (self *Noise) GetEx(f FloatArray, noiseType NoiseType) float32 {
	return float32(C.TCOD_noise_get_ex(self.Data, (*C.float)(unsafe.Pointer(&f[0])), C.TCOD_noise_type_t(noiseType)))
}

func (self *Noise) SetType(noiseType NoiseType) {
	C.TCOD_noise_set_type(self.Data, C.TCOD_noise_type_t(noiseType))
}

func (self *Noise) GetFbmEx(f FloatArray, octaves float32, noiseType NoiseType) float32 {
	return float32(C.TCOD_noise_get_fbm_ex(self.Data, (*C.float)(unsafe.Pointer(&f[0])), C.float(octaves),
		C.TCOD_noise_type_t(noiseType)))
}

func (self *Noise) GetTurbulenceEx(f FloatArray, octaves float32, noiseType NoiseType) float32 {
	return float32(C.TCOD_noise_get_turbulence_ex(self.Data, (*C.float)(unsafe.Pointer(&f[0])), C.float(octaves),
		C.TCOD_noise_type_t(noiseType)))
}

func (self *Noise) Get(f FloatArray) float32 {
	return float32(C.TCOD_noise_get(self.Data, (*C.float)(unsafe.Pointer(&f[0]))))
}

func (self *Noise) GetFbm(f FloatArray, octaves float32) float32 {
	return float32(C.TCOD_noise_get_fbm(self.Data, (*C.float)(unsafe.Pointer(&f[0])), C.float(octaves)))
}

func (self *Noise) GetTurbulence(f FloatArray, octaves float32) float32 {
	return float32(C.TCOD_noise_get_turbulence(self.Data, (*C.float)(unsafe.Pointer(&f[0])), C.float(octaves)))
}

//
// Zip
//

type Zip struct {
	Data C.TCOD_zip_t
}

func deleteZip(zip *Zip) {
	C.TCOD_zip_delete(zip.Data)
}

func NewZip() *Zip {
	result := &Zip{C.TCOD_zip_new()}
	runtime.SetFinalizer(result, deleteZip)
	return result
}

// output interface

func (self *Zip) PutChar(val byte) {
	C.TCOD_zip_put_char(self.Data, C.char(val))
}

func (self *Zip) PutInt(val int) {
	C.TCOD_zip_put_int(self.Data, C.int(val))
}

func (self *Zip) PutFloat(val float32) {
	C.TCOD_zip_put_float(self.Data, C.float(val))
}

func (self *Zip) PutString(val string) {
	cval := C.CString(val)
	defer C.free(unsafe.Pointer(cval))
	C.TCOD_zip_put_string(self.Data, cval)
}

func (self *Zip) PutColor(val Color) {
	cval := fromColor(val)
	C._TCOD_zip_put_color(self.Data, (*C.TCOD_color_t)(&cval))
}

func (self *Zip) PutImage(val *Image) {
	C.TCOD_zip_put_image(self.Data, val.Data)
}

func (self *Zip) PutConsole(val *Console) {
	C.TCOD_zip_put_console(self.Data, val.Data)
}

func (self *Zip) PutData(nbBytes int, data unsafe.Pointer) {
	C.TCOD_zip_put_data(self.Data, C.int(nbBytes), data)
}

func (self *Zip) GetCurrentBytes() uint32 {
	return uint32(C.TCOD_zip_get_current_bytes(self.Data))
}

func (self *Zip) SaveToFile(filename string) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	C.TCOD_zip_save_to_file(self.Data, cfilename)
}

// input interface

func (self *Zip) LoadFromFile(filename string) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	C.TCOD_zip_load_from_file(self.Data, cfilename)
}

func (self *Zip) GetChar() byte {
	return byte(C.TCOD_zip_get_char(self.Data))
}

func (self *Zip) GetInt() int {
	return int(C.TCOD_zip_get_int(self.Data))
}

func (self *Zip) GetFloat() float32 {
	return float32(C.TCOD_zip_get_float(self.Data))
}

func (self *Zip) GetString() string {
	return C.GoString(C.TCOD_zip_get_string(self.Data))
}

func (self *Zip) GetColor() Color {
	return toColor(C.TCOD_zip_get_color(self.Data))
}

func (self *Zip) GetImage() *Image {
	return &Image{C.TCOD_zip_get_image(self.Data)}
}

func (self *Zip) GetConsole() *Console {
	return &Console{C.TCOD_zip_get_console(self.Data)}
}

func (self *Zip) GetData(nbBytes int, data unsafe.Pointer) int {
	return int(C.TCOD_zip_get_data(self.Data, C.int(nbBytes), data))
}

func (self *Zip) GetRemainingBytes() uint32 {
	return uint32(C.TCOD_zip_get_remaining_bytes(self.Data))
}

func (self *Zip) SkipBytes(nbBytes uint32) {
	C.TCOD_zip_skip_bytes(self.Data, C.uint32(nbBytes))
}
