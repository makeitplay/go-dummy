#!/bin/sh

if [ -z "$1" ]
  then
    echo "Please, pass the team's Docker image name as first argument"
    exit 1
fi

if [ -z "$2" ]
  then
    echo "Please, pass the the team side (home or away) as second argument"
    exit 1
fi
HOST_IP=`hostname -I | awk '{print $1}'`
for i in `seq 1 11`
do
    docker run $1 -team=$2 -number=$i -wshost=$HOST_IP &
done
