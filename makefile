# Go parameters
GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOINSTALL=$(GOCMD) install
GOGET=$(GOCMD) get
DOCKERCMD=docker-compose
BINARY_NAME=bin/goratelimiter

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

install:
	$(GOINSTALL)

test:
	$(GOTEST) -v ./... --cover

run-docker:
	${DOCKERCMD} up --build

all: test run