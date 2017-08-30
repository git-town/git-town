#!/bin/bash

VER_NUM=$(echo $TRAVIS_TAG | cut -d 'v' -f 2) # extract version number
TIMESTAMP=$(date "+%a, %d %b %Y %H:%M:%S %z")

cat << EOF > package/debian/debian/changelog
gittown ($VER_NUM) RELEASED; urgency=low

  * See https://github.com/Originate/git-town/releases for a complete changelog

 -- allonsy <alec.snyder@originate.com> $TIMESTAMP
EOF

cp dist/git-town-linux-amd64 package/debian/git-town
cd package/debian
sudo apt-get -qq update
sudo apt-get install -y build-essential debhelper devscripts
./debian_build.sh
cd ../..
cp package/gittown_"$VER_NUM"_amd64.deb dist/git-town-amd64.deb
