FROM golang:1.21 AS builder

WORKDIR /go/src/app

ENV GOPROXY=https://goproxy.io,direct

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
RUN CGO_ENABLED=0 go build -ldflags="-extldflags=-static" ./cmd/server

FROM busybox AS unzip

WORKDIR /clientportal.gw

COPY clientportal.gw.zip ./
RUN unzip clientportal.gw.zip

FROM amazoncorretto:21 AS runner

WORKDIR /app

COPY --from=builder /go/src/app/server .
COPY --from=unzip /clientportal.gw ./clientportal.gw

EXPOSE 8000

ENTRYPOINT ["./server"]
