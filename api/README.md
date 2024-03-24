# Go API sample

goのAPIのsampleです

## Quick Start
```
$ go run .
```

get albums
```
$ curl http://localhost:8080/albums
```

get specific album
```
$ curl http://localhost:8080/albums/2
```

post album
```
$ curl http://localhost:8080/albums \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"id": "4","title": "The Modern Sound of Betty Carter","artist": "Betty Carter","price": 49.99}'
```





## using it
https://go.dev/doc/tutorial/web-service-gin

