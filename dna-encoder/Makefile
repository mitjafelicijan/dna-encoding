run:
	find . \( -name "*.go" \) | entr -r go run *.go

build-linux:
	mkdir -p ../dist/linux_x86_64
	CGO_ENABLED=0 GOOS=linux go build -v -a \
		-ldflags '-extldflags "-static"' \
		-o ../dist/linux_x86_64/dnae-encode

build-darwin:
	mkdir -p ../dist/darwin_x86_64
	CGO_ENABLED=0 GOOS=darwin go build -v -a \
		-ldflags '-extldflags "-static"' \
		-o ../dist/darwin_x86_64/dnae-encode
