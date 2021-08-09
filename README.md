<div align="center">
    <img src = "static/yakv.png">
    <br><br>
    <a href="http://makeapullrequest.com"><img src ="https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square"></a>
    <a href = "https://github.com/burntcarrot/yakv/actions?workflow=Tests"><img src = "https://github.com/burntcarrot/yakv/workflows/Tests/badge.svg"></a>
    <a href="https://pkg.go.dev/github.com/burntcarrot/yakv"><img src="https://godoc.org/github.com/burntcarrot/yakv?status.svg" /></a>
    <br><br>
    <img src = "static/term-preview-yakv.svg">
    <br><br>
</div>

yakv (*yak-v. (originally intended to be "yet-another-key-value store")*) is a simple, in-memory, concurrency-safe key-value store for hobbyists.

yakv provides persistence by appending transactions to a transaction log and restoring data from the transaction log on startup.

yakv is designed with simplicity as the main purpose and has *almost zero external dependencies*.

<h2>Table of Contents:</h2>

- **[Installation](#installation)**
- [Methods](#methods)
- [Options](#options)
- [Transaction Log](#transaction-log)
- [Security](#security)
- [Benchmarks](#benchmarks)
- [FAQ](#faq)
- **[Contributing Guide](#contributing-guide)**
- [Attributions](#attributions)

## Installation

Install using:

- **One-Script Installation (Linux):**
    ```
    curl https://raw.githubusercontent.com/burntcarrot/yakv/main/install.sh | bash
    ```

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

## Methods

yakv exposes a HTTP/HTTPS API and provides 3 methods to deal with data:

- **GET**:
    - On a HTTPS server without certificate:
    ```
    curl -X GET --header "Content-Type: application/json" -d '{"key": "yakv"}' https://0.0.0.0:8080/yakv/v0/get --insecure
    ```
    - On a HTTP server:
    ```
    curl -X GET --header "Content-Type: application/json" -d '{"key": "yakv"}' http://0.0.0.0:8080/yakv/v0/get
    ```
- **PUT**:
    - On a HTTPS server without certificate:
    ```
    curl -X PUT --header "Content-Type: application/json" -d '{"key": "yakv", "value": "Hello, yakv!"}' https://0.0.0.0:8080/yakv/v0/put --insecure
    ```
    - On a HTTP server:
    ```
    curl -X PUT --header "Content-Type: application/json" -d '{"key": "yakv", "value": "Hello, yakv!"}' http://0.0.0.0:8080/yakv/v0/put
    ```
- **DELETE**:
    - On a HTTPS server without certificate:
    ```
    curl -X DELETE --header "Content-Type: application/json" -d '{"key": "yakv"}' https://0.0.0.0:8080/yakv/v0/delete --insecure
    ```
    - On a HTTP server:
    ```
    curl -X DELETE --header "Content-Type: application/json" -d '{"key": "yakv"}' http://0.0.0.0:8080/yakv/v0/delete
    ```

yakv currently accepts request bodies in the form of JSON.

## Options

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

## Transaction Log

All of the transactions are backed up in a transaction log, which are automatically loaded up by yakv on start-up.

## Security

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
    go run main.go -port 8080 -secure
    ```
- From binary:
    ```
    ./yakv -port 8080 -secure
    ```

## Benchmarks

Benchmarks are done using [vegeta](https://github.com/tsenart/vegeta).

All of the benchmarks are performed under these device specifications:

```
Processor:
Intel(R) Core(TM) i5-8265U CPU @ 1.60GHz   1.80 GHz

Installed RAM:
8.00 GB (7.85 GB usable)
```

> **NOTE: rate is set manually. This does not denote the maximum number of requests yakv can handle.**

<h4>5,000 PUT requests in 100 seconds (rate = 50 requests/second):</h4>

*Available RAM while performing benchmark: `3.8 GB`*

```
Requests      [total, rate, throughput]         5000, 50.01, 50.01
Duration      [total, attack, wait]             1m40s, 1m40s, 1.13ms
Latencies     [min, mean, 50, 90, 95, 99, max]  467.4µs, 2.687ms, 2.949ms, 3.579ms, 3.652ms, 3.723ms, 14.932ms
Bytes In      [total, mean]                     0, 0.00
Bytes Out     [total, mean]                     185000, 37.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      201:5000
Error Set:
```

![PUT-Benchmark](benchmarks/put-requests-100s.png)

## FAQ:

#### Why a database-based transaction log isn't available?

`yakv` was designed with simplicity as the main purpose, although this doesn't mean resistence to addition of new features, the addition of database-based transaction log is scheduled for future releases.

## Contributing Guide

Read the contributing guide [here.](https://github.com/burntcarrot/yakv/blob/main/CONTRIBUTING.md)

## Attributions

The yak vector is provided by [OpenClipart/FreeSVG](https://freesvg.org/vector-drawing-of-a-yak) under the [Public Domain](https://creativecommons.org/licenses/publicdomain/).
