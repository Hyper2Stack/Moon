#!/bin/bash

sudo /usr/sbin/moon-config -key "abc"
sudo service moon start

rm -rf /tmp/mock-server.log
`dirname $0`/../bin/mock-server > /tmp/mock-server.log 2>&1 &
