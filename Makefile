BINARY_NAME = ctt
VERSION_NUMBER ?= $(shell git describe --tags --always | sed -E 's#^v##')

ARCH = x86_64
DESTDIR ?=
BINDIR ?= /usr/local/bin
SHAREDIR ?= /usr/local/share

all: clean test build

clean:
	go clean
	rm -f $(BINARY_NAME)

version_number:
	@echo $(VERSION_NUMBER)

test:
	go test -v -timeout=30s ./...

build:
	go build -trimpath -buildmode=pie -mod=readonly -modcacherw -v -ldflags "-linkmode external -extldflags \"${LDFLAGS}\" -X main.version=${VERSION_NUMBER}" -o $(BINARY_NAME) ./cmd

install:
	install -Dm755 $(BINARY_NAME) $(DESTDIR)$(BINDIR)/$(BINARY_NAME)
	cp .assets/ctt.desktop ${DESTDIR}${SHAREDIR}/applications/${BINARY_NAME}.desktop

uninstall:
	rm -f $(DESTDIR)$(BINDIR)/$(BINARY_NAME)
	rm -f ${DESTDIR}${SHAREDIR}/applications/${BINARY_NAME}.desktop
