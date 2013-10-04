To run:

```
git submodule update --init # https://github.com/nfnt/resize
go build server.go
./server
```

add a PNG or JPEG image to ./images and visit:

```
http://localhost:8080/WIDTH/HEIGHT/FILENAME
```

Requires Go 1.1.2
