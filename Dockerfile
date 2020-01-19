FROM golang:alpine3.11 as builder
WORKDIR /go/src/github.com/spawn2kill/bme280-exporter
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN go build -ldflags="-w -s" -tags netgo -a -o bme280 .

RUN addgroup -g 500 devices && adduser -H -D -g '' -G devices pi

FROM scratch
# Import from builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/src/github.com/spawn2kill/bme280-exporter/bme280 /app/
ENTRYPOINT ["/app/bme280"]