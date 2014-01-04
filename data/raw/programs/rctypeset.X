#! /bin/sh
# rctypeset - typeset a recipe
#
# This shell script extracts a recipe from an alt.gourmand netnews article and 
# typesets it with troff. It first strips all communication headers
# (such as netnews or mail or notesfile headers) from the front and back
# of the recipe.
#
# Brian Reid, November 1985
# Copyright (C) 1986, USENET Community Trust
#
TROFF="DEFTROFF"
PATH=DEFPATH
export PATH
TMPFILE=/tmp/xart.$$
FRACTIONS=1
TMAC=/tmp/tmac.$$
METFLAG="-rMMETRIC"
VERBATIM=0

trap "rm -f $TMPFILE $TMAC; exit" 0 1 2 3 15
while [ $# -gt 0 ];
do
    case $1 in
	    -t) TFLAG="$1";;
	    -v) VERBATIM=1;;
	    -m) METFLAG="-rM1";;
	    -M) METFLAG="-rM0";;
	    -f) FRACTIONS=0;;
	    -*) echo $0: Unknown option $1 1>&2 ; exit;;
    esac
    shift
done

cat > $TMPFILE << 'extractor.awk'
BEGIN {Seen = 0}
/^\.RH MOD\.RECIPES/ {Seen = 1}
/^-- *$/ {Seen = 0}
/^\/\* End of text from/ {Seen = 0}
Seen == 1 {print $0}
extractor.awk

cat > $TMAC << 'tmac.recip'
TMAC.RECIP
tmac.recip
case $FRACTIONS in
  1) cat ;;
  *) sed -e 's/\([0-9a-zA-Z]\)\(\\([1-3]\)/\1 \2/g' \
	 -e 's|\\(12|1/2|' -e 's|\\(14|1/4|' -e 's|\\(34|3/4|' \
	 -e 's|``|"|'      -e "s|''|\"|" ;;
esac |
  case $VERBATIM in
      0) awk -f $TMPFILE;;
      1) awk '{print}';;
  esac |\
  $TROFF $TFLAG $METFLAG $TMAC -
