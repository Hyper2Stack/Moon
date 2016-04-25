#!/bin/bash

echo "--- get key ---"
curl -i http://localhost:8080/key
echo

echo "--- get agent info ---"
curl -i "http://localhost:8080/test/agent-info"
echo

echo "--- get node info ---"
curl -i "http://localhost:8080/test/node-info"
echo

echo "--- success to exec shell ---"
curl -i -X POST -d '{"commands":[{"command":"ls", "args":["/"], "restrict":true}]}' "http://localhost:8080/test/shell"
echo

echo "--- fail to exec shell ---"
curl -i -X POST -d '{"commands":[{"command":"ls", "args":["/not-exist"], "restrict":true}]}' "http://localhost:8080/test/shell"
echo
