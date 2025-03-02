FROM golang:alpine

WORKDIR /src

COPY . ./
RUN  go mod download
RUN GOOS=linux GOARCH=amd64 GO111MODULE=on CGO_ENABLED=0 \
    go build -o service ./cmd/main.go

CMD ./service
