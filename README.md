useful-cookery
==============
Copyright &copy; 2014-2017 Carl A. Dunham, All Rights Reserved

Installing
==========

Install `mongodb` and start it

In your favorite virtual environment:
```
 % pip install -r requirements.txt
 % mongo < config/mongodb/setup.mongodb
 % scripts/loadfiles.py --source usenet --copyright usenet
 % scripts/loadfiles.py --source usenet-published --copyright usenet data/raw/recipes/*
 % scripts/loadfiles.py --source usenet-unpublished --copyright usenet data/raw/recipes-unpub/*
 % mongo useful-cookery < config/mongodb/recipe-images.mongodb
```
A sample Apache2 config file is provided.

Development
===========

To run the code locally, additional set up is needed:
```
 % pip install -r requirements-dev.txt
```
The built-in Bottle server can be run to do testing:
```
 % cd web/
 % PYTHONPATH=.:../lib ./usefulcookery.py -d 3
```
The `--debug` flag will enable additional logging both from the usefulcookery app and Bottle.
