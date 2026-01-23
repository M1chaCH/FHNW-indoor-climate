#!/bin/zsh

SRC_DIR="../protobuf"
DST_DIR="./"

FILE_NAME=$1

protoc -I=$SRC_DIR --go_out=$DST_DIR "$SRC_DIR/$FILE_NAME"