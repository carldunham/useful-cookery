.TH COOKBOOK 1 "12 Sep 87" "USENET Cookbook"
.SH NAME
cookbook: rckeep, rckeepnew, rcbook.t, rcbook.n, rctypeset, rcnroff \- Store and print the alt.gourmand USENET Cookbook
.SH SYNOPSIS
.na
.if t .ta 2.5i
.B rckeep < filename	
copies a recipe to your keep directory
.br
.B s\ |rckeep	
(same thing, from a news-reading program.)
.br
.B rckeepnew [ directory ]	
finds and keeps new recipes
.br
\fBrcbook.t [ \-t ] [ \-m ] [ \-M ] [ \-f ]	\fR
typesets a cookbook from kept recipes
.br
\fBrcbook.n [ \-m ] [ \-M ]	\fR
nroffs a cookbook from kept recipes
.br
\fBrctypeset [ \-m ] [ \-M ] [ \-f ]	\fR
typesets one recipe from std input
.br
\fBrcnroff [ \-m ] [ \-M ]	\fR
nroffs one recipe from std input
.br
\fBrcnew.t [ \-m ] [ \-M ] [ \-f ]	\fR
typesets new kept recipes
.br
\fBrcnew.n [ \-m ] [ \-M ]	\fR
nroffs new kept recipes
.ad
.\" This next line supposedly sets default tabs
.DT
.SH DESCRIPTION
The USENET Cookbook is an online database distributed in the 
.B alt.gourmand
newsgroup. This software makes it easy for you to clip recipes that interest
you, to store them in your own directory, and to make printed cookbooks from
those stored recipes.
.SH SAVING RECIPES ONE AT A TIME
.PP
Read the news with your favorite news reading program. If you don't have a
favorite, try 
.B rn.
When you see a recipe that interests you, type
.br
.B 	s |rckeep
.br
The ``s'' command means ``save''; the vertical bar means ``save into a
program''. 
.B rckeep
is the program that keeps recipes. It will save recipes into a directory
named Recipes, and it will create that directory for you if it does not
currently exist. If you would like recipes saved in some other place, you can
type
.br
.B 	s |rckeep PartyRecipes
.br
.B 	s |rckeep ~/PartyRecipes
.br
.B 	s |rckeep /usr/local/lib/Recipes
.br
and so forth.
.SH SAVING NEW RECIPES AUTOMATICALLY
If you don't want to take the time and trouble to run \fBrckeep\fR on every
recipe, you can run \fBrckeepnew\fR from time to time. When \fBrckeepnew\fR
is run, it looks in the netnews spooling directory to find the files that
hold the articles that hold the recipes, and automatically runs \fBrckeep\fR
on every file that contains a recipe and that has arrived since the last time
you ran \fBrckeepnew\fR. Like \fBrckeep\fR, \fBrckeepnew\fR can be given the
name of a directory to use as the keep directory:
.br
.B	rckeepnew
.br
or
.br
.B	rckeepnew /usr/local/lib/Recipes
.br
and so forth.
.SH HOW RECIPES ARE STORED
Every recipe is given a code word by the editor of alt.gourmand. That code
word is contained in the first few lines of the article body. The 
.B rckeep
and
.B rckeepnew
programs use that code word as a file name in your keep directory. For
example, if you are saving a recipe whose 1-line description is
.br
.B 	PHEASANT-PERRY \- Pheasant for Thanksgiving
.br
it will be saved in your keep directory under the filename ``pheasant-perry''.
These filenames are forced to lower case by
.B rckeep.
.SH MAKING A COOKBOOK
To typeset a cookbook, just type
.br
.B 	rcbook.t
.br
and all of the recipes in your keep directory will be typeset in alphabetical
order, with an index, introduction, and title page added. They will come out
either in English units (cups and teaspoons) or in metric units (grams and
deciliters) depending on how the alt.gourmand software has been installed at
your site. If you want to print the cookbook in cups and spoons regardless of
what the software installer did, then type
.br
.B	rcbook.t  \-M
.br
and if you want to print the cookbook in metric units regardless of what the
software installer did, then type
.br
.B	rcbook.t  \-m
.br
If you don't have a
typesetter or if you don't have the
.B troff(1)
program, then use ``\fBrcbook.n\fR'' instead; it uses 
.B nroff(1).
and sends the results of the formatting to standard output. It supports the
same \-m and \-M options that \fBrcbook.t\fR does.
.PP
If you have more than one keep directory, for specialty cookbooks of
different types, you can give rcbook a directory name as an argument, e.g.:
.br
.B	rcbook.t /usr/local/lib/Recipes
.PP
You can capture the output of \fBrcbook.t\fR in a file by using the ``\-t''
option:
.br
.B	rcbook.t \-t > FormattedFile
.SH PRINTING ONE RECIPE
The program \fBrctypeset\fR typesets one recipe, from its standard input.
It is able to discard news and mail headers from the beginning of the
recipe; you may use \fBrctypeset\fR on a netnews message in the form that it
is stored in a spooling directory, or you may use it directly from a
news-reading program. \fBrctypeset\fR supports the metric (\-m) and
non-metric (\-M) options explained above.
.PP
The program \fBrcnroff\fR performs the same function as \fBrctypeset\fR,
except that it formats to standard output instead of to the typesetter.
.SH PRINTING NEW RECIPES
It is a nuisance to keep track of which recipes you have printed and which
you have not. It is expensive and slow to print the entire cookbook every time
a few more recipes arrive. The \fBrcnew.t\fR program works very much like
\fBrcbook.t\fR, except that it typesets only those recipes that were put
into the keep directory after the last time you ran \fBrcnew.t\fR. It
produces a complete new index, though. It uses the creation time of the file
INDEX in your current directory to determine the last time you created an
index, which is presumably the last time you ran \fBrcbook.t\fR or
\fBrcnew.t\fR. \fBrcnew.t\fR supports the metric (\-m) and non-metric (\-M)
options explained above.
.PP
The program \fRrcnew.n\fR does the same thing as \fRrcnew.t\fR, except that
it formats to standard output instead of to the typesetter.
.SH THE INDEX
The cookbook index is a ``permuted index''. This means that each recipe is
indexed by every major word in its title. Some people seem to find this
format jarring, and more than one person has filed a bug report upon seeing
the index, certain that the output that it produces cannot possibly be
correct.
.PP
Since your notebook of recipes will grow each week, the pages cannot have
numbers. Instead, the pages are identified by the code name for each recipe,
with the expectation that you will keep your notebook in alphabetical order
by this name.
.SH SPECIAL CHARACTERS
Troff has a number of special characters, such as such as
typeset fractions (``\(12'' for 1/2), which are not available on some brands
of typesetters and laser printers. If your typesetter handles those
characters incorrectly, then you should add the ``\-f'' option to rcbook.t,
rctypeset, and rcnew.t. The ``\-f'' option causes the recipe software to
translate those special characters into standard ASCII characters before
printing.
.SH COPYRIGHT
The entire USENET cookbook is copyrighted by the USENET Community Trust. The
purpose of this copyright is to prevent commercial use of the material. The
USENET Community Trust permits any noncommercial use of the contents of the
recipe database, and forbids any commercial use.
.SH FILES
.ta 3i
NEWSDIR/alt/gourmand	USENET recipe directory
.br
DEFDIR	Your keep directory (saved recipes)
.br
DEFDIR/.keepnew	Record of when rckeepnew last run
.br
INDEX	Record of when rcnew last run
.br
OBJDIR/rc*	The recipe software
.SH SEE ALSO
recipes(5), rn(1)
.SH AUTHOR
Brian K. Reid, DEC Western Research Laboratory
