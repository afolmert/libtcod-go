CMDS=sample hmtool

GOMAKE=gomake

TARG=libtcod-go

SUB=$(LIBS:%=pkg/%) $(CMDS:%=cmd/%)

all: deps build.cmds

run: $(TARG).run

build.cmds: $(addsuffix .build, $(CMDS))
clean.cmds: $(addsuffix .clean, $(CMDS))

build.libs:
	$(GOMAKE) -C pkg all

test: deps build.libs
	$(GOMAKE) -C pkg test

# XXX: Hardwired to clean the command before build to hack around problems
# specifying library dependencies to the command.
%.run: build.libs
	$(GOMAKE) -C cmd/$* clean
	$(GOMAKE) -C cmd/$* all
	(cd ./cmd/$*; ./$* ${ARGS})

%.build: build.libs
	$(GOMAKE) -C cmd/$*

%.clean:
	$(GOMAKE) -C cmd/$* clean

clean: clean.cmds
	$(GOMAKE) -C pkg clean

nuke: clean.cmds
	$(GOMAKE) -C pkg nuke

deps:
	$(GOMAKE) -C pkg deps
