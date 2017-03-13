.PHONY: clean dist

BINARY:=shex
WINBINARY:=$(BINARY).exe
WINTARASSETS:=open_console.cmd
WINZIPASSETS:=assets/open_console.cmd

VERSION_FILE=VERSION
VERSION:=$(shell cat ${VERSION_FILE})
#VERSION:=`git describe --VERSIONs`
#VERSION=0.0.1-alpha

LDFLAGS=-ldflags "-X main.buildVersion=$(VERSION)"

all:
	make release
	make dist

release:
	mkdir -p dist/linux/i386 && GOOS=linux GOARCH=386 go build $(LDFLAGS) -o dist/linux/i386/$(BINARY)
	mkdir -p dist/linux/amd64 && GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/linux/amd64/$(BINARY)
	mkdir -p dist/windows/i386 && GOOS=windows GOARCH=386 go build $(LDFLAGS) -o dist/windows/i386/$(WINBINARY)
	mkdir -p dist/windows/amd64 && GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/windows/amd64/$(WINBINARY)

dist:
	# linux 32/64
	mkdir -p dist/release
	tar -cvzf dist/release/$(BINARY)-linux-amd64-$(VERSION).tar.gz -C dist/linux/amd64 $(BINARY)
	tar -cvzf dist/release/$(BINARY)-linux-i386-$(VERSION).tar.gz -C dist/linux/i386 $(BINARY)

	# WIN32
	tar -cvf dist/release/$(BINARY)-win32-$(VERSION).tar -C dist/windows/i386 $(WINBINARY)
	tar -uvf dist/release/$(BINARY)-win32-$(VERSION).tar -C assets $(WINTARASSETS)
	gzip -f dist/release/$(BINARY)-win32-$(VERSION).tar

	# WIN64
	tar -cvf dist/release/$(BINARY)-win64-$(VERSION).tar -C dist/windows/amd64 $(WINBINARY)
	tar -uvf dist/release/$(BINARY)-win64-$(VERSION).tar -C assets $(WINTARASSETS)
	gzip -f dist/release/$(BINARY)-win64-$(VERSION).tar

	# WIN 32/64 ZIP
	zip -j dist/release/$(BINARY)-win32-$(VERSION).zip dist/windows/i386/$(WINBINARY) $(WINZIPASSETS)
	zip -j dist/release/$(BINARY)-win64-$(VERSION).zip dist/windows/amd64/$(WINBINARY) $(WINZIPASSETS)

#push_release: dist
#	API_JSON:=$(printf '{"VERSION_name": "v%s","target_commitish": "master","name": "v%s","body": "Release of version %s","draft": false,"prerelease": false}')
#	echo API_JSON

clean:
	rm -rf dist
