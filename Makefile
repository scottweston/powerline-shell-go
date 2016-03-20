all: osx linux windows

options = -ldflags "-X main.build=`date -u +%Y%m%d.%H%M%S`"

linux:
	GOOS=linux GOARCH=amd64 go build $(options) -o ./build/linux/x86-64/powerline-shell-go
	GOOS=linux GOARCH=386 go build $(options) -o ./build/linux/x86/powerline-shell-go
	GOOS=linux GOARCH=arm go build $(options) -o ./build/linux/arm/powerline-shell-go

osx:
	GOOS=darwin GOARCH=amd64 go build $(options) -o ./build/osx/x86-64/powerline-shell-go

windows:
	GOOS=windows GOARCH=amd64 go build $(options) -o ./build/windows/x86-64/powerline-shell-go
	GOOS=windows GOARCH=386 go build $(options) -o ./build/windows/x86/powerline-shell-go

clean:
	go clean
	rm -rf build

install:
	go install $(options)
