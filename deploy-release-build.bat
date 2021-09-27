REM @ECHO off

REM Release
@RD /S /Q .\release

SET ROOTPATH=.\release

MKDIR %ROOTPATH%\
COPY README.md %ROOTPATH%
COPY LICENSE %ROOTPATH%
COPY empty.conf.toml %ROOTPATH%\conf.toml

MKDIR %ROOTPATH%\tpl
COPY .\tpl\author.tpl %ROOTPATH%\tpl
COPY .\tpl\description.tpl %ROOTPATH%\tpl
COPY .\tpl\footer.tpl %ROOTPATH%\tpl
COPY .\tpl\title.tpl %ROOTPATH%\tpl

REM Docker
@RD /S /Q .\deploy\data

SET ROOTPATH=.\deploy\data\

MKDIR %ROOTPATH%\
COPY empty.conf.toml %ROOTPATH%\conf.toml
COPY NUL %ROOTPATH%\sqlite.db

MKDIR %ROOTPATH%\tpl
COPY .\tpl\author.tpl %ROOTPATH%\tpl
COPY .\tpl\description.tpl %ROOTPATH%\tpl
COPY .\tpl\footer.tpl %ROOTPATH%\tpl
COPY .\tpl\title.tpl %ROOTPATH%\tpl

go build -ldflags "-s -w" -o .\release\app.exe src/main.go