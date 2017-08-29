#!/bin/bash
cp dist/git-town-linux-amd64 package/debian/git-town
cd package/debian
sudo apt-get -qq update
sudo apt-get install -y build-essential debhelper devscripts
./debian_build.sh
cd ../..
cp package/gittown_4.2.1_amd64.deb dist
