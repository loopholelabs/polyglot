ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
VERSION=$(shell cat $(ROOT_DIR)/version/current_version)
VERSION_MAJOR=$(shell echo $(VERSION) | cut -c2)
VERSION_MINOR=$(shell echo $(VERSION) | cut -c4)
VERSION_MICRO=$(shell echo $(VERSION) | cut -c6)
CLIB_SO_DEV=libpolyglot.so
CLIB_SO_MAN=$(CLIB_SO_DEV).$(VERSION_MAJOR)
CLIB_SO_FULL=$(CLIB_SO_DEV).$(VERSION)
CLIB_HEADER=polyglot.h
CLIB_SO_DEV_RELEASE=c_bindings/target/release/$(CLIB_SO_DEV)
CLIB_SO_DEV_DEBUG=c_bindings/target/debug/$(CLIB_SO_DEV)
CLIB_PKG_CONFIG=polyglot.pc
PREFIX ?= /usr/local

outdir ?= $(ROOT_DIR)

CPU_BITS = $(shell getconf LONG_BIT)
ifeq ($(CPU_BITS), 32)
    LIBDIR ?= $(PREFIX)/lib
else
    LIBDIR ?= $(PREFIX)/lib$(CPU_BITS)
endif

INCLUDE_DIR ?= $(PREFIX)/include
PKG_CONFIG_LIBDIR ?= $(LIBDIR)/pkgconfig
MAN_DIR ?= $(PREFIX)/share/man

RELEASE ?=0

.PHONY: clib_debug
clib_debug:
	cd c_bindings
	cargo build --all
	ln -sfv $(CLIB_SO_DEV) target/debug/$(CLIB_SO_FULL)
	ln -sfv $(CLIB_SO_DEV) target/debug/$(CLIB_SO_MAN)
	cd -

$(CLIB_SO_DEV_RELEASE):
	cargo build --all --release

$(CLIB_SO_DEV_DEBUG): clib_debug

clib: $(CLIB_HEADER) $(CLIB_SO_DEV_RELEASE) $(CLIB_PKG_CONFIG)

.PHONY: $(CLIB_HEADER)
$(CLIB_HEADER): $(CLIB_HEADER).in
	cp $(CLIB_HEADER).in $(CLIB_HEADER)
	sed -i '' 's/@_VERSION_MAJOR@/$(VERSION_MAJOR)/' \
		$(CLIB_HEADER)
	sed -i '' 's/@_VERSION_MINOR@/$(VERSION_MINOR)/' \
		$(CLIB_HEADER)
	sed -i '' 's/@_VERSION_MICRO@/$(VERSION_MICRO)/' \
		$(CLIB_HEADER)

.PHONY: $(CLIB_PKG_CONFIG)
$(CLIB_PKG_CONFIG): $(CLIB_PKG_CONFIG).in
	cp $(CLIB_PKG_CONFIG).in $(CLIB_PKG_CONFIG)
	sed -i '' 's|@_VERSION_MAJOR@|$(VERSION_MAJOR)|' $(CLIB_PKG_CONFIG)
	sed -i '' 's|@_VERSION_MINOR@|$(VERSION_MINOR)|' $(CLIB_PKG_CONFIG)
	sed -i '' 's|@_VERSION_MICRO@|$(VERSION_MICRO)|' $(CLIB_PKG_CONFIG)
	sed -i '' 's|@PREFIX@|$(PREFIX)|' $(CLIB_PKG_CONFIG)
	sed -i '' 's|@LIBDIR@|$(LIBDIR)|' $(CLIB_PKG_CONFIG)
	sed -i '' 's|@INCLUDE_DIR@|$(INCLUDE_DIR)|' $(CLIB_PKG_CONFIG)

install: clib
	install -p -D -m755 $(CLIB_SO_DEV_RELEASE) \
		$(DESTDIR)$(LIBDIR)/$(CLIB_SO_FULL)
	ln -sfv $(CLIB_SO_FULL) $(DESTDIR)$(LIBDIR)/$(CLIB_SO_MAN)
	ln -sfv $(CLIB_SO_FULL) $(DESTDIR)$(LIBDIR)/$(CLIB_SO_DEV)
	install -p -v -D -m644 $(CLIB_HEADER) \
		$(DESTDIR)$(INCLUDE_DIR)/$(shell basename $(CLIB_HEADER))
	install -p -v -D -m644 $(CLIB_PKG_CONFIG) \
		$(DESTDIR)$(PKG_CONFIG_LIBDIR)/$(shell basename $(CLIB_PKG_CONFIG))