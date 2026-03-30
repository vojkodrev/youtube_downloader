go build -o backend.exe .
powershell -NoProfile -Command { $env:GIN_MODE = "release"; ./backend.exe }
