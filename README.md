# Go Rate Limiter

This is my own implementation of an API rate-limiter middleware in Go using different algorithms. Coding Challenges series by John Crickett https://codingchallenges.fyi/challenges/challenge-rate-limiter/

- Token Bucket
- Fixed Window
- Sliding Window Log
- Sliding Window Counter
- Sliding Window Counter across multiple servers using Redis

## Requirements

- Install [GO v1.21](https://go.dev/dl/)
- Install Make [Windows](https://gnuwin32.sourceforge.net/packages/make.htm) | [Mac] brew install make

## Setup

```
Make build
```

```
Make install
```

## Run

```
go-rate-limiter <algorithm> --flag1 --flag2
```

## Tests

```
go test -v ./... --cover
```

## Stress Tests

I am using the `loadtest` package to run stress tests against the API server directly from the CLI. This can also be achieved by installing `Postman` or .

1. Install loadtest package using NPM

```
npm install -g loadtest
```

2. Run the server

```
go-rate-limiter tokenBucket --capacity 20 --refillRate 1
```

3. Run the tests

```
loadtest -c <number of concurrent users> --cores <number of cpu cores> --rps <target request per second by client> -k -n <number of requests to perform> http://localhost:8080/limited
```

## Packages

1. Gin: Fast Web/API web server

2. Cobra: CLI: Used to create CLI commands

3. Go-redis/v9: Read/write data to Redis

4. Testify: Testing utilities, easy assertions, mocking, etc
