#!/bin/bash

sudo /usr/sbin/moon-config -key "abc"
sudo /usr/sbin/moon
`dirname $0`/../bin/mock-server &
