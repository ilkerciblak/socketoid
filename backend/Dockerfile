FROM golang:1.25-alpine AS dev

WORKDIR ./app

COPY go.mod ./
RUN go mod download


COPY . .
RUN go build -o /bin/app/runner

EXPOSE 8000

CMD ["/bin/app/runner"]
