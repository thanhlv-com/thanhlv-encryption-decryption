# syntax=docker/dockerfile:1

FROM golang:1.21-alpine AS builder
WORKDIR /src

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -trimpath -ldflags="-s -w" -o /out/thanhlv-ed main.go

FROM alpine:3.20
RUN apk add --no-cache ca-certificates

COPY --from=builder /out/thanhlv-ed /usr/local/bin/thanhlv-ed

ENTRYPOINT ["/usr/local/bin/thanhlv-ed"]
CMD ["--help"]
