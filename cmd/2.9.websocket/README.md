### websocket example

Run the go server:

``` go
go run cmd/websocket/main.go
```

then make a echo request:

``` bash
websocat ws://localhost:5432/echo
hey
```

to make a chat-like message

```bash
# make this from two processes at the same time to simulate chat
websockat ws://localhost:5432/chat
het
```

Ref: https://github.com/gorilla/websocket
https://github.com/vi/websocat
