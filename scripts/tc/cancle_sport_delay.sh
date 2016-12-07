#! /usr/bin/env bash

# cancel packet delay

sudo tc qdisc del dev eth0 parent 1:2
sudo tc qdisc del dev eth0 root
