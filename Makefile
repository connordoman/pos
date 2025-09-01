build-linux:
	mkdir -p bin
	GOOS=linux GOARCH=arm64 go build -o bin/pos-linux main.go

clean-linux:
	rm -rf bin/pos-linux