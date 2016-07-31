#!/bin/bash

set -e

CUR_DIR=`dirname $0`
TARGET_DIR=${CUR_DIR}/../target

[ -d ${TARGET_DIR} ] && rm -rf ${TARGET_DIR}
mkdir -p ${TARGET_DIR}

mkdir -p ${TARGET_DIR}/usr/sbin
cp ${CUR_DIR}/../bin/moon ${TARGET_DIR}/usr/sbin/moon
cp ${CUR_DIR}/../bin/moon-config ${TARGET_DIR}/usr/sbin/moon-config

mkdir -p ${TARGET_DIR}/etc/init.d
cp ${CUR_DIR}/moon ${TARGET_DIR}/etc/init.d/moon

tar -zcf ${TARGET_DIR}/moon.tar.gz -C ${TARGET_DIR} usr etc

rm -rf ${TARGET_DIR}/usr ${TARGET_DIR}/etc

echo "The moon package is locate:"
echo "`cd ${TARGET_DIR};pwd`/moon.tar.gz"
