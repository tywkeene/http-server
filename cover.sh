#!/bin/bash
#Somewhat copied from https://gist.github.com/hailiang/0f22736320abe6be71ce

echo "mode: count" > profile.out
for dir in $(find . -type d); do
    if [ -e $dir/*_test.go ]; then
        go test -v -covermode=count -coverprofile=$dir.tmp $dir || exit -1
        cat $dir.tmp | tail -n +2 >> profile.out || exit -1
        rm $dir.tmp
    fi
done

go tool cover -html=profile.out
rm profile.out
