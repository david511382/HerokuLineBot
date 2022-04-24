#!/bin/sh

COMMAND=swag

if ! command -v $COMMAND &> /dev/null
then
    go install github.com/swaggo/swag/cmd/swag@v1.8.1
fi

swag init