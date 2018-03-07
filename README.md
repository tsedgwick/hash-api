# hash-api 

This small golang app takes an input (password) and encodes it with SHA512.  

### [cli](cmd/cli/) and [server](cmd/server/) 

    go get github.com/tsedgwick/hash-api

It has two entrypoints into functionality:

CLI [cli](cmd/cli/) covers:

* The basic SHA512 encoding of a password via command line

    go run cmd/cli/main.go -password=angryMonkey

Server [server](cmd/server/) covers:

* Hash and Encode Passwords over HTTP (returns the encoded password back to the user)
    curl -X POST \
        http://localhost:8080/v1/hash \
        -H 'cache-control: no-cache' \
        -H 'content-type: application/x-www-form-urlencoded' \
        -d password=angryMonkey

* Hash End-Point Returns Identifier (returns a key to retrieve the encoded value from the GET below)
    curl -X POST \
        http://localhost:8080/v2/hash \
        -H 'cache-control: no-cache' \
        -H 'content-type: application/x-www-form-urlencoded' \
        -d password=angryMonkey

* GET a Hashed Password (using the key provided from above, retrieve the encoded value)
    curl -X GET \
        http://localhost:8080/v3/hash/7887

* Graceful Shutdown

    curl -X POST \
        http://localhost:8080/shutdown

