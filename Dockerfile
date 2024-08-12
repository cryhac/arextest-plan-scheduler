FROM golang:1.18 AS builder

WORKDIR /app

ENV GOPROXY https://goproxy.cn/
ENV GO111MODULE on

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest

ENV APP_ID=APP_ID
ENV TARGET_HOST=TARGET_HOST

COPY --from=builder /app/app /app

ENTRYPOINT ["/app"]