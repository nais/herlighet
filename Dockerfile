FROM golang:1.15-alpine as builder

COPY . /src
WORKDIR /src
RUN go get . && go build .

FROM alpine:3.10
RUN apk add --no-cache ca-certificates git curl

WORKDIR /app
COPY --from=builder /src/herlighet /app/herlighet

CMD ["/app/herlighet"]
