build:
	@go build -C backend/ -o ./bin/runner
run: build
	@./backend/bin/runner

run-frontend:
	@cd frontend/ &&\
	npm run dev
