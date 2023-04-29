@echo off

rmdir /s /q .tmp
mkdir .tmp\data
git clone --depth 1 https://github.com/wwsheng009/chatgpt-web.git .tmp\chatgpt-web

rmdir /s /q .tmp\chatgpt-web\.git
del /f .tmp\chatgpt-web\.gitignore
del /f .tmp\chatgpt-web\LICENSE
del /f .tmp\chatgpt-web\README.md

cd .tmp\chatgpt-web && pnpm install --no-frozen-lockfile && pnpm run build && cd ../../

xcopy /e /y .tmp\chatgpt-web\dist\ .tmp\data\web\

go-bindata -fs -pkg ui -o ui/bindata.go -prefix ".tmp/data/" .tmp/data/...


set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=windows

go build

set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=linux
go build -o gpt_stream_server