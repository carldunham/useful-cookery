#!/usr/bin/python3

##----------------------------------------------------------------------
## Copyright (c) 2014 Carl A. Dunham, All Rights Reserved
##----------------------------------------------------------------------
##
## useful-cookery.py
##
## Created: 2014-Jan-09 by carl
##
##----------------------------------------------------------------------

"""
Main bottle app for usefulcookery.com
"""

import sys
import os
from argparse import ArgumentParser
import functools

from bottle import *
import jinja2

import recipes


DEBUG = 0

_unitsType = 'us' # default


def main():
    """
    If run as a script, parse command line arguments and start embedded server
    """

    parser = ArgumentParser(description="Server for useful cookery website")
    parser.add_argument("-d", "--debug", type=int, default=0, help="set debug level to DEBUG (default: %(default)s)")
    group = parser.add_mutually_exclusive_group()
    group.add_argument("-q", "--quiet", action="store_true")
    group.add_argument("-v", "--verbose", action="store_true")

    opts = parser.parse_args()

    global DEBUG
    DEBUG = opts.debug

    if DEBUG >= 3:
        print('opts="%s"' % opts, file=sys.stderr)

    if DEBUG > 0:
        debug(True)

    recipes.setDebug(DEBUG)
    
    run(reloader=(DEBUG>0))


# set up some basic stuff
#

# from http://jpvanoosten.nl/blog/2013/03/24/custom-jinja2-filters-when-using-bottle/

filter_dict = {}
view = functools.partial(jinja2_view, template_settings={'filters': filter_dict})

def filter(func):
	"""Decorator to add the function to filter_dict"""
	filter_dict[func.__name__] = func
	return func


def setcookies():
    """Check for setting of cookies, eg us vs. metric. should be done for each page"""

    global _unitsType
    
    if 'unitstype' in request.query:
        _unitsType = request.query['unitstype']
    elif 'unitstype' in request.cookies:
        _unitsType = request.cookies['unitstype']

    response.set_cookie('unitstype', _unitsType, path='/', max_age=1000*24*60*60)   # a long time
    
    
def chooseUnits(aUS, aMetric=None, anAddon=None):
    ret = aUS

    if aMetric and (_unitsType == 'metric'):
        ret = aMetric
    
    if anAddon:
        ret += anAddon

    return ret


def getCategory(aCategoryCode):
    ret = 'Unkown'

    code = aCategoryCode.upper()
    prefix = ''

    if (len(code) > 1) and code.endswith('V'):
        code = code[:-1]
        prefix = 'Vegetarian '

    if code in recipes.CATEGORIES:
        ret = prefix + recipes.CATEGORIES[code] 
    
    return ret


def convertToISO8601(aTime):
    ret = 'PT1H'

    return ret


@filter
def interpret(aValue):
    ret = aValue

    global DEBUG

    if '{{' in ret:
        e = jinja2.Environment()
        t = e.from_string(ret)
        ret = t.render(chooseUnits=chooseUnits,
                       pre = '<pre>',
                       endpre = '</pre>'
            )

        if DEBUG >= 4: print('filter called, value="%s", units="%s" -> "%s"' % (aValue, _unitsType, ret), file=sys.stderr)

    return ret


@route('/<:re:(index|home|default)\.html?>')
def sendhome():
    """
    redirect for home page aliases
    """
    redirect('/')


@route('/')
@view('index.tpl')
def home():
    """
    home page
    """
    ret = {
        'title': 'Useful Cookery - Recipes lovingly restored for the global village',
        }

    setcookies()

    rname = 'CHOC-CHIP-1'
    rotd = recipes.get(rname) or { 'name': rname,
                                   'error': True,
                                   'errorText': 'Unable to find recipe',
                                   }

    ret['rotd'] = rotd

    return ret


@route('/index')
@view('permindex.tpl')
def index():
    """
    index page
    """
    ret = {
        'title': 'Useful Cookery Permuted Index - Recipes lovingly restored for the global village',
        }

    setcookies()

    all = recipes.permutedIndex() or { 'error': True,
                                       'errorText': 'Unable to get permuted index',
                                       }

    ret['recipes'] = all

    return ret


@route('/search')
@view('search.tpl')
def search():
    """
    search page
    """
    ret = {
        'title': 'Search Useful Cookery - Recipes lovingly restored for the global village',
        }

    setcookies()

    query = request.query.q

    if query:
        if DEBUG >= 4: print('search for "%s"' % (query), file=sys.stderr)

        ret['q'] = query

        results = recipes.search(query) or { 'error': True,
                                             'errorText': 'Unable to get search results',
                                             }

        ret['results'] = results

    return ret


@route('/recipe/<aName>')
@view('recipe.tpl')
def recipe(aName):
    """
    recipe page
    """
    ret = {
        'rawname': aName,
        'chooseUnits': chooseUnits,
        'getCategory': getCategory,
        'convertToISO8601': convertToISO8601,
        }

    setcookies()

    recipe = recipes.get(aName.strip().upper()) 

    if recipe:
        ret['title'] = (recipe['title'] or recipe['description'] or  '') + ' Recipe - Useful Cookery'
        ret['recipe'] = recipe

    else:
        #ret['error'] = True
        #ret['errorText'] = 'Unable to find a recipe named "%s"' % aName
        abort(404, 'File does not exist.')

    return ret


# put this last, routes are processed in order
#
@route('/<aPath:path>')
def staticfiles(aPath):
    """
    handle static content
    """
    return static_file(aPath, root='static')


# if run as a script...
#
if __name__ == "__main__":
    main()
