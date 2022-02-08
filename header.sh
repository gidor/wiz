#!/bin/bash

SOURCE_DIR=`dirname $0`
shopt -s globstar
for i in **/*.go # or whatever other pattern...
do
    echo $i
    cat $i | awk -f $SOURCE_DIR/rmheader.awk > $i.new && cat $SOURCE_DIR/copyright.txt $i.new > $i 
    rm $i.new
done
