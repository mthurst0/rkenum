.PHONY: all
all: build

.PHONY: build
build:
	./rk-go-bld

.PHONY: install
install: build
	mkdir -p ~/bin
	cp -v ./rkenum ~/bin
