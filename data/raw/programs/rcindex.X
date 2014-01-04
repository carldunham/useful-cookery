#! /bin/sh
# rcindex - produce index of recipes
#
# usage:
#   rcindex  [ -options ] [ keepdir ]
#
# This shell script produces an index of all recipes in a keep directory.
# That index is an nroff/troff document that will be used as a part
# of the book produced from that directory.
#
# If you don't specify a keep directory, it will use DEFDIR, which
# is as good a place as any.
#
# Brian Reid, November 1985
#	v2, June 1986
# Copyright (C) 1986, USENET Community Trust
#
PATH=DEFPATH
export PATH
BADWDS=TEMPDIR/rcindex.$$
PUNCHR=TEMPDIR/rcipunc.$$
trap "rm -f TEMPDIR/*.$$; exit" 0 1 2 3 15
while [ $# -gt 1 ];
do
    case $1 in
	    -*) OPTIONS="$OPTIONS $1";;
	    *) OPTIONS="$OPTIONS $1";;
    esac
    shift
done
case $1 in
	"") KEEPDIR=DEFDIR;;
	*) KEEPDIR=$1;;
esac
if [ ! -d $KEEPDIR ]; then
	echo rcindex: Your keep directory $KEEPDIR does not seem to exist.
	exit 1
fi;
cat > $PUNCHR << 'endflag'
()/.,
endflag
tr ' ' '\12' << 'endflag' | sort -u > $BADWDS
& - a about addictive after against all also amazingly an and
another any are as at back basic be because been before being best-ever
between both but by came can center classic cloying come con cook
cooked cooking could countries day decadent delicate delicious did
dinner dish dishes do down each earthy easy-to-make even extremely fattening
filling first flavor flavorful food for from get go good great had has have
he hearty her here him his hot how i if in incredible into intoxicating is
it its just know large last life light like little long luscious made
make makes man many may me men might more most mr much must my never
new no not now of off old on one only or other our out over own people
picture-perfect prizewinning recipe recipes refreshing rich right said
same see she should since smothered so some state still succulent such
take than that the their them then there these they thick this those
three through time to too two ultimate under universe up us used very
was way we well were what when where which while who will with
wonderful work world would year years you your
endflag
DIR=`pwd`
cd $KEEPDIR
grep \^.RZ [a-z]* |\
  sed -e 's/\: *.RZ //' |\
  awk -F\" '{print $1,$4}' |\
  ptx -f -r -i $BADWDS -b $PUNCHR $OPTIONS |\
  sed -e 's/^\.xx/.XX/' > $DIR/INDEX
