# Copyright 2009 The Go Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

TARG=tcod

GOFILES:=tcod_defs.go gui.go
CGOFILES:=tcod.go


CGO_LDFLAGS=-ltcod

CGO_CFLAGS=-I./include
CGO_CFLAGS+=-I.

include $(GOROOT)/src/Make.pkg


tcod_defs.go: tcod.c
	godefs -g tcod tcod.c -I./include > tcod_defs.go

