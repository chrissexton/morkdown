#!/bin/bash
set -e

go build
set +e
pkill morkdown
set -e
./morkdown &
