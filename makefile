.PHONY: run
run: templ.generate
	go run *.go

rt:
	templ generate --watch --proxy="http://localhost:8080" --cmd="go run ."

.PHONY: templ.install
templ.install:
	go install github.com/a-h/templ/cmd/templ@latest

.PHONY: templ.generate
templ.generate:
	templ generate

lint.install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.2

lint.run:
	golangci-lint run -c .golangci.yml -v --fix

