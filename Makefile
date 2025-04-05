build:
	go build -o backend-exec main.go

run: build
	./backend-exec

watch:
	reflex -s -r '\.go$$' make run
