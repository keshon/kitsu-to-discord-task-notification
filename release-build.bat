go build -ldflags "-s -w" -o .\release\app.exe src/main.go

COPY README.md .\release
COPY LICENSE .\release
COPY empty.conf.toml .\release\conf.toml

mkdir .\release\tpl
COPY .\tpl\author.tpl .\release\tpl
COPY .\tpl\description.tpl .\release\tpl
COPY .\tpl\footer.tpl .\release\tpl
COPY .\tpl\title.tpl .\release\tpl
