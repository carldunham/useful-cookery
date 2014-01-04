#! /bin/sh 
# Script to configure and install the other shell scripts for alt.gourmand
#	Brian Reid, December 1985
#	v2	    April 1986
#	v3	    June  1986
#	v3.1	    Oct   1987
#
SYSTYPE="$1"
OBJDIR="$2"
NEWSDIR="$3"
DEFDIR="$4"
DEFPATH="$5"
DEFTROFF="$6"
DEFPAGER="$7"
METRIC="$8"
TEMPDIR="$9"

case $SYSTYPE in
    XXX) if [ -d /usr/ucb ]
	 then
	     SYSTYPE="BSD"
	 else
	     SYSTYPE="USG"
	 fi;;
esac;


# make rcnroff.X from rctypeset.X
sed -e "s+DEFTROFF+nroff -man+" \
    -e "s+typesets it with troff+formats it with nroff+" \
    -e 's+^FRACTIONS=1$+FRACTIONS=0+' \
    -e "/-t)/d" \
    -e '/^\$TROFF/s+$+ > Update.nr+' \
    -e 's/$TFLAG//g' \
    < rctypeset.X > rcnroff.X

# Make a tmac file tuned to rthe type of system we are on, and strip comments.
sed -e "/^@@@$SYSTYPE/s/^@@@$SYSTYPE//" \
    -e '/^@@@/d' \
    < tmac.recip | grep -v '^\.\\"' > tmac.tmp

# now configure the executable programs
for i in *.X
do
    j=`basename $i .X`
    echo Processing $i into $j
    sed -e "s+OBJDIR+$OBJDIR+g" \
        -e "s+NEWSDIR+$NEWSDIR+g" \
        -e "s+TEMPDIR+$TEMPDIR+g" \
	-e "s+DEFDIR+$DEFDIR+g" \
	-e "s+DEFPATH+$DEFPATH+g" \
	-e "s+DEFPAGER+$DEFPAGER+g" \
	-e "s+METRIC+$METRIC+g" \
	-e '/^TMAC.RECIP$/r tmac.tmp' \
	-e '/^TMAC.RECIP$/d' \
	-e "s+DEFTROFF+$DEFTROFF+g" < $i > $j
    chmod 755 $j
done
rm -f rcbook.n; ln rcbook.t rcbook.n
rm -f rcnew.n; ln rcnew.t rcnew.n
