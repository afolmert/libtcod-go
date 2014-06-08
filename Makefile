all: deps install

DIRS=\
	tcod\

NOTEST=\
	tcod\

TEST=$(filter-out $(NOTEST),$(DIRS))

GOMAKE=gomake

clean.dirs: $(addsuffix .clean, $(DIRS))
install.dirs: $(addsuffix .install, $(DIRS))
nuke.dirs: $(addsuffix .nuke, $(DIRS))
test.dirs: $(addsuffix .test, $(TEST))

%.clean:
	$(GOMAKE) -C $* clean

%.install:
	$(GOMAKE) -C $* install

%.nuke:
	$(GOMAKE) -C $* nuke

%.test:
	$(GOMAKE) -C $* test

clean: clean.dirs

install: install.dirs

test: test.dirs

nuke: nuke.dirs

deps:
	$(GOROOT)/src/pkg/deps.bash

-include Make.deps
