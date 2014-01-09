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

from bottle import run, debug as bottle_debug, route, redirect, abort, static_file, jinja2_view as view, jinja2_template as template

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
    
    run(reloader=True)


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
