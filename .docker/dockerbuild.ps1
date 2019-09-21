param (
    [Parameter(Mandatory=$true)][string]$version
)

$env:GOOS = "linux"
go build -o webserver.o
$env:GOOS = "windows"
docker build -t "webserver:$version" ./