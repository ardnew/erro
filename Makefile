# ------------------------------------------------------------------------------
#  project configuration (symbols exported verbatim via Go linker)

PROJECT   ?= erro
IMPORT    ?= github.com/ardnew/$(PROJECT)
VERSION   ?= 0.2.3
BRANCH    ?= $(shell git symbolic-ref --short HEAD)
REVISION  ?= $(shell git rev-parse --short HEAD)
BUILDTIME ?= $(shell date -u '+%FT%TZ')
TARGET    ?= $(shell uname -s | tr 'A-Z' 'a-z')
PLATFORM  ?= $(TARGET)-amd64
DEFAULT   ?= zip # only for "release" target

# default output paths
BINPATH ?= bin
PKGPATH ?= pkg

# consider all Go source files recursively from working dir
SOURCES ?= $(shell find . -type f -iname '*.go')

# Go modules support
GOMODFILE ?= go.mod
GOSUMFILE ?= go.sum

# other non-Go source files that may affect build staleness
METASOURCES ?= Makefile $(GOMODFILE) $(GOSUMFILE)

# other files to include with distribution packages
EXTRAFILES ?= LICENSE README.md

# Go package import path where the exported symbols will be defined
MAINPKGPATH ?= main

# Makefile identifiers to export (as strings) via Go linker
EXPORTS ?= PROJECT VERSION BRANCH REVISION BUILDTIME PLATFORM

# ------------------------------------------------------------------------------
#  constants and derived variables

# supported platforms (GOARCH-GOOS)
platforms :=                                           \
  linux-amd64 linux-386 linux-arm64 linux-arm          \
  darwin-amd64 darwin-arm64                            \
  windows-amd64 windows-386                            \
  freebsd-amd64 freebsd-386 freebsd-arm                \
  android-arm64

show-platforms := $(addprefix show-,$(platforms))

# invalid build target provided
ifeq "" "$(strip $(filter $(platforms),$(PLATFORM)))"
$(error unsupported PLATFORM "$(PLATFORM)" (see: "make help"))
endif

# parse OS (linux, darwin, ...) and arch (386, amd64, ...) from PLATFORM
os   := $(word 1,$(subst -, ,$(PLATFORM)))
arch := $(word 2,$(subst -, ,$(PLATFORM)))

# output file extensions
binext := $(if $(filter windows,$(os)),.exe,)
tgzext := .tar.gz
tbzext := .tar.bz2
zipext := .zip

# system commands
echo  := echo
test  := test
make  := make
cd    := cd
rm    := rm -rvf
mv    := mv -v
cp    := cp -rv
mkdir := mkdir -pv
chmod := chmod -v
tail  := tail
grep  := command grep
go    := GOOS="$(os)" GOARCH="$(arch)" go
git   := git
gh    := gh
tgz   := tar -czvf
tbz   := tar -cjvf
zip   := zip -vr

# go build flags: export variables as strings to the selected package
goflags ?= -v -ldflags='-w -s $(foreach %,$(EXPORTS),-X "$(MAINPKGPATH).$(%)=$($(%))")'

# output paths
bindir := $(BINPATH)/$(PLATFORM)
binexe := $(bindir)/$(PROJECT)$(binext)
pkgver := $(PKGPATH)/$(VERSION)
triple := $(PROJECT)-$(VERSION)-$(PLATFORM)

# Since it isn't possible to pass arguments from make to the target executable
# (without, e.g., inline variable definitions), we simply use a separate shell
# script that builds the project and calls the executable.
# You can thus call this shell script, and all arguments will be passed along.
# Use the 'make run' target to generate this script in the project root.
runsh := run.sh
define RUNSH
#!/bin/sh
# Description:
# 	Rebuild and run $(binexe) with command-line arguments.
# 
# Usage:
# 	./$(runsh) [arg ...]
# 
if make build > /dev/null; then
	"$(binexe)" "$${@}"
fi
endef
export RUNSH

# ------------------------------------------------------------------------------
#  make targets

.PHONY: all
all: build

.PHONY: clean
clean:
	$(rm) "$(bindir)" "$(pkgver)/$(triple)"
	$(go) clean

.PHONY: tidy
tidy: $(GOMODFILE) $(GOSUMFILE)
	$(go) mod tidy

.PHONY: build
build: $(binexe)

.PHONY: vet
vet: $(SOURCES) $(METASOURCES)
	$(go) vet

.PHONY: run
run: $(runsh)

.PHONY: $(show-platforms)
$(show-platforms):
	@$(echo) $(subst show-,,$(@))

.PHONY: show-platforms
show-platforms: | $(show-platforms)

$(GOMODFILE):
	$(go) mod init $(IMPORT)

$(GOSUMFILE):
	$(go) mod download

$(bindir) $(pkgver) $(pkgver)/$(triple):
	@$(test) -d "$(@)" || $(mkdir) "$(@)"

$(binexe): $(SOURCES) $(METASOURCES) tidy | $(bindir)
	$(go) build -o "$(@)" $(goflags)
	@$(echo) " -- success: $(@)"

$(runsh):
	@$(echo) "$$RUNSH" > "$(@)"
	@$(chmod) +x "$(@)"
	@$(echo) " -- success: $(@)"
	@$(echo)
	@$(tail) -n +2 "$(@)" | $(grep) -oP '^#\K.*'

# ------------------------------------------------------------------------------
#  targets for creating versioned packages (.zip, .tar.gz, or .tar.bz2)

.PHONY: release
release: artifact = $(or $(RELEASE),$(DEFAULT))
release:
	$(test) -z "$$( $(git) status --porcelain=v1 )" || \
	  { $(echo) "working tree contains modified files"; false; }
	@# make target $(RELEASE) for all platforms (or $(DEFAULT) if undefined)
	for p in $$( $(make) show-platforms ); do \
	  $(make) PLATFORM="$${p}" $(artifact); done
	$(gh) release create v$(VERSION) --generate-notes $(pkgver)/*.$(artifact)

.PHONY: zip
zip: $(EXTRAFILES) $(pkgver)/$(triple)$(zipext)

$(pkgver)/%$(zipext): $(binexe) $(pkgver)/%
	$(cp) "$(<)" $(EXTRAFILES) "$(@D)/$(*)"
	@$(cd) "$(@D)" && $(zip) "$(*)$(zipext)" "$(*)"

.PHONY: tgz
tgz: $(EXTRAFILES) $(pkgver)/$(triple)$(tgzext)

$(pkgver)/%$(tgzext): $(binexe) $(pkgver)/%
	$(cp) "$(<)" $(EXTRAFILES) "$(@D)/$(*)"
	@$(cd) "$(@D)" && $(tgz) "$(*)$(tgzext)" "$(*)"

.PHONY: tbz
tbz: $(EXTRAFILES) $(pkgver)/$(triple)$(tbzext)

$(pkgver)/%$(tbzext): $(binexe) $(pkgver)/%
	$(cp) "$(<)" $(EXTRAFILES) "$(@D)/$(*)"
	@$(cd) "$(@D)" && $(tbz) "$(*)$(tbzext)" "$(*)"

