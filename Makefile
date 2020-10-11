mrtg_exporter: $(shell find . -type f -name *.go)
	CGO_ENABLED=0 go build \
	  -o "$@" ./cmd/mrtg_exporter
	@if ! ldd "$@" 2> /dev/null; then echo "OK: not a dynamic executable!"; fi

.PHONY: clean
clean:
	rm -Rfv mrtg_exporter

.PHONY: cl
cl:
	find . -type f -name *.go | xargs wc -l

.PHONY: test
test:
	go test ./...
