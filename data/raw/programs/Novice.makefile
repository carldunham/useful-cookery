### This file is used as part of the Novice installation procedure. Please
### read file "Novice".
##=
# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
# (1) Where is the news spooling directory? If you don't have any idea, then
# leave this line alone--this is the best guess. If you get this wrong,
# all of the programs will still work, though they may give strange error
# messages at times.

NEWSDIR=/usr/spool/news

# (2) What version of troff works best at your site? Change the word "troff" 
# to the name of your version, if it is different.
#
# If you don't know, then try one of these guesses:
# If you have an Apple LaserWriter, your best guess is "ptroff"
# If you have an Imagen printer, your best guess is "itroff"
# If you have a QMS or Talaris printer, your best guess is "qtroff"
# Otherwise try just "troff", or ask for help.
# If you get this wrong, you will be able to save recipes but you won't be
# able to typeset them or print them out.

MYTROFF=troff

# (3) What is your favorite pager and what are its options? If you don't have
# any idea, then leave this line alone. If you have a "pg" program, then edit
# the word "cat" into "pg" below. If you have a "more" program, then edit the
# word "cat" into "more". Even if you get this completely wrong, almost
# everything will still work.

DEFPAGER=cat

# (4) Do you want the recipes to come out in metric units or English units?
# If you think that recipes should use teaspoons and cups, then leave the
# METRIC=0 line alone. If you think that recipes should use grams and
# milliliters, then edit this line to say METRIC=1

METRIC=0

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # 
# You should not have to change anything below this line, but you are
# welcome to make changes if you think you know what you are doing.

TEMPDIR=/tmp
DEFTROFF=$(MYTROFF) -man
DEFDIR=$$HOME/Recipes
OBJDIR=$$HOME
DEFPATH=:$(OBJDIR):XXDEFPATH
OBJECTS=rcbook.t rcbook.n rcextract rcindex rckeep rcnroff rcshow rctypeset \
	rcintro

configure: Makefile 
	@sh configure.sh XXX '$(OBJDIR)' '$(NEWSDIR)' '$(DEFDIR)' \
	 	'$(DEFPATH)' '$(DEFTROFF)' '$(DEFPAGER)' '$(METRIC)' \
		'$(TEMPDIR)'
	@cat Novice.step2

install:
	@Novice.install $(OBJECTS) $(OBJDIR)
