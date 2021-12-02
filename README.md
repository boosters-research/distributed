# Go Distributed [![godoc](https://godoc.org/github.com/m3o/go-distributed?status.svg)](https://godoc.org/github.com/m3o/go-distributed) [![Go Report Card](https://goreportcard.com/badge/github.com/m3o/go-distributed)](https://goreportcard.com/report/github.com/m3o/go-distributed)

How remote teams stay in sync

## Overview

Build a knowledge base, ask questions and stay in sync with your team all while doing it asynchronously.
Go Distributed is a building block for remote teams and communities. Write posts, leave comments and 
upvote the most relevant content.

## Demo

Find the demo running at [go-distributed.org](https://go-distributed.org)

## Usage

Distributed is built as a single Go program.

### API Key

Get an API key from [m3o.com](https://m3o.com/) and export as

```sh
export M3O_API_TOKEN=xxxxxx
```

### Server

Download and install

```sh
go get github.com/m3o/go-distributed/cmd/distributed
```

Run the API

```sh
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
        http.HandleFunc("/signup", Signup)
        http.HandleFunc("/readSession", ReadSession)
	http.ListenAndServe(":8080", nil)
}
```
