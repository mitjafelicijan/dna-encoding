run:
	find . \( -name "*.go" \) | entr -r go run *.go

pprof:
	-rm *.pdf
	-rm *.pprof
	go build
	./dna-encoder -i ../sample-data/hello -o ../encoded-data/hello.fa -profile
	go tool pprof -pdf ./dna-encoder cpu.pprof > cpu.pdf

build:
	mkdir -p ../dist
	CGO_ENABLED=0 GOOS=linux go build -v -a \
		-ldflags '-extldflags "-static"' \
		-o ../dist/dnae-encode