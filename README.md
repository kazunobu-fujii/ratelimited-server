# ratelimited-server

rate limited server

# Usage

```
$ go get github.com/kazunobu-fujii/ratelimited-server/cmd/ratelimited-server
```

```
$ ratelimited-server -h
Usage of ratelimited-server:
  -c int
        chunk size (default 8)
  -d int
        bucket depth (default 1)
  -o int
        timeout (s) (default 30)
  -p string
        response body file path (default "response.dat")
  -r int
        request rate limit (per second) (default 1)
  -s string
        listening server (default "localhost:8080")
  -t string
        content type (default "text/html; charset=UTF-8")
  -w int
        chunk delay (ms) (default 10)
```
