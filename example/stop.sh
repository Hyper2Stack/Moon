#!/bin/bash

ps -ef | grep mock-server | grep -v grep | awk '{print $2}' | xargs kill -9
sudo /usr/sbin/moon -s quit
