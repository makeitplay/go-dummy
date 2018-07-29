#!/bin/sh
go build -o player main.go
for i in `seq 1 11`
do
    docker run the-dummies -team=home -number=$i -wshost=$1 &
    docker run the-dummies -team=away -number=$i -wshost=$1 &
done



