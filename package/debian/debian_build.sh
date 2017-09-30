#!/bin/sh
# makes a fakeroot structure and copies files to locations
mkdir -pv usr/bin &&
cp git-town usr/bin &&
chmod +x usr/bin/git-town &&
debuild -us -uc
