#!/bin/sh
go build -o player main.go
for i in `seq 1 11`
do
  ./player -team=away -number=$i&
  ./player -team=home -number=$i&
done



