@echo off
set GOARCH=amd64
echo Building...
go build -o .\build\qrswap-server.exe -mod=vendor && echo Launching... && .\build\qrswap-server.exe -d .\templates
