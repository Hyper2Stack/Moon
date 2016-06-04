#!/bin/bash

sudo /usr/sbin/moon-config -key "abc"
sudo service moon start
`dirname $0`/../bin/mock-server &
