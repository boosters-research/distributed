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
