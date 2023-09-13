ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
VERSION=$(shell cat $(ROOT_DIR)/version/current_version)
VERSION_MAJOR=$(shell echo $(VERSION) | cut -c2)
VERSION_MINOR=$(shell echo $(VERSION) | cut -c4)
VERSION_MICRO=$(shell echo $(VERSION) | cut -c6)
CLIB_SO=libpolyglot
CLIB_SO_DEV=$(CLIB_SO).so
CLIB_SO_MAN=$(CLIB_SO_DEV).$(VERSION_MAJOR)
CLIB_SO_FULL=$(CLIB_SO_DEV).$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_MICRO)
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
	cd c_bindings && cargo build --all && ([ -f ./target/debug/$(CLIB_SO).dylib ] && (echo "DYLIB exists, copying" && cp ./target/debug/$(CLIB_SO).dylib  ./target/debug/$(CLIB_SO_DEV)) || echo "DYLIB does not exist, ignoring") && ln -sfv $(CLIB_SO_DEV) target/debug/$(CLIB_SO_FULL) && ln -sfv $(CLIB_SO_DEV) target/debug/$(CLIB_SO_MAN) && cd -

$(CLIB_SO_DEV_RELEASE):
	cd c_bindings && cargo build --all --release && ([ -f ./target/release/$(CLIB_SO).dylib ] && (echo "DYLIB exists, copying" && cp ./target/release/$(CLIB_SO).dylib  ./target/release/$(CLIB_SO_DEV)) || echo "DYLIB does not exist, ignoring") && cd -

$(CLIB_SO_DEV_DEBUG): clib_debug

clib: $(CLIB_HEADER) $(CLIB_SO_DEV_RELEASE) $(CLIB_PKG_CONFIG)

.PHONY: $(CLIB_HEADER)
$(CLIB_HEADER): $(CLIB_HEADER).in
	cp $(CLIB_HEADER).in $(CLIB_HEADER)
	sed -i -e 's/@_VERSION_MAJOR@/$(VERSION_MAJOR)/' \
		$(CLIB_HEADER)
	sed -i -e 's/@_VERSION_MINOR@/$(VERSION_MINOR)/' \
		$(CLIB_HEADER)
	sed -i -e 's/@_VERSION_MICRO@/$(VERSION_MICRO)/' \
		$(CLIB_HEADER)

.PHONY: $(CLIB_PKG_CONFIG)
$(CLIB_PKG_CONFIG): $(CLIB_PKG_CONFIG).in
	cp $(CLIB_PKG_CONFIG).in $(CLIB_PKG_CONFIG)
	sed -i -e 's|@_VERSION_MAJOR@|$(VERSION_MAJOR)|' $(CLIB_PKG_CONFIG)
	sed -i -e 's|@_VERSION_MINOR@|$(VERSION_MINOR)|' $(CLIB_PKG_CONFIG)
	sed -i -e 's|@_VERSION_MICRO@|$(VERSION_MICRO)|' $(CLIB_PKG_CONFIG)
	sed -i -e 's|@PREFIX@|$(PREFIX)|' $(CLIB_PKG_CONFIG)
	sed -i -e 's|@LIBDIR@|$(LIBDIR)|' $(CLIB_PKG_CONFIG)
	sed -i -e 's|@INCLUDE_DIR@|$(INCLUDE_DIR)|' $(CLIB_PKG_CONFIG)

.PHONY: clib_test
clib_test: $(CLIB_SO_DEV_DEBUG) $(CLIB_HEADER)
	$(eval TMPDIR := $(shell mktemp -d))
	cp $(CLIB_SO_DEV_DEBUG) $(TMPDIR)/$(CLIB_SO_FULL)
	ln -sfv $(CLIB_SO_FULL) $(TMPDIR)/$(CLIB_SO_MAN)
	ln -sfv $(CLIB_SO_FULL) $(TMPDIR)/$(CLIB_SO_DEV)
	cp $(CLIB_HEADER) $(TMPDIR)/$(shell basename $(CLIB_HEADER))
	ls -a $(TMPDIR)
	gcc -g -Wall -Wextra -L$(TMPDIR) -I$(TMPDIR) -o $(TMPDIR)/polyglot_test c_bindings/tests/polyglot_test.c -lpolyglot
	$(TMPDIR)/polyglot_test
	rm -rf $(TMPDIR)

install: clib
	mkdir -p $(DESTDIR)$(LIBDIR)/$(CLIB_SO_FULL)
	install -p -m755 $(CLIB_SO_DEV_RELEASE) $(DESTDIR)$(LIBDIR)/$(CLIB_SO_FULL)
	ln -sfv $(CLIB_SO_FULL) $(DESTDIR)$(LIBDIR)/$(CLIB_SO_MAN)
	ln -sfv $(CLIB_SO_FULL) $(DESTDIR)$(LIBDIR)/$(CLIB_SO_DEV)
	mkdir -p $(DESTDIR)$(INCLUDE_DIR)/$(shell basename $(CLIB_HEADER))
	install -p -v -m644 $(CLIB_HEADER) $(DESTDIR)$(INCLUDE_DIR)/$(shell basename $(CLIB_HEADER))
	mkdir -p $(DESTDIR)$(PKG_CONFIG_LIBDIR)/$(shell basename $(CLIB_PKG_CONFIG))
	install -p -v -m644 $(CLIB_PKG_CONFIG) $(DESTDIR)$(PKG_CONFIG_LIBDIR)/$(shell basename $(CLIB_PKG_CONFIG))

uninstall:
	- rm -rfv $(DESTDIR)$(LIBDIR)/$(CLIB_SO_DEV)
	- rm -rfv $(DESTDIR)$(LIBDIR)/$(CLIB_SO_MAN)
	- rm -rfv $(DESTDIR)$(LIBDIR)/$(CLIB_SO_FULL)
	- rm -rfv $(DESTDIR)$(INCLUDE_DIR)/$(shell basename $(CLIB_HEADER))
	- rm -rfv $(DESTDIR)$(INCLUDE_DIR)/$(shell basename $(CLIB_PKG_CONFIG))
