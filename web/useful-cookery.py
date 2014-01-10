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

from bottle import run, debug as bottle_debug, route, redirect, abort, static_file, jinja2_view, jinja2_template as template

import jinja2

import recipes

DEBUG = 0


def main():
    """
    If run as a script, parse command line arguments and start embedded server
    """

    parser = ArgumentParser(description="A really interesting program")
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
        bottle_debug(True)

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


@filter
def interpret(aValue, aUnitsType='us'):
    ret = aValue

    global DEBUG

    if '{{' in ret:
        def chooseUnits(aUS, aMetric=None, anAddon=None):
            ret = aUS

            # todo: check a param, cookie, config setting or something

            if aMetric and (aUnitsType == 'metric'):
                ret = aMetric
            
            if anAddon:
                ret += anAddon

            return ret

        e = jinja2.Environment()
        t = e.from_string(ret)
        ret = t.render(chooseUnits=chooseUnits,
                       pre = '<pre>',
                       endpre = '</pre>'
            )

        if DEBUG >= 3: print('filter called, value="%s", units="%s" -> "%s"' % (aValue, aUnitsType, ret), file=sys.stderr)

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

    rname = 'CHOC-CHIP-1'
    rotd = recipes.get(rname) or { 'name': rname,
                                   'error': True,
                                   'errorText': 'Unable to find recipe',
                                   }

    ret['rotd'] = rotd

    return ret


@route('/index')
@view('permindex.tpl')
def home():
    """
    index page
    """
    ret = {
        'title': 'Useful Cookery Permuted Index - Recipes lovingly restored for the global village',
        }

    # todo: paging windows
    
    all = recipes.getSummary() or { 'error': True,
                                    'errorText': 'Unable to get recipe list',
                                    }

    ret['recipes'] = all

    return ret


@route('/search')
@view('search.tpl')
def home():
    """
    search page
    """
    ret = {
        'title': 'Search Useful Cookery - Recipes lovingly restored for the global village',
        }


    return ret


@route('/recipe/<aName>')
@view('recipe.tpl')
def home(aName):
    """
    recipe page
    """
    ret = {
        'rawname': aName,
        }

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
