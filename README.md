# telegrambotgo
Tiny client for Telegram Bot Api v2

See examples_test.go to learn how to use this client properly

## Getting the client

You can use the client in your applications either using go get 

```
go get https://github.com/xordi/telegrambotgo
```

or using [govendor](https://github.com/kardianos/govendor) from your application directory in order to get the latest stable commit

```
govendor fetch https://github.com/xordi/telegrambotgo@v0.1
```

## Currently implemented methods
- getUpdates
- sendMessage
- sendPhoto
- sendAudio
- sendDocument
- sendVideo
