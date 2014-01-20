#!/usr/bin/python3

##----------------------------------------------------------------------
## Copyright (c) 2014 Carl A. Dunham, All Rights Reserved
##----------------------------------------------------------------------
##
## mongoutil.py
##
## Created: 2014-Jan-19 by carl
##
##----------------------------------------------------------------------

"""
Some helpful code for working with MongoDB
"""

import sys
import shlex


DEBUG = 0

def setDebug(aDebugLevel=1):
    global DEBUG
    DEBUG = aDebugLevel


def andSearch(aDB, aCollection, aQuery, aFilter={}, aProjection={}, aLanguage='english'):
    ret = None

    terms = [ ('"%s"' % t) if (' ' in t) else t for t in shlex.split(aQuery) ]

    if DEBUG >= 4: print('query="%s", terms=%s' % (aQuery, terms), file=sys.stderr)

    # can't support limit (well, maybe todo: get all and truncate sorted results, but why?)
    #
    ret = aDB.command('text', aCollection, search=terms[0], filter=aFilter, project=aProjection, language=aLanguage)

    if len(terms) > 1:
        possibles = { r['obj']['_id']:r for r in ret['results'] }

        realterms = [ terms[0] ] if ret['queryDebugString'].split('|')[0] else []

        for term in terms[1:]:
            next = aDB.command('text', aCollection, search=term, filter=aFilter, project=aProjection, language=aLanguage)

            for key in next.keys():

                # keep stats for all terms, even stop words
                #
                if key == 'stats':
                    ret[key]['timeMicros'] += next[key]['timeMicros']
                    ret[key]['nfound'] += next[key]['nfound']
                    # TODO: more

                # but only keep results for non-stopword terms
                #
                elif (key == 'results') and next['queryDebugString'].split('|')[0]:

                    if not realterms:
                        possibles = { r['obj']['_id']:r for r in next['results'] }
                    else:
                        possibles = { r['obj']['_id']:r for r in next['results'] if r['obj']['_id'] in possibles }

                    realterms.append(term)

        # this sort isn't really accurate, as it only uses the last search score
        #
        ret['results'] = sorted(possibles.values(), key=lambda x: x['score'])

        ret['stats']['nfound'] = len(ret['results'])

        if DEBUG >= 4: print('realterms=%s, ret=%s' % (realterms, ret), file=sys.stderr)
        

    return ret
    
