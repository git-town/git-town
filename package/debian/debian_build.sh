#!/bin/sh
# makes a fakeroot structure and copies files to locations
mkdir -pv usr/lib/git-core &&
mkdir -pv usr/share/man/man1 &&
cp -r ../../src/* usr/lib/git-core &&
cp -r ../../man/man1 usr/share/man/man1 &&
chmod +x usr/lib/git-core/drivers/code_hosting/*.sh &&
chmod +x usr/lib/git-core/helpers/*.sh &&
chmod +x usr/lib/git-core/helpers/git_helpers/*.sh &&
debuild -us -uc &&
rm -rf usr
