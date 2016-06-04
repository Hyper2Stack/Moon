#!/bin/bash

set -e

CUR_DIR=`dirname $0`
TARGET_DIR=${CUR_DIR}/../target

[ -d ${TARGET_DIR} ] && rm -rf ${TARGET_DIR}/moon
mkdir -p ${TARGET_DIR}/moon

mkdir -p ${TARGET_DIR}/moon/usr/sbin
cp ${CUR_DIR}/../bin/moon ${TARGET_DIR}/moon/usr/sbin/moon
cp ${CUR_DIR}/../bin/moon-config ${TARGET_DIR}/moon/usr/sbin/moon-config

mkdir -p ${TARGET_DIR}/moon/etc/init.d
cp ${CUR_DIR}/moon ${TARGET_DIR}/moon/etc/init.d/moon

tar -zcf ${TARGET_DIR}/moon.tar.gz -C ${TARGET_DIR}  moon

rm -rf ${TARGET_DIR}/moon

echo "The moon package is locate:"
echo "`readlink -f ${TARGET_DIR}`/moon.tar.gz"
