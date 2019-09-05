#!/usr/bin/env bash

# -I 指定import路径，可以指定多个-I参数，编译时按顺序查找，不指定时默认查找当前目录
# --go_out ：golang编译支持，支持以下参数
# plugins=plugin1+plugin2 - 指定插件，目前只支持grpc，即：plugins=grpc
protoc -I . --go_out=plugins=grpc:. *.proto
