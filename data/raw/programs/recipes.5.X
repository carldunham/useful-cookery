.TH RECIPES 5 "28 Jun 87" "USENET Cookbook"
.SH NAME
recipes \- USENET Cookbook (alt.gourmand) format and submitting procedure.
.SH SYNOPSIS
Mail your recipe to recipes@decwrl.dec.com. Try not to plagiarize.
If you want to put in your own troff commands, use the Unix manual macros.
.SH DESCRIPTION
Alt.gourmand is a ``moderated newsgroup''. This means that you mail your
submissions to the moderator, and he distributes them.
The moderator's mailbox is
.br
	decwrl!recipes,
.br
if you believe in that kind of address, or
.br
	{decvax,ucbvax,sun,pyramid,hplabs,ames}!decwrl!recipes
.br
if you believe in the other kind of address. The addresses
.br
	recipes@decwrl.DEC.COM
.br
	decwrl.dec.com!recipes
.br
also work, if your computer knows how to mail to them (many don't).
.PP
If your netnews administrator has set things up properly and if you are
running the right news software, you can just post to alt.gourmand as you
would any other group, but your posting will be mailed to the moderator
insted of being placed directly in the newsgroup.
.SH HOW TO FORMAT A RECIPE
Every recipe that goes out will be formatted with a small set of troff
commands. If you don't know any troff, then just send the text of your
recipe. If you can do the formatting yourself, then please do! That will
save the editor a lot of time, and your recipe will go out sooner.  Some
hints for how to do it are in a later section of this document.  If you
don't know about xroff/troff/nroff but you do want to be helpful, then you
can help a lot by doing these things:
.PP
.IP "(1)" 5
Put your recipe in the standard sequence.
.IP "(2)" 5
Don't put any tab characters (^I) in the file if you can possibly avoid it.
They do strange things on typesetters.
.IP "(3)" 5
If you are a troff wizard, please don't use any of that wizardry in these
recipes. If you stick to the set of commands used in the Unix manual macros
(see \fIman(7)\fR), and the set of commands that are part of the USENET
Cookbook package (see following section), then things should work OK. 
.br
There are plenty of places in troff where you can get away with not
using quotes around macro arguments. Please use quotes, even when they are
redundant, such as in ``.IG "1" "onion"''  This is because the indexing and
cross-referencing programs expect to find the quotes, even though troff can
work without them.
.RE
.SH HOW TO WRITE A RECIPE
Please try to put your recipe in the standard sequence. Different cookbooks
use different standards. This is the sequence that the \fIUSENET Cookbook\fR
uses.
.RS
.IP (a) 5
Title and 1-line description. The 1-line description will be used to index
the recipe, so make it as descriptive as possible. Avoid words like
``delicious'' or ``yummy''. We expect all of these recipes to be delicious.
.IP (b) 5
Introductory commentary, explaining (if possible) where you got
the recipe from and what you like about it. If you got the recipe from a
cookbook, give the title and author of that cookbook.
.IP (c) 5
List of ingredients, using ``Tbsp'', ``tsp'', ``cup'', ``oz'', and ``lb''
for English units, or ``ml'', ``cl'', ``dl'', ``l'', ``g'', and ``kg'' for metric
units. Stay away from ``pint'', ``quart'', and ``gallon'', because they have
different meanings in different countries. Do not use 1-letter abbreviations
for US measures: don't use "c" for "cup" or "T" for "Tbsp". An ``oz'' is a
fluid ounce or an avoirdupois ounce depending on context.
The ingredients should be listed in the order they will
be used. Don't capitalize ingredient names unless they are proper nouns. 
Avoid terms like ``1 box'' or ``1 can'' or ``1 package'', because packaging
conventions vary widely from place to place. If you must say ``1 can of
soup'', then at least tell me how big you think a can of soup is.
.IP (d) 5
Numbered sequence of recipe steps. Be very careful to mention every
ingredient somewhere. The most common mistake made in recipes is to omit one
or more ingredients from the procedure steps. (The second-most-common mistake
is to leave an ingredient out of the ingredient list).
.IP (e) 5
Notes (if any). Comments on how you like to make it, ingredient
availability, comments about specific brands, etc.
.IP (f) 5
Your ``signature''. This should include your name and net address, the
organization that you are a part of, and the name of the city it is in.
.RE
.SH COPYRIGHT NOTES
Tell us where you got the recipe from. It's ok if you cribbed it from a book
or magazine or newspaper, but if you copy exactly the words that you found
there, there might be a problem with copyright violation. The copyright of a
recipe is not on the formula, but on the words. If you have copied the words
out of a copyrighted cookbook, then you are infringing its copyright.
.PP
While the main purpose of the USENET cookbook is to let us all make our own
custom cookbooks, we can't ignore the reality of the copyright law. Surely
you have noticed that every modern book says ``\fINo part of this publication
may be reproduced, stored in a retrieval system, or transmitted, in any form
or by any means, electronic, ..., or otherwise without the prior written
permission of the publisher.\fR''
.PP
Copyright law is complex, and only a lawyer who specializes in copyright
can reliably advise you on whether or not you are violating it, but in
general if you rewrite a recipe, in your own words, even if you don't change
the formula, then you are not infringing the copyright by submitting that
recipe to the network. The copyright is on the words that explain the recipe,
and on the title of the recipe, not the formula itself.
.PP
The USENET Cookbook itself is copyrighted. Every recipe is sent out with a
copyright notice, and the macros that print the cookbooks add a copyright
notice to the bottom of each page. The purpose of this is to prevent
unauthorized commercial use of the USENET Cookbook.
.SH DO YOUR OWN FORMATTING
If you put formatting commands in the recipes that you submit, they will go
out more quickly, since the cookbook editor won't have to put them in.
.PP
The basic principle is to use as few commands as possible, and in general to
use only commands that are defined in the manual macros or the recipe macro
package. Various processing programs search through these files and look for
string matches on things like ``.IG'' and ``.RZ''.
Print out a few of the recipe source files from NEWSDIR/alt/gourmand,
and then use OBJDIR/rctypeset
or OBJDIR/rcnroff to produce 
formatted versions of them. Have that at hand when you are formatting
your recipe. You can test your formatting by using rctypeset or rcnroff on
your draft versions before you submit the recipe.
.PP
A ``skeleton'' recipe follows. .RH is the recipe header command; it must be
the first line of each recipe and it must have 4 arguments. 
``RECIPE-ID'' is the file name under which the recipe will be stored. Normally
the editor chooses the recipe id, but you can put a suggestion there if you
like. It must be 14 characters or less, and must not duplicate any previous
or pending recipe id. The ``?'' is a ``what kind of recipe'' code from the
code table below.
.PP
.nf
.ta 3i
.RS 4
\&.RH MOD.RECIPES-SOURCE RECIPE-ID ? "22 Dec 83" 
\&.RZ "RECIPE TITLE IN CAPITALS" "One-line description of it"
Introductory comments; use .PP between paragraphs.
\&.IH "4 cups"	\(<- Ingredients Header
\&.IG "1/2 cup" "butter"	\(<- Ingredient (please use quotes)
\&.IG "1" "onion"
(medium to large, chopped fine. Don't try to use instant onion
in this recipe)
\&.PH	\(<- Procedure header
\&.SK 1	\(<- Procedure step
Boil the water.	\(<- Text for that step
\&.SK 2	\(<- and so forth.
\&.NX	\(<- Notes header
Notes (commentary) goes here; use .PP to separate paragraphs.
\&.WR	\(<- Wrapup
Signature information goes here. As a minimum you should list your name,
network address, organization (company, university, etc.), and the city you
live or work in.
.RE 4
.fi
Remember that when you post to a moderated newsgroup, the news software
usually forgets to include your .signature file, so you should be sure to
include it manually. If it manages to get there twice, I will remove the
extra copy.
.PP
You can also use the following -man macros; see \fIman(7)\fR for the complete
list.
.nf
.RS 4
.ta 3i
\&.I "italic words"
\&.B "boldface words"
\&.SM "small words"
\&.PP	\(<- paragraph break
\&.PD <distance>	\(<- paragraph distance
\&.IP "indented paragraph".
\&.RS	\(<- relative start: move things to the right
\&.RE	\(<- relative end: move things left again
.RE 4
.fi
You can also use these nroff/troff commands:
\\fI,
\\fB,
\\fR,
\&.if, 
\&.ds,
\&.br,
\&.nf,
\&.fi,
and .ta
.SH MET\&RIC OR INTUITIVE MEASUREMENT?
Most countries in the world use metric measurements for their recipes. The
U.S. defiantly uses a system that they call ``the English system of
measurement,'' teaspoons and cups. The English never used that system,
though\(emthey used a similar system, with spoons and cups of a different
size, and they called it the Imperial system. Every attempt to convert U.S.
cooks to the metric system has failed. Most Americans have never even seen a
metric recipe, and virtually no American cook owns a kitchen scale or metric
measuring spoons.
.PP
But alt.gourmand is international in scope, and we want the recipes to be
accessible to everyone. Therefore the recipes that are posted all have a
dual measurement system in them, both English and metric. When a cookbook
page is printed, a troff/nroff option determines whether or not it will be
printed with English units or metric units.
.PP
The dual-system measurement scheme is accomplished by having both sets of
values stored in the body of the recipe, and having the software select one
or the other when the recipe is printed. This means that somebody has
converted the recipe to both sets of units, and has edited in both sets of
measurements.
.PP
There are four places in a recipe where the system of measurement matters:
the ingredient list, commentary about ingredients, references to ingredients
in the text, and references to temperature in the text. There is a
troff/nroff command for each one of those situations.
Because alt.gourmand originated in the U.S., the U.S. (``English'')
measurement system is always given first in all of these commands.
.RS
\fBIngredient header.\fR  The ``.IH'' command takes two arguments. The first
is the English-unit specification of how much the recipe produces, and the
second (if present) is the metric specification of the same thing.
.PP
\fBIngredient list.\fR  The ``.IG'' command takes an optional third
argument, which is the metric description of the quantity of the ingredient.
For counts\(em3 onions, 2 eggs\(emthe metric description will be the same
and you can leave it off.
.PP
\fBReferences to ingredients.\fR  In the text you might want to say
something like ``set aside \(12 Tbsp of the ginger'' or ``reserve 100 ml of
the sauce''. For that purpose there is a macro ``.AB'', that takes two
arguments and prints one or the other, but never both. 
.PP
\fBReferences to temperature.\fR  The U.S. uses Fahrenheit degrees; nearly
everyone else uses Celsius or Centigrade. Troff can print a ``degree'' sign,
but nroff cannot. To solve these two problems simultaneously, there is a
``.TE'' macro, for indicating temperature. It takes two numeric arguments,
the first a Fahrenheit temperature and the second a Celsius temperature. 
.RE
In case the .AB or .TE expression needs to be the end of a sentence, the
macro can take a third argument, which is the punctuation character at
the end of the sentence. For example, you would type
.nf
.RS 4
Then add butter
\&.AB "1 Tbsp" "30 g"
at a time. Lick the spoon to use the last
\&.AB "tablespoon" "few grams" .
.fi
.PP
Here is the example of the previous section, updated to include
international multi-unit arguments to all of the macro calls:
.nf
.ta 3i
.RS 4
\&.RH MOD.RECIPES-SOURCE RECIPE-ID ? "22 Dec 83" 
\&.RZ "RECIPE TITLE IN CAPITALS" "One-line description of it"
Introductory comments; use .PP between paragraphs.
\&.IH "4 cups" "1 liter"	\(<- Ingredients Header
\&.IG "1/2 cup" "butter" "100 g"	\(<- Ingredient (please use quotes)
(or use a mixture of
\&.AB "1/4 cup" "50 g"	\(<- In-text reference to two systems
of margarine and
\&.AB "1/4 cup" "50 g"
of butter.
\&.IG "1" "onion"	\(<- No need for 2 systems here
(medium to large, chopped fine. Don't try to use instant onion
in this recipe)
\&.PH	\(<- Procedure header
\&.SK 1	\(<- Procedure step
Preheat the oven to	\(<- Text for that step
\&.TE 350 175	\(<- Dual-system temperature
before soaking the rice.
Boil the water.
\&.SK 2	\(<- and so forth.
\&.NX	\(<- Notes header
Notes (commentary) goes here; use .PP to separate paragraphs.
\&.WR	\(<- Wrapup
Signature information goes here. As a minimum you should list your name,
network address, organization (company, university, etc.), and the city 
and country you live or work in.
.RE 4
.fi
.SH CONVERTING RECIPES TO/FROM MET\&RIC
Don't try to convert a recipe to metric units unless you have some
experience cooking with metric, and don't try to convert a recipe to English
units unless you have some experience cooking with them. Submit your recipe
in the units that you are comfortable with, and let the alt.gourmand editor
do the conversion for you. It's not just a simple matter of unit conversion,
because most ingredients are specified by weight in metric recipes and
by volume in English recipes.
.SH CATEGORY CODES
.nf
.ta .4i 1i 3i 3.6i
	M	Main dish	SL	Salad
	A	Appetizer or snack	SP	Soup
	B	Bread/pasta	D	Dessert
	L	Beverage (Liquid)	V	Vegetable dish
	C	Cookie or cake	O	Other
	S	Sauce
.fi
The suffix ``V'' on any category means that it is vegetarian; for example, a
vegetarian main dish recipe would be marked ``MV''.
.SH SEE ALSO
cookbook(1), rn(1)
.SH AUTHOR
Brian Reid, DEC Western Research Laboratory
