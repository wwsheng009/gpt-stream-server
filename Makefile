bindata:
	rm -rf .tmp/data
	mkdir -p .tmp/data
	
	git clone --depth 1 -b main https://github.com/wwsheng009/chatgpt-web.git .tmp/chatgpt-web

	cd .tmp/chatgpt-web && git checkout go-stream-upgrade
	cd ../../
	rm -rf .tmp/chatgpt-web/.git
	rm -rf .tmp/chatgpt-web/.gitignore
	rm -rf .tmp/chatgpt-web/LICENSE
	rm -rf .tmp/chatgpt-web/README.md

	cd .tmp/chatgpt-web  && pnpm install --no-frozen-lockfile && pnpm run build

	cp -r .tmp/chatgpt-web/dist .tmp/data/web

	go-bindata -fs -pkg ui -o ui/bindata.go -prefix ".tmp/data/" .tmp/data/...

release:
	CGO_ENABLED=0 go build
	
linux-windows:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build