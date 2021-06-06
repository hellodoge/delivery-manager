FROM golang:1.16

COPY ./ ./

ENV GOPATH=/

RUN go mod download
RUN go build -o delivery-manager ./cmd/main.go

CMD ["./delivery-manager"]