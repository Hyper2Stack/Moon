#!/bin/bash

export GOPATH=$(readlink -f $(dirname $0))
go install moon
go install moon-config
go install mock-server
