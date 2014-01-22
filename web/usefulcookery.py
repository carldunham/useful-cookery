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

    parser = ArgumentParser(description="Serves up the useful-cookery.com website")
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

    if code in recipes.CATEGORIES:
        ret = recipes.CATEGORIES[code] 
    
    return ret


def convertToISO8601(aTime):
    ret = 'PT1H'

    return ret


@filter
def interpret(aValue):
    ret = aValue

    global DEBUG

    if '{{' in ret:
        t = jinja2.Template(ret)

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
        'title': 'Useful Cookery - Recipes restored for the global village',
        }

    setcookies()

    rname = 'CHOC-CHIP-1'
    rotd = recipes.get(rname) 

    if not rotd:
        ret.update({ 'name': rname,
                     'error': True,
                     'errorText': 'Unable to find recipe',
                     })
    else:
        ret['rotd'] = rotd

    return ret


@route('/index')
@view('permindex.tpl')
def index():
    """
    index page
    """
    ret = {
        'title': 'Useful Cookery Permuted Index - Recipes restored for the global village',
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
        'title': 'Search Useful Cookery - Recipes restored for the global village',
        }

    setcookies()

    query = request.query.q

    if query:
        if DEBUG >= 4: print('search for "%s"' % (query), file=sys.stderr)

        ret['q'] = query

        results = recipes.search(query) or { 'error': True,
                                             'errorText': 'Unable to get search results',
                                             }

        if DEBUG >= 4: print('search for "%s", results=[%s]' % (query, results), file=sys.stderr)

        ret['results'] = results

    return ret


@route('/categories')
@view('categories.tpl')
def category():
    """
    listing of recipes categories
    """
    ret = {
        'title': 'Useful Cookery Categories - Recipes restored for the global village',
        }

    categories = recipes.CATEGORIES

    if not categories:
        ret.update({ 'error': True,
                     'errorText': 'Unable to get categories',
                     })
    else:
        ret['categories'] = categories

    return ret


@route('/category/<aCategory>')
#@route('/category/<aCategory>/<aSortKey>')
@view('category.tpl')
def category(aCategory, aSortKey='title'):
    """
    listing of recipes in a given category
    """
    ret = {
        'title': 'Useful Cookery Recipes by Category - Recipes restored for the global village',
        'category': {
            'key': aCategory,
            }
        }

    if aCategory in recipes.CATEGORIES.keys():
        ret['recipes'] = recipes.getSummary(aCategory=aCategory, aSortKey=aSortKey)
        ret['category']['name'] = getCategory(aCategory)

    else:
        abort(404, 'File does not exist.')
        
    return ret


@route('/random')
@route('/random/<aCategory>')
def randRecipe(aCategory=None):
    """
    Just redirect to some random recipe
    """

    r = recipes.getRandom(aCategory)

    if r and r['name']:
        redirect('/recipe/%s' % r['name'])
    else:
        abort(404, 'File does not exist.')

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

    r = recipes.get(aName.strip().upper()) 

    if r:
        ret['title'] = (r['title'] or r['description'] or  '') + ' Recipe - Useful Cookery'
        ret['recipe'] = r

    else:
        #ret['error'] = True
        #ret['errorText'] = 'Unable to find a recipe named "%s"' % aName
        abort(404, 'File does not exist.')

    return ret


@error(404)
@view('error404.tpl')
def error404(anError):
    ret = {
        'title': 'Useful Cookery - File Not Found',
        'statusText': anError.status_line,
        }

    return ret


@route('/favicon.ico')
def favicon():
    #print('looking for favicon', file=sys.stderr)
    return static_file('images/favicon.ico', root='static')
    #redirect('/images/favicon.ico')


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
