#!/bin/bash

echo "--- get client info ---"
curl -i http://localhost:8080/client
echo

echo "--- test ---"
curl -i "http://localhost:8080/test"
echo
