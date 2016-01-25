PREFIX=/usr/local
VERSION=$(shell git describe)
GOBUILD=go build -ldflags "-X main.Version=${VERSION}"

all: async

async: $(shell find . -name '*.go')
	${GOBUILD} -o $@

install: async
	mkdir -p $(DESTDIR)${PREFIX}/bin/
	install async ${DESTDIR}${PREFIX}/bin/

.PHONY: clean
clean:
	rm -f async
