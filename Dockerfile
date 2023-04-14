# -----------------------------------------------------------------------------
# Build...
FROM golang:1.20 -alpine3.14 AS build

WORKDIR /hercules

COPY go.mod go.sum main.go Makefile ./
COPY internal internal
COPY cmd cmd
COPY types types
COPY pkg pkg
RUN apk --no-cache add make git gcc libc-dev curl ca-certificates && make build

# -----------------------------------------------------------------------------
# # Image...
# FROM alpine:3.14.0

# COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# COPY --from=build /hercules/execs/hercules /bin/hercules

ENTRYPOINT [ "/bin/hercules" ]
