set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=linux

go mod tidy
go mod vendor

go build -ldflags "-w -s" -o go-plugins main.go

upx go-plugins