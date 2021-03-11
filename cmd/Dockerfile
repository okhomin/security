#
# STEP 1 build an executable binary
#
FROM golang:1.16 AS builder
WORKDIR /source
COPY ./cmd /source/cmd
COPY ./internal /source/internal
COPY ./go.mod /source/go.mod
RUN CGO_ENABLED=0 GOOS=linux go build -o server /source/cmd/main.go

#
# STEP 2 build an alpine image with the binary
#
FROM alpine:latest
WORKDIR /root
COPY --from=builder /source/server .
CMD ["./server"]