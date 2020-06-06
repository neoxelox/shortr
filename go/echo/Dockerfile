FROM golang:1.14.4 AS builder

ENV CGO_ENABLED=0 \
    GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

# Download dependencies before build in order to cache them
COPY go/echo/go.mod go/echo/go.sum /build/

RUN go mod download

COPY go/echo/ /build

RUN go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o shortr ./main.go

FROM alpine AS app

RUN apk add ca-certificates

WORKDIR /app

COPY static /static
COPY --from=builder /build/shortr /app

# App
EXPOSE 80

CMD [ "/app/shortr" ]