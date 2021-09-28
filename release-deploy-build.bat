REM @ECHO off

REM Release
SET ROOTPATH=.\release

@RD /S /Q %ROOTPATH%

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
SET ROOTPATH=.\deploy\data

@RD /S /Q %ROOTPATH%

MKDIR %ROOTPATH%\
COPY empty.conf.toml %ROOTPATH%\conf.toml
COPY NUL %ROOTPATH%\sqlite.db

MKDIR %ROOTPATH%\tpl
COPY .\tpl\author.tpl %ROOTPATH%\tpl
COPY .\tpl\description.tpl %ROOTPATH%\tpl
COPY .\tpl\footer.tpl %ROOTPATH%\tpl
COPY .\tpl\title.tpl %ROOTPATH%\tpl

go build -ldflags "-s -w" -o .\release\app.exe src/main.go