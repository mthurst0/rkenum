.PHONY: all
all: build

.PHONY: build
build:
	./rk-go-bld

.PHONY: install
install: build
	cp -v ./rkenum ~/bin
