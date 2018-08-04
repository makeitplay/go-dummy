#!/bin/sh

if [ -z "$1" ]
  then
    echo "Please, pass the first argument (home or away) to set the team side"
    exit 1
fi

go build -o the-dummies main.go || { echo "building has failed"; exit 1; }
for i in `seq 1 11`
do
  ./the-dummies -team=$1 -number=$i -wshost=$2 &
done
