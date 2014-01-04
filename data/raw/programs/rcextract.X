#! /bin/sh
# rcextract - extract article from alt.gourmand
#
# This shell script extracts the portion of an alt.gourmand
# article that is supposed to be fed into xroff. It strips off the
# article header and any signature, and sends the resulting text to stdout.
#
# Brian Reid, November 1985
# Paul Pomes, March 1986: updated to work with notesfiles
#
# Copyright (C) 1986, USENET Community Trust

PATH=DEFPATH
export PATH
TMPFILE=/tmp/xart.$$
trap "rm -f $TMPFILE; exit" 0 1 2 3 15
cat > $TMPFILE << 'extractor.awk'
BEGIN {Seen = 0}
/^\.RH MOD\.RECIPES/ {Seen = 1}
/^-- *$/ {Seen = 0}
/^\/\* End of text from/ {Seen = 0}
Seen == 1 {print $0}
extractor.awk
awk -f $TMPFILE
