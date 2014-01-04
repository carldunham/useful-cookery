#! /bin/sh
# rckeepnew - keep all recipes that arrived since the last time this
#	      program was run.
#
# usage:
#   rckeepnew [ keepdir ]
#
# This program looks in the USENET spooling directory and finds all files
# that contain recipes and that were created since the last time it was run.
# For each such file, it runs the program "rckeep". The record of the last
# time the program was run is kept as the creation date of the file
# ".keepnew" in your keep directory.
#
# Brian Reid, March 1986
# Copyright (C) 1986, USENET Community Trust
#
PATH=DEFPATH
export PATH

# if you are not running B news, you will need to edit the following line to
# tell it the name of the directory in which alt.gourmand articles are stored.
RECIPEDIR=NEWSDIR/alt/gourmand
export RECIPEDIR

KEEPDIR=${1-DEFDIR}
if [ ! -d $RECIPEDIR ]
then
    echo rckeepnew: Cannot find ${RECIPEDIR}. 
    exit 1
fi
cd $RECIPEDIR
if [ ! -d $KEEPDIR ]
then
    echo rckeepnew: Directory $KEEPDIR does not exist. I will try to create it.
    echo mkdir $KEEPDIR
    if mkdir $KEEPDIR
    then
	echo rckeepnew: mkdir succeeded. Continuing.
    else
	echo rckeepnew: mkdir failed. Cannot continue.
	exit 1
    fi
fi

# Find the files that have changed 
for j in `if test -f $KEEPDIR/.keepnew
then
    find . -newer $KEEPDIR/.keepnew -exec grep -l "^Subject: RECIPE:" {} \;
else
    grep -l "^Subject: RECIPE:" *
fi`
# run rckeep on each of the changed files
do
    rckeep $KEEPDIR < $j
done
# update the date-time flag
touch $KEEPDIR/.keepnew
echo Date-and-time flag in $KEEPDIR updated.
