#!/bin/sh
#
# Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
# Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
# Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
# Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
# Vestibulum commodo. Ut rhoncus gravida arcu.
#

rm -rf profile-manager
go build .
BUILD_VERSION=latest
docker build --no-cache -t edgegallery/profile-manager:${BUILD_VERSION} .