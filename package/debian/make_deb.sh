#!/bin/bash
cp dist/git-town-linux-amd64 package/debian/git-town
cd package/debian
sudo apt-get install build-essential debhelper devscripts
./debian_build
echo "success"
