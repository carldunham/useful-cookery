#! /bin/sh
# rckeep - keep a recipe by adding it to a keep directory
#
# usage:
#   rckeep [ keepdir ]  [ filename ]
#
# This program lets you snatch recipes out of newsgroups and store them in
# a "keep" directory. They are stored in a file name that matches the 
# recipe id, which is found in the .RH command in the recipe.
#
# If you don't specify a keep directory, it will use DEFDIR which
# is as good a place as any. You can use ~/ to stand for $HOME/ if you want.
#
# Brian Reid, November 1985
# Copyright (C) 1986, USENET Community Trust
#
PATH=DEFPATH
export PATH
TEMPFILE=/tmp/rckeep.$$
case $1 in
	"") KEEPDIR=DEFDIR;;
	~/*) KEEPDIR=$HOME/`expr $1 : '~/\(.*\)'`;;
	~*) echo Sorry, I cannot process ~ notation, except for ~/...;
	       exit 1;;
	*) KEEPDIR=$1;;
esac
if [ ! -d $KEEPDIR ]; then
	echo Creating new keep directory $KEEPDIR;
	mkdir $KEEPDIR;
fi;
case $2 in
	"") rcextract > $TEMPFILE;;
	*) rcextract < $2 > $TEMPFILE;;
esac
# the [A-Z] in tr arg below is for compatibility with SysV
IDLINE=`awk '{print $2,$3; exit;}' < $TEMPFILE | tr "[A-Z]" "[a-z]" `
case $IDLINE in
    "") echo This is not an alt.gourmand datafile. I do not know what to do with it.
	exit 1;;
esac
set $IDLINE
KEEPNAME=$2
FILETYPE=$1
case $FILETYPE in
    mod.recipes-source) if [ -f $KEEPDIR/$KEEPNAME ]; then
	echo Overwriting previous version of `basename $KEEPNAME`:
	ls -l $KEEPDIR/$KEEPNAME
	fi
	mv $TEMPFILE $KEEPDIR/$KEEPNAME
	echo Saving $KEEPDIR/$KEEPNAME;;
    *) echo This is not an alt.gourmand datafile. I do not know what to do with it.;
       exit 1;;
esac
