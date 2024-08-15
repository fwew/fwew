SOURCES = completer.go config.go fwew.go lib.go numbers.go numbers_test.go txt.go
TAG = build
OS = nix
ifeq ($(OS),nix)
CP = sudo cp
RM = sudo rm
BINDEST = /usr/local/bin
else ifeq ($(OS),termux)
CP = cp
RM = rm
BINDEST = /data/data/com.termux/files/usr/bin
endif
WGET := $(shell command -v wget 2> /dev/null)
CURL := $(shell command -v curl 2> /dev/null)
DOWNLOAD := $(if $(WGET),wget,$(if $(CURL),curl,echo "Error: wget or curl not found"; exit 1;))

fwew: format compile

all: test docker cross-compile

download:
	$(DOWNLOAD) -O .fwew/dictionary-v2.txt https://tirea.learnnavi.org/dictionarydata/dictionary-v2.txt

format:
	gofmt -w $(SOURCES)

compile: download
	go build -o bin/fwew ./...

cross-compile: download
	GOOS=darwin go build -o bin/mac/fwew ./...
	GOOS=linux go build -o bin/linux/fwew ./...
	GOOS=windows go build -o bin/windows/fwew.exe ./...

test: install
	go test -v -cover

docker: download
	docker build -t tirea/fwew:$(TAG) .
	docker run -it --rm tirea/fwew:$(TAG) -v -r test

copy:
	@test -n "$(BIN)" || (echo "Error: BIN variable not set. BIN must be set to one of the following:" ; /bin/ls bin | grep -v fwew ; exit 1)
	$(CP) bin/$(BIN)/fwew $(BINDEST)/
	cp -r .fwew ~/

install: fwew
	$(CP) bin/fwew $(BINDEST)/
	cp -r .fwew ~/

uninstall:
	$(RM) $(BINDEST)/fwew
	rm -rf ~/.fwew

clean:
	rm -f bin/*

