FROM golang:1.23-alpine AS builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change

COPY go.mod go.sum ./
RUN go mod download && go mod verify


COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /manul ./cmd/manul/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /manul /manul

CMD ["/manul"]
