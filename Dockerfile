FROM golang:1.25-alpine AS builder

WORKDIR /src
COPY go.mod go.sum /
RUN go mod download
COPY . .
RUN go build .

FROM alpine:3.22
RUN apk add --no-cache ca-certificates git curl

WORKDIR /app
COPY --from=builder /src/herlighet /app/herlighet

CMD ["/app/herlighet"]
