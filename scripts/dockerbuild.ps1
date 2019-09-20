$env:GOOS = "linux"
go build -o webserver.o
$env:GOOS = "windows"
docker build -t "webserver:1.0" ./