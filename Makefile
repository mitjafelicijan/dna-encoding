default:
	@echo "Use 'make linux' or 'make darwin'"

linux:
	cd dna-encoder && make build-linux
	cd fasta-to-image && make build-linux

darwin:
	cd dna-encoder && make build-darwin
	cd fasta-to-image && make build-darwin
