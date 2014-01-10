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

DEBUG = 0


def setDebug(aDebugLevel=1):
    global DEBUG
    DEBUG = aDebugLevel


def get(aName):
    """
    Return a recipe with a given name, or None if it doesn't exist
    """
    client = MongoClient()

    db = client['useful-cookery']

    recipe = db.recipes.find_one({ 'name': aName })

    return recipe
