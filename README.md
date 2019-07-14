# Rakuten take home assignment - Implement a Chat Server

Rakuten chat server written in Go and Node.js.

## This project contains three components:
* Chat server (written in Go)
* Migration scripts
* Test scripts (written in Node.js with mocha framework)

### Build chat server and test scripts docker images

1. cd to project's root directory

2. Build docker images:

`./docker/build.sh`

3. Start chat server + run migrations

`./docker/run_chat_server.sh`

4. Run test scripts

`./docker/run_test_cases.sh`
