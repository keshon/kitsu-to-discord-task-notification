REM @ECHO off

REM Release
SET ROOTPATH=.\release

@RD /S /Q %ROOTPATH%

MKDIR %ROOTPATH%\
COPY README.md %ROOTPATH%
COPY LICENSE %ROOTPATH%
COPY empty.conf.toml %ROOTPATH%\conf.toml

MKDIR %ROOTPATH%\tpl
MKDIR %ROOTPATH%\tpl\eng
MKDIR %ROOTPATH%\tpl\rus

COPY .\tpl\eng\author.tpl %ROOTPATH%\tpl\eng
COPY .\tpl\eng\description.tpl %ROOTPATH%\tpl\eng
COPY .\tpl\eng\footer.tpl %ROOTPATH%\tpl\eng
COPY .\tpl\eng\title.tpl %ROOTPATH%\tpl\eng
COPY .\tpl\eng\footer.tpl %ROOTPATH%\tpl\eng
COPY .\tpl\eng\title.tpl %ROOTPATH%\tpl\eng

COPY .\tpl\rus\author.tpl %ROOTPATH%\tpl\rus
COPY .\tpl\rus\description.tpl %ROOTPATH%\tpl\rus
COPY .\tpl\rus\footer.tpl %ROOTPATH%\tpl\rus
COPY .\tpl\rus\title.tpl %ROOTPATH%\tpl\rus
COPY .\tpl\rus\footer.tpl %ROOTPATH%\tpl\rus
COPY .\tpl\rus\title.tpl %ROOTPATH%\tpl\rus

REM Docker
SET ROOTPATH=.\deploy\data

@RD /S /Q %ROOTPATH%

MKDIR %ROOTPATH%\
COPY empty.conf.toml %ROOTPATH%\conf.toml
COPY NUL %ROOTPATH%\sqlite.db

MKDIR %ROOTPATH%\tpl
MKDIR %ROOTPATH%\tpl\eng
MKDIR %ROOTPATH%\tpl\rus

COPY .\tpl\eng\author.tpl %ROOTPATH%\tpl\eng
COPY .\tpl\eng\description.tpl %ROOTPATH%\tpl\eng
COPY .\tpl\eng\footer.tpl %ROOTPATH%\tpl\eng
COPY .\tpl\eng\title.tpl %ROOTPATH%\tpl\eng

COPY .\tpl\rus\author.tpl %ROOTPATH%\tpl\rus
COPY .\tpl\rus\description.tpl %ROOTPATH%\tpl\rus
COPY .\tpl\rus\footer.tpl %ROOTPATH%\tpl\rus
COPY .\tpl\rus\title.tpl %ROOTPATH%\tpl\rus

go build -ldflags "-s -w" -o .\release\app.exe src/main.go