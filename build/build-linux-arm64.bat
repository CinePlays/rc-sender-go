cd ../src

set GOOS=linux
set GOARCH=arm64

go build -o ./build/out/rc-sender-linux-arm64 ./src/cmd/sender/