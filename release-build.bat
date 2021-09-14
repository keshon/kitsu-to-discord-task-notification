go build -ldflags "-s -w" -o .\release\app.exe src/main.go

COPY README.md .\release
COPY LICENSE .\release
COPY empty.conf.toml .\release\conf.toml
