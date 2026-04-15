#!/usr/bin/env bash
go mod tidy && go mod vendor && go build cmd/main.go && rm main

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o minigo cmd/main.go

ssh ubuntu@152.136.253.39 "rm /home/ubuntu/minigo"

scp minigo ubuntu@152.136.253.39:/home/ubuntu/minigo
scp ./config/release.yml ubuntu@152.136.253.39:/home/ubuntu/config/release.yml
scp ./config/release.yml ubuntu@152.136.253.39:/home/ubuntu/config/production.yml

rm minigo