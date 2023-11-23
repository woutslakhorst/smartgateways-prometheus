FROM golang:1.21-alpine as builder

ENV GOPATH /

RUN mkdir /opt/kp && cd /opt/kp
COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /opt/kp/kp

# alpine 3.10.3
FROM alpine:3.11
COPY --from=builder /opt/kp/kp /usr/bin/kp
ENTRYPOINT ["/usr/bin/kp"]
