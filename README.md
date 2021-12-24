# Go Distributed [![godoc](https://godoc.org/github.com/asim/go-distributed?status.svg)](https://godoc.org/github.com/asim/go-distributed) [![Go Report Card](https://goreportcard.com/badge/github.com/asim/go-distributed)](https://goreportcard.com/report/github.com/asim/go-distributed)

How remote teams stay in sync

## Overview

Build a knowledge base, ask questions and stay in sync with your team all while doing it asynchronously.
Go Distributed is a building block for remote teams and communities. Write posts, leave comments and 
upvote the most relevant content.

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
go get github.com/asim/go-distributed/cmd/distributed
```

Run the API

```sh
distributed
```

Your API should be running on `localhost:8080`
