# GPT Stream Server

This project serves as an API proxy for the OpenAI GPT API, allowing for the use of stream mode. It utilizes the `yao-chatgpt` interface management tool to manage API keys and content.

## Installation

To install this project, clone the `chatgpt-web` repository as front-end:

```sh
git clone https://github.com/wwsheng009/chatgpt-web

pnpm i && pnpm build

```

Install the `yao` application-engine:

download the yao-0.10.3 release from github actions. install the yao bin to your local machine.

https://github.com/YaoApp/yao/actions/workflows/release-linux.yml

Then, install the `yao-chatgpt` application for backend content management and API key management:

```sh

git clone https://github.com/wwsheng009/yao-chatgpt

cd yao-chatgpt && yao start
```

At last,install the `gpt-stream-server`

```sh
git clone https://github.com/wwsheng009/gpt-stream-server

go mod tidy

/bin/bash build
```

## Usage

update the OpenAI key use the `yao-chatgpt` application. login the yao application and change the api key refer to the `yao-chatgpt` document

https://github.com/wwsheng009/yao-chatgpt/blob/main/readme.md

Once all repositories are installed, start the `chatgpt-web` frontend with the modified version that uses `fetch` as the request library to support stream mode:

```sh
npm start
```

You can then interact with the OpenAI GPT API through the `chatgpt-web` interface management tool.

## Contributing

If you would like to contribute to this project, please fork the repository and submit a pull request with your changes.

## License

This project is licensed under the MIT license.
