FROM golang:1.24-alpine

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
WORKDIR /app

COPY ./configs/certs/russian_trusted_root_ca_pem.crt /usr/local/share/ca-certificates/
COPY ./configs/certs/russian_trusted_sub_ca_2024_pem.crt /usr/local/share/ca-certificates/
COPY ./configs/certs/russian_trusted_sub_ca_pem.crt /usr/local/share/ca-certificates/

RUN update-ca-certificates

COPY go.mod go.sum ./
RUN go mod download
COPY . .

WORKDIR /app/
RUN go get -v ./... \
  && go install -v ./... \
  && go build -v -o tgbot

ENTRYPOINT [ "/app/tgbot" ]
