#!/bin/bash

user="$(id -un 2>/dev/null || true)"
if [ "${user}" != "root" ]; then
    echo "Need root permission to uninstall package"
    exit 1
fi

/usr/sbin/moon -s quit > /dev/null 2>&1

rm -rf /etc/moon
rm -rf /var/run/moon.pid
rm -rf /var/log/moon
rm -rf /var/run/moon

rm -f /usr/sbin/moon
rm -f /usr/sbin/moon-config
