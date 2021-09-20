FROM golang:1.17-alpine as build
COPY . /tmp/hal/
WORKDIR /tmp/hal
RUN go mod tidy
RUN go build -o /tmp/build/hal ./cmd/hal

FROM alpine:latest
COPY --from=build --chown=root:root /tmp/build/hal /usr/local/bin/hal
CMD ["/usr/local/bin/hal"]
