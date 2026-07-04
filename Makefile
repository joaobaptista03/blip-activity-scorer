BINARY := blip-activity-scorer

.PHONY: build test run bench lint fmt githooks clean

build:
	go build -o $(BINARY) .

test:
	go test -v ./...

run: build
	./$(BINARY)

bench:
	go test -bench=. -benchmem ./...

lint:
	go vet ./...
	gofmt -l .

fmt:
	gofmt -w .

githooks:
	git config core.hooksPath .githooks
	chmod +x .githooks/pre-commit

clean:
	rm -f $(BINARY) ranking_full.csv
