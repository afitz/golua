#!/bin/sh

LUAVER=$(pkg-config --print-provides lua | cut -d ' ' -f1)

if test "$LUAVER" = ""; then
    LUAVER=lua5.1
fi

echo $LUAVER
