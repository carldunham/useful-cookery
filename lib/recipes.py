#!/usr/bin/python3

##----------------------------------------------------------------------
## Copyright (c) 2014 Carl A. Dunham, All Rights Reserved
##----------------------------------------------------------------------
##
## recipes.py
##
## Created: 2014-Jan-09 by carl
##
##----------------------------------------------------------------------

"""
functions for managing recipes
"""

import sys
import os
from argparse import ArgumentParser
from pymongo import MongoClient
from jinja2 import evalcontextfilter

DEBUG = 0

_client = MongoClient()

_db = _client['useful-cookery']


def setDebug(aDebugLevel=1):
    global DEBUG
    DEBUG = aDebugLevel


def get(aName):
    """
    Return a recipe with a given name, or None if it doesn't exist
    """
    ret = _db.recipes.find_one({ 'name': aName })

    return ret


def getSummary():
    """
    Return a summary list of all the recipes
    """
    ret = _db.recipes.find({}, { 'name': 1, 'title': 1, 'description': 1 }).sort('title')

    return ret


