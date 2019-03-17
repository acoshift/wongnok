FROM golang:1.12.1 as stage

WORKDIR /src

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.moonrhythm.io

ADD . .

RUN go build -o wongnok -ldflags '-w -s' main.go

# ------

FROM alpine:3.9

RUN apk add ca-certificates tzdata

WORKDIR /app

COPY --from=stage /src/wongnok .

EXPOSE 8080
ENTRYPOINT /app/wongnok
