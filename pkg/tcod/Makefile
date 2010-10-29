# Copyright 2009 The Go Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

TARG=tcod

# Can have plain GOFILES too, but this example doesn't.

GOFILES=tcod_defs.go gui.go
CGOFILES=tcod.go

CGO_LDFLAGS=-ltcod
# -ltcodgui

CGO_CFLAGS=-I./include
CGO_CFLAGS+=-I.
# CGO_LDFLAGS+=-L/usr/lib


LDPATH_linux=-Wl,-R,`pwd`
#
# To add flags necessary for locating the library or its include files,
# set CGO_CFLAGS or CGO_LDFLAGS.  For example, to use an
# alternate installation of the library:
#	CGO_CFLAGS=-I/home/rsc/gmp32/include
#	CGO_LDFLAGS+=-L/home/rsc/gmp32/lib
# Note the += on the second line.

CLEANFILES+=

include $(GOROOT)/src/Make.pkg

tcod_defs.go: tcod.c
	godefs -g tcod tcod.c -I./include > tcod_defs.go



