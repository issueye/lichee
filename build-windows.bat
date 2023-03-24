set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=windows

go mod vendor

go build -ldflags "-w -s" -o go-plugins.exe main.go

upx go-plugins.exe