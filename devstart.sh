#!/usr/bin/env bash

nodemon --exec go run tikiserver.go --signal SIGTERM
