# Distributed [![godoc](https://godoc.org/github.com/m3o/go-distributed?status.svg)](https://godoc.org/github.com/m3o/go-distributed) [![Go Report Card](https://goreportcard.com/badge/github.com/m3o/go-distributed)](https://goreportcard.com/report/github.com/m3o/go-distributed) [![Apache 2.0 License](https://img.shields.io/github/license/m3o/go-distributed)](https://github.com/m3o/go-distributed/blob/master/LICENSE)

A federated community API

## Overview

Build federated communities using the same API. Distributed provides a single reusable API for building multiple communities 
on multiple platforms. Whether its for work, social, or anything else, quickly bring up the backend API for it and spin up 
your own frontend experience.

## Usage

Distributed is built as a single program in Go.

## Server

Download and install

```sh
go get github.com/m3o/go-distributed/cmd/distributed
```

Run the API

```
distributed
```

Your API should be running on `localhost:8080`

### Library

Alternatively import and use it directly

```go
package main

import (
	"net/http"

	"github.com/m3o/go-distributed"
)

func main() {
        http.HandleFunc("/upvotePost", VoteWrapper(true, false))
        http.HandleFunc("/downvotePost", VoteWrapper(false, false))
        http.HandleFunc("/upvoteComment", VoteWrapper(true, true))
        http.HandleFunc("/downvoteComment", VoteWrapper(false, true))
        http.HandleFunc("/posts", Posts)
        http.HandleFunc("/post", NewPost)
        http.HandleFunc("/comment", NewComment)
        http.HandleFunc("/comments", Comments)
        http.HandleFunc("/login", Login)
        http.HandleFunc("/readSession", ReadSession)
	http.ListenAndServe(":8080", nil)
}
```
