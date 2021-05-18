#!/bin/bash

# go build

if [ $1 == 'client' ]
then
	./tuntap client 4 &
	ifconfig wg2 172.16.0.2/30
else
	./tuntap server 4 &
	ifconfig wg2 172.16.0.1/30
fi
ifconfig wg2 mtu 1470
ifconfig wg2 txqueuelen 2000

