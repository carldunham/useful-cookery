#! /bin/sh
# rcnew.t - make new cookbook pages with nroff or the typesetter
#   linked to rcnew.n
#
# Spencer Thomas, January 1986
# after rcbook.t by Brian Reid, November 1985
#
# Modifed by:  David Messer, Febuary 1987
#
TROFF="nroff -man"
FRACTIONS=0
RCINDEX="-w 56"
case $0 in
	*.t)	TROFF="DEFTROFF"
		FRACTIONS=1
		RCINDEX=-t
		;;
	esac
PATH=DEFPATH
export PATH
TMAC=TEMPDIR/tmac.$$
INDEX=TEMPDIR/Index.$$
ENDEX=TEMPDIR/Endex.$$
RECIPES=TEMPDIR/Recipes.$$
trap "rm -f TEMPDIR/*.$$; exit" 0 1 2 3 15
TFLAG=
FRACTIONS=1
METFLAG="-rMMETRIC"
KEEPDIR=DEFDIR
while [ $# -gt 0 ];
do
    case $1 in
	    -t) TFLAG="$1";;
	    -m) METFLAG="-rM1";;
	    -M) METFLAG="-rM0";;
	    -f) FRACTIONS=0;;
	    -*) echo $0: Unknown option $1 1>&2 ; exit;;
	    "") KEEPDIR=DEFDIR;;
	    *) KEEPDIR=$1;;
    esac
    shift
done
if [ ! -d $KEEPDIR ]; then
	echo rcnew: Your keep directory $KEEPDIR does not seem to exist. 1>&2
	exit 1
fi;
# if Index.keep exists, then a previous rcnew bombed.  Use the old one.
if [ ! -f Index.keep ]; then
	mv INDEX Index.keep
fi;
echo Typesetting cookbook update 1>&2
echo from files newer than \
    `ls -l Index.keep | awk '/./ {printf("%s %s %s", $5, $6, $7);}'` 1>&2
echo in $KEEPDIR 1>&2
rcindex $RCINDEX $KEEPDIR
cat > $TMAC << 'tmac.recip'
TMAC.RECIP
tmac.recip
set `date`
#  Thu Mar 27 22:56:49 PST 1986
YEAR="$6"
YR=`expr $YEAR - 1900`
TODAY="$3 $2 $YR"
sed -e "s/DATE/$TODAY/g" -e "s/YEAR/$YEAR/g" > $INDEX << 'Index'
.RH MOD.RECIPES-SOURCE INDEX Recipes "DATE" YEAR
.ll 6.5i
.nf
.br
.de XX
.tr ~ 
.ta 2.4iR 2.6i 5.3i
.ie !"\\$2""	\\$1~~~~~\\$2	\\$3 \\$4	\\$5
.el	\\$1	\\$3 \\$4	\\$5
..
.ps 8
.vs 10p
Index
cat > $ENDEX << 'Endex'
.vs
Endex
# This next kludge (using file $RECIPES) is because some versions of troff
# can take only a very small number of argument files.
cat `find $KEEPDIR -name '[a-z]*' -newer Index.keep -print | grep '^./[^/]*$' | sort` > $RECIPES
case $FRACTIONS in
  1) cat $INDEX INDEX $ENDEX $RECIPES;;
  *) cat $INDEX INDEX $ENDEX $RECIPES |\
	 sed -e 's/\([0-9a-zA-Z]\)\(\\([1-3]\)/\1 \2/g' \
	     -e 's|\\(12|1/2|' -e 's|\\(14|1/4|' -e 's|\\(34|3/4|' \
	     -e 's|``|"|'     -e "s|''|\"|" ;;
esac | \
$TROFF $METFLAG $TFLAG $TMAC - > Update
echo Update created 2>&1
rm Index.keep
