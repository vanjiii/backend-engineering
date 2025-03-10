### server sent events example

Run the go server:

``` go
go run cmd/server_sent_events/main.go
```

Then in chrome:
- open 'http://localhost:5432/'
- open console with F12
- under 'Console' create event source with: `const sse = new EventSource("http://localhost:5432/stream")`
- under 'Network' locate the '/stream' path and open 'EventStream' tab
