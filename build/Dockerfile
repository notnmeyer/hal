FROM golang:1.17-alpine as build

# RUN apk add --no-cache git make musl-dev go

WORKDIR /tmp/hal
COPY . .
RUN cd /tmp/hal && go mod tidy
RUN go build -o /tmp/build/hal ./cmd/hal

FROM alpine:latest
COPY --from=build --chown=root:root /tmp/build/hal /usr/local/bin/hal
CMD ["/usr/local/bin/hal"]
