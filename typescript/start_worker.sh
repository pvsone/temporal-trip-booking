#!/bin/bash
# TODO: remove this when the typescript sdk implements the default path on macOS properly
export TEMPORAL_CONFIG_FILE=~/Library/Application\ Support/temporalio/temporal.toml
npm install
npm run start.watch
