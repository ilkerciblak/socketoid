build:
	@go build -C backend/ -o ./bin/runner
run: build
	@./backend/bin/runner
