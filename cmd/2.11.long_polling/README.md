### Simple long polling example

Run the go server:

``` go
go run cmd/short_polling/main.go
```

then make a request:

``` bash
# submit a job
curl "localhost:5432/submit"

# check for progress
curl 'localhost:5432/checkstatus?job_id=1738791277'
```

Note: this will persist the tcp connection until the job is finished.
