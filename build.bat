@echo off
echo IDI_ICON1 ICON "./server/html/static/favicon.ico" > main.rc
windres -o main.syso main.rc
go build -ldflags="-s -w -H windowsgui" .
del main.rc
del main.syso
upx .\PushMeClient.exe