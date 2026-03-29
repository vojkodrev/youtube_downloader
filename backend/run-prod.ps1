go build -o backend.exe .
$env:GIN_MODE = "release"
./backend.exe
