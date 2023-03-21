

bindata:
	rm -rf .tmp/data
	mkdir -p .tmp/data
	git clone -b go-axios https://github.com/wwsheng009/chatgpt-web.git .tmp/chatgpt-web

	rm -rf .tmp/chatgpt-web/.git
	rm -rf .tmp/chatgpt-web/.gitignore
	rm -rf .tmp/chatgpt-web/LICENSE
	rm -rf .tmp/chatgpt-web/README.md

	cd .tmp/chatgpt-web  && pnpm install && pnpm run build

	cp -r .tmp/chatgpt-web/dist .tmp/data/web


	go-bindata -fs -pkg ui -o ui/bindata.go -prefix ".tmp/data/" .tmp/data/...

release:
	export CGO_ENABLED=0
	go build