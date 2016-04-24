#!/bin/bash

set -e

user="$(id -un 2>/dev/null || true)"
if [ "${user}" != "root" ]; then
    echo "Need root permission to install package"
    exit 1
fi

ROOT_DIR=`dirname $0`

mkdir -p /etc/moon
mkdir -p /var/run/moon
mkdir -p /var/log/moon

install ${ROOT_DIR}/moon /usr/sbin/moon
install ${ROOT_DIR}/moon-config /usr/sbin/moon-config
