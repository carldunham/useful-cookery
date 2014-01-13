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
import re
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


def getSummary(aSortKey=None):
    """
    Return a summary list of all the recipes
    """
    ret = _db.recipes.find({}, { 'name': 1, 'title': 1, 'description': 1 })

    if ret and aSortKey:
        ret.sort(aSortKey)

    return ret


def search(aTerm):
    ret = None

    # db.recipes.runCommand( "text", { search: "\"bacon\"", project: { "name": 1 }, language: "english" } )

    ret = _db.command('text', 'recipes', search=aTerm, project={ 'name': 1, 'title': 1, 'description': 1 }, language='english')

    return ret


def permutedIndex():
    ret = []

    STOP_WORDS = (
        'a', '(a', '(and', 'an', 'the', 'of', 'in', 'and', 'or', 'for', 'made', 
        'make', 'with', 'from', '-', '&', 'like', 'to', 'very', 'as', 'who', 
        'what', 'where', 'when', 'at', 'be', 'no', 'not', 'on', 'piece', 'pieces',
        'top', 'truly', 'two', 'used', 'you', 'ever', 
        )

    for recipe in getSummary():
        words = recipe['description'].split()

        for i in range(len(words)):

            if words[i].lower() not in STOP_WORDS:
                pre = ' '.join(words[:i])
                post = ' '.join(words[i+1:])

                ret.append((pre, words[i], post, recipe))

    def getkey(aWord):
        ret = re.sub(r'(?:&[lr]dquo;|[ \(\)\"\"])', '', aWord.lower())

        return ret
    
    ret.sort(key = lambda x: getkey(x[1]))
    
    return ret
