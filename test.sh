#!/bin/sh
set -e

echo append
go test -- features/append
echo branch
go test -- features/branch
echo completions
go test -- features/completions
echo compress
go test -- features/compress
echo config
go test -- features/config
echo contribute
go test -- features/contribute
echo delete
go test -- features/delete
echo detach
go test -- features/detach
echo diff_parent
go test -- features/diff_parent
echo down
go test -- features/down
echo end_to_end
go test -- features/end_to_end
echo feature
go test -- features/feature
echo hack
go test -- features/hack
echo help
go test -- features/help
echo init
go test -- features/init
echo merge
go test -- features/merge
echo observe
go test -- features/observe
echo offline
go test -- features/offline
echo park
go test -- features/park
echo perennial
go test -- features/perennial
echo prepend
go test -- features/prepend
echo propose
go test -- features/propose
echo prototype
go test -- features/prototype
echo rename
go test -- features/rename
echo repo
go test -- features/repo
echo set_parent
go test -- features/set_parent
echo shared
go test -- features/shared
echo ship
go test -- features/ship
echo status
go test -- features/status
echo swap
go test -- features/swap
echo switch
go test -- features/switch
echo sync
go test -- features/sync
echo up
go test -- features/up
echo version
go test -- features/version
echo walk
go test -- features/walk
echo
go test -- features/
