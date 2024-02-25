.PHONY: run
run: templ.generate
	go run *.go

.PHONY: templ.install
templ.install:
	go install github.com/a-h/templ/cmd/templ@latest

.PHONY: templ.generate
templ.generate:
	templ generate


