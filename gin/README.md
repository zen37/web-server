https://golang.org/doc/tutorial/web-service-gin

curl commands
curl http://localhost:8080/albums

curl http://localhost:8080/albums \
    --header "Content-Type: application/json" \
    --request "GET"

curl http://localhost:8080/albums \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"id": "4","title": "The Wall","artist": "Pink Floyd","price": 49.99}'

 curl http://localhost:8080/albums/2 