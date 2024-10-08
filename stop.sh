#!/bin/bash
kill -9 `pgrep -f 'go run'`
kill -9 `pgrep -f 'registrar'`
kill -9 `pgrep -f 'resolver'`
kill -9 `pgrep -f 'exe/main'`
