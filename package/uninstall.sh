#!/bin/bash

user="$(id -un 2>/dev/null || true)"
if [ "${user}" != "root" ]; then
    echo "Should have root permission, good luck!"
    exit 1
fi

service moon stop > /dev/null 2>&1
update-rc.d -f moon remove > /dev/null 2>&1

rm -rf /etc/init.d/moon

rm -rf /etc/moon
rm -rf /var/run/moon.pid
rm -rf /var/log/moon
rm -rf /var/run/moon

rm -f /usr/sbin/moon
rm -f /usr/sbin/moon-config
