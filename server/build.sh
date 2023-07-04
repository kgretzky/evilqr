#!/bin/bash
go build -ldflags="-s -w" -o ./build/evilqr-server -mod=vendor main.go
