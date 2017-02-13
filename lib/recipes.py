
# ----------------------------------------------------------------------
#  Copyright (c) 2014-2017 Carl A. Dunham, All Rights Reserved
# ----------------------------------------------------------------------
#
#  recipes.py
#
#  Created: 2014-Jan-09 by carl
#
# ----------------------------------------------------------------------

"""
functions for managing recipes
"""

import re
import shlex
import random

from pymongo import MongoClient


CATEGORIES = {
    'M': 'Main Dish',
    'MV': 'Vegetarian Main Dish',
    'A': 'Appetizer/Snack',
    'AV': 'Vegetarian Appetizer/Snack',
    'B': 'Bread/Pasta',
    'L': 'Beverage',
    'C': 'Cookie/Cake',
    'S': 'Sauce',
    'SL': 'Salad',
    'SP': 'Soup',
    'SPV': 'Vegetarian Soup',
    'D': 'Dessert',
    'V': 'Vegetable',
    'O': 'Other',
}


_client = MongoClient()

_db = _client['useful-cookery']

DEBUG = 0


def setDebug(aDebugLevel=1):
    global DEBUG
    DEBUG = aDebugLevel


def get(aName):
    """
    Return a recipe with a given name, or None if it doesn't exist
    """
    ret = _db.recipes.find_one({'name': aName})

    return ret


def getRandom(aCategory=None):
    """
    Return a random recipe name
    """
    ret = None

    cond = {'category': aCategory} if aCategory else {}

    n = _db.recipes.find(cond).count()

    i = random.randint(0, n-1)

    ret = _db.recipes.find_one(cond, {'name': 1}, skip=i)

    return ret


def getSummary(aCategory=None, aSortKey=None):
    """
    Return a summary list of all the recipes
    """
    ret = None

    cond = {'category': aCategory} if aCategory else {}

    ret = _db.recipes.find(cond, {'name': 1, 'title': 1, 'description': 1})

    if ret and aSortKey:
        ret.sort(aSortKey)

    return ret


def search(search_input):
    ret = None

    terms = [('"%s"' % t) for t in shlex.split(search_input)]

    ret = _db.recipes.find({
        "$text": {
            "$search": ''.join(terms),
        },
    }, {
        "name": 1, "title": 1, "description": 1,
        "score": {"$meta": "textScore"},
    }).sort([
        ('score', {'$meta': 'textScore'}),
    ])

    return ret


def permutedIndex():
    ret = []

    STOP_WORDS = (
        'a', '(a', '(and', 'an', 'the', 'of', 'in', 'and', 'or', 'for', 'made',
        'make', 'with', 'from', '-', '&', 'like', 'to', 'very', 'as', 'who',
        'what', 'where', 'when', 'at', 'be', 'no', 'not', 'on', 'piece', 'pieces',
        'top', 'truly', 'two', 'used', 'you', 'ever', 'recipe',
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

    ret.sort(key=lambda x: getkey(x[1]))

    return ret
