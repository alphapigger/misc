#! /usr/bin/env bash

# emulate network latency in etcd cluster.
# set 30ms delay before sending packet that comes from source port 2380

sudo tc qdisc add dev eth0 root handle 1: htb default 1
sudo tc class add dev eth0 parent 1: classid 1:1 htb rate 1000kbps
sudo tc class add dev eth0 parent 1: classid 1:2 htb rate 1000kbps
sudo tc qdisc add dev eth0 parent 1:2 handle 10: netem delay 30ms
sudo tc filter add dev eth0 protocol ip parent 1: prio 1 u32 match ip sport 2380 0xffff flowid 1:2
