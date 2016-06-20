#!/bin/bash

set -e

user="$(id -un 2>/dev/null || true)"
if [ "${user}" != "root" ]; then
    echo "Should have root permission, good luck!"
    exit 1
fi

ROOT_DIR=`dirname $0`/..

mkdir -p /etc/moon
mkdir -p /var/run/moon
mkdir -p /var/log/moon

install ${ROOT_DIR}/config/moon.cfg /etc/moon/moon.cfg
install ${ROOT_DIR}/bin/moon /usr/sbin/moon
install ${ROOT_DIR}/bin/moon-config /usr/sbin/moon-config

install ${ROOT_DIR}/package/moon /etc/init.d/moon

# TODO: for OSX??
if [ -f /etc/debian_version ]; then
    update-rc.d moon defaults > /dev/null
elif [ -f /etc/redhat-release ]; then
    chkconfig --add moon
    chkconfig moon on
else
    echo "Unknown distribution, service auto startup not enabled!"
fi
