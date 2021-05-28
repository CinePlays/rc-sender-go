cd ../src

set GOOS=linux
set GOARCH=arm

go build -o ./build/out/rc-sender-linux-arm32 ./src/cmd/sender/