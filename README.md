<div align="center">
    <img src = "yakv.png">
    <br>
</div>

yakv (*yak-v. (originally intended to be "yet-another-key-value store")*) is a simple, in-memory, concurrency-safe key-value store for hobbyists.

yakv provides persistence by appending transactions to a transaction log and restoring data from the transaction log on startup.

yakv is designed with simplicity as the main purpose and has *almost zero external dependencies*.

## Installation:

Spin up a server using:

- **Docker**:
    ```
    git clone https://github.com/burntcarrot/yakv
    cd yakv
    docker build --tag yakv .
    docker run -p 8080:8080 yakv /bin/sh -c "/yakv -host 0.0.0.0"
    ```
- **Install from source:**

    You can run directly from the source files:
    ```
    git clone https://github.com/burntcarrot/yakv
    cd yakv
    go run main.go -port 8080
    ```
    Or, you can build the binary on your own:
    ```
    git clone https://github.com/burntcarrot/yakv
    cd yakv
    go build
    ```

## Methods:

yakv exposes a HTTP/HTTPS API and provides 3 methods to deal with data:

- **GET**:
    - On a HTTPS server without certificate:
    ```
    curl -X GET --header "Content-Type: application/json" -d '{"key": "yakv"}' http://0.0.0.0:8080/yakv/v0/get --insecure
    ```
    - On a HTTP server:
    ```
    curl -X GET --header "Content-Type: application/json" -d '{"key": "yakv"}' http://0.0.0.0:8080/yakv/v0/get
    ```
- **PUT**:
    - On a HTTPS server without certificate:
    ```
    curl -X PUT --header "Content-Type: application/json" -d '{"key": "yakv", "value": "Hello, yakv!"}' http://0.0.0.0:8080/yakv/v0/put --insecure
    ```
    - On a HTTP server:
    ```
    curl -X PUT --header "Content-Type: application/json" -d '{"key": "yakv", "value": "Hello, yakv!"}' http://0.0.0.0:8080/yakv/v0/put
    ```
- **DELETE**:
    - On a HTTPS server without certificate:
    ```
    curl -X DELETE --header "Content-Type: application/json" -d '{"key": "yakv"}' http://0.0.0.0:8080/yakv/v0/delete --insecure
    ```
    - On a HTTP server:
    ```
    curl -X DELETE --header "Content-Type: application/json" -d '{"key": "yakv"}' http://0.0.0.0:8080/yakv/v0/delete
    ```

yakv currently accepts request bodies in the form of JSON.

## Options:

Here are the list of options or the command line flags provided by yakv:

```
yakv [OPTIONS]

OPTIONS:
    - port
        Port number for starting yakv.
    - host
        Host address for starting yakv.

    -secure
        Enable TLS-encrypted connection.
    - cert
        Filename for certificate.
    - key
        Filename for private key.

    -filename
        Filename for transaction log.
```

## Transaction Log:

All of the transactions are backed up in a transaction log, which are automatically loaded up by yakv on start-up.

## Security:

yakv provides a TLS-encrypted HTTPS connection using the `-secure` flag.

A certificate and a matching private key for the server must be provided through the `-cert` and `-key` flags respectively.

If the flags are not provided, yakv assumes the certificate and key to be named as `cert.pem` and `key.pem` in the current directory.

Example:

**On Docker:**

```
docker run -p 8080:8080 yakv /bin/sh -c "/yakv -host 0.0.0.0 -secure tls"
```

**Locally:**
- From source code:
    ```
    go run main.go -port 8080 -secure tls
    ```
- From binary:
    ```
    ./yakv -port 8080 -secure tls
    ```

## Attributions:

The yak vector is provided by [OpenClipart/FreeSVG](https://freesvg.org/vector-drawing-of-a-yak) under the [Public Domain](https://creativecommons.org/licenses/publicdomain/).
