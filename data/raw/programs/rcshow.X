#! /bin/sh
# rcshow - show a recipe
#
# This shell script extracts a recipe from an alt.gourmand netnews article and 
# nroffs it for display on the screen. This differs slightly from "rcnroff"
# in that it is not page-oriented. Actually, this program is almost too slow
# to use, because it has to run nroff every time.
#
# Brian Reid, November 1985
# Copyright (C) 1986, USENET Community Trust
#
# If you don't have "less", you will want to 
# add "ul" to the pipeline to get the underlines right.

TROFF="nroff -man"
PAGER="DEFPAGER"

PATH=DEFPATH
export PATH
TMPFILE=/tmp/xart.$$
TMAC=/tmp/tmac.$$
trap "rm -f $TMPFILE $TMAC; exit" 0 1 2 3 15
METFLAG="-rMMETRIC"

cat > $TMPFILE << 'extractor.awk'
BEGIN {Seen = 0}
/^\.RH MOD\.RECIPES/ {Seen = 1}
/^-- *$/ {Seen = 0}
Seen == 1 {print $0}
extractor.awk
cat > $TMAC << 'tmac.recip'
TMAC.RECIP
tmac.recip
awk -f $TMPFILE |\
  sed -e 's/\([0-9a-zA-Z]\)\(\\([1-3]\)/\1 \2/g' \
	 -e 's|\\(12|1/2|' -e 's|\\(14|1/4|' -e 's|\\(34|3/4|' \
	 -e 's|``|"|'      -e "s|''|\"|" |\
  $TROFF $METFLAG  $TMAC - | $PAGER
# if you don't have "less", you will probably want to change the above line
# into $TROFF $TMAC - | ul | $PAGER
