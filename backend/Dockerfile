FROM golang:1.24 AS builder

WORKDIR /src
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

RUN CGO_ENABLED=1 GOOS=linux go build -o /app -a -ldflags '-linkmode external -extldflags -static' -tags timetzdata -tags netgo nostalgia/cmd/web
RUN chmod +x /app

FROM alpine:latest
RUN apk update && apk add --no-cache ffmpeg
COPY --from=builder /app /app
COPY migrations ./migrations
EXPOSE 8080

ENTRYPOINT ["/app"]
