#!/usr/bin/python3

##----------------------------------------------------------------------
## Copyright (c) 2014 Carl A. Dunham, All Rights Reserved
##----------------------------------------------------------------------
##
## loadfiles.py
##
## Created: 2014-Jan-07 by carl
##
##----------------------------------------------------------------------

"""
Convert troff file(s) and load into the database
"""

import sys
import os
import argparse

from urllib.parse import urlparse
from pymongo import MongoClient

from troff2json import readRecipe

DEBUG = 0


def main():
    """
    Given one or more files, load each one
    """

    parser = argparse.ArgumentParser(description='Load recipe files into a database')
    parser.add_argument('-d', '--debug', type=int, default=0, help='set debug level to DEBUG (default: %(default)s)')
    group = parser.add_mutually_exclusive_group()
    group.add_argument('-q', '--quiet', action='store_true')
    group.add_argument('-v', '--verbose', action='store_true')

    parser.add_argument('-u', '--url', type=str, default='mongodb://localhost:27017/', help='MongoDB URL (default: %(default)s)')
    parser.add_argument('-D', '--database', type=str, default='useful-cookery', help='MongoDB database name (default: %(default)s)')

    parser.add_argument('infiles', metavar='filename', nargs='*', type=str, default=['-'], help='file(s) to read from (default: <stdin> if missing or "-")')

    opts = parser.parse_args()

    global DEBUG
    DEBUG = opts.debug

    if DEBUG >= 3:
        print('opts="%s"' % opts, file=sys.stderr)

    url = urlparse(opts.url)

    assert url.scheme == 'mongodb', 'expected MongoDB URL'

    client = MongoClient(url.geturl())
    db = client[opts.database]

    for fname in opts.infiles:
        fd = sys.stdin if (not fname) or (fname == '-') else open(fname, 'r')

        recipe = readRecipe(fd)

        res = db.recipes.update({ 'name': recipe['name'] }, recipe, upsert=True)

        id = res['upserted'] if 'upserted' in res else res

        if DEBUG >= 3: print('inserted "%s" as record [%s]' % (fname, id), file=sys.stderr)
        

if __name__ == '__main__':
    main()
