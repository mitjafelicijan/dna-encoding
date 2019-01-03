run:
	find . \( -name "*.go" \) | entr -r go run *.go

build:
	mkdir -p ../dist
	CGO_ENABLED=0 GOOS=linux go build -v -a \
		-ldflags '-extldflags "-static"' \
		-o ../dist/fasta-to-image