#!/bin/bash
set -eo pipefail

if grep --color=always "if l.buf.peek() ==" ./**/*.go; then
    echo "peek() is not free. these peek token can maybe be reused"
fi

if grep --color=always "for peek :=" ./**/*.go; then
    echo "peek() is not free. these peek token can maybe be reused"
fi

if grep --color=always "peek.Kind != token.KindEOF &&" ./**/*.go; then
    echo "EOF is maybe less probable than the following condition(s):"
fi