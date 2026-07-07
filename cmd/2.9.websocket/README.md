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

Note: It is highly recommended the lib of choice to be highly compatible with the web socket protocol [RFC](https://datatracker.ietf.org/doc/html/rfc6455)

There is test suite for this compliance: https://github.com/crossbario/autobahn-testsuite

Ref: https://github.com/gorilla/websocket
https://github.com/vi/websocat
