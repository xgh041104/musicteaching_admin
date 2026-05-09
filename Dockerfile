FROM golang:1.23 AS builder

WORKDIR /app
COPY ./server_go/go.mod ./server_go/go.sum ./
RUN go mod download

COPY ./server_go ./
RUN cd cmd/server && GOOS=linux go build -o /app/server_go

FROM ubuntu:latest
EXPOSE 7890
WORKDIR /service

RUN apt update && \
    apt -y upgrade && \
    apt-get install -y ffmpeg openjdk-8-jdk libreoffice p7zip-full && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server_go /service/server_go
COPY ./server_go/config/ /service/config/
COPY ./server_go/Resources /service/Resources

CMD ["/service/server_go"]
