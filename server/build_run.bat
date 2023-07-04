@echo off
set GOARCH=amd64
echo Building...
go build -o .\build\evilqr-server.exe -mod=vendor && echo Launching... && .\build\evilqr-server.exe -d .\templates
