#!/bin/bash

# start fledge
/usr/local/fledge/bin/fledge start
/usr/local/fledge/bin/fledge status

while :
do
	sleep 300
done