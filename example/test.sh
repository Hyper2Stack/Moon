#!/bin/bash

echo "--- get key ---"
curl -i http://localhost:8080/key
echo

echo "--- test ---"
curl -i "http://localhost:8080/test"
echo
