#!/bin/bash
cd package/debian
sudo apt-get install build-essential debhelper devscripts
./debian_build
echo "success"
