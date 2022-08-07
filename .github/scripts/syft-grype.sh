#!/bin/sh

syft -q -o json "$@" | grype -q
