#!/bin/bash

export GOPATH=$(cd `dirname $0`; pwd)
go install moon
go install moon-config
go install mock-server
