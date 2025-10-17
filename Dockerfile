FROM golang:1.25-alpine AS builder

COPY . /src
WORKDIR /src
RUN go get . && go build .

FROM alpine:3.22
RUN apk add --no-cache ca-certificates git curl

WORKDIR /app
COPY --from=builder /src/herlighet /app/herlighet

CMD ["/app/herlighet"]
