#!/usr/bin/python3

##----------------------------------------------------------------------
## Copyright (c) 2014 Carl A. Dunham, All Rights Reserved
##----------------------------------------------------------------------
##
## troff2json.py
##
## Created: 2014-Jan-02 by carl
##
##----------------------------------------------------------------------

"""
Read a troff (man page) version of a recipe and output in JSON
"""

import sys
import os
import argparse

import json
import shlex

DEBUG = 0


def main():
    """
    This is the main function for the module. 
    """

    parser = argparse.ArgumentParser(description='Read a troff-formatted recipe and output in JSON format')
    parser.add_argument('-d', '--debug', type=int, default=0, help='set debug level to DEBUG (default: %(default)s)')
    group = parser.add_mutually_exclusive_group()
    group.add_argument('-q', '--quiet', action='store_true')
    group.add_argument('-v', '--verbose', action='store_true')

    parser.add_argument('infile', metavar='filename', nargs='?', default=sys.stdin, type=argparse.FileType('r'), help='file to read from (<stdin> if missing or "-")')
    #parser.add_argument('-o', '--outfile', default=sys.stdout, type=argparse.FileType('w'), help='file to write to (<stdout> if missing or "-")')

    opts = parser.parse_args()

    global DEBUG
    DEBUG = opts.debug

    if DEBUG >= 3:
        print('opts="%s"' % opts, file=sys.stderr)


    json.dump(readRecipe(opts.infile), sys.stdout, indent=2, separators=(',', ': '))
    print()


def readRecipe(anFP):
    ret = {
        'name': '',
        'category': '',
        'date': '',
        'title': '',
        'description': '',

        'sections': []
        }

    global DEBUG

    for block in _blockIter(anFP):
        if DEBUG >= 4: print(block, file=sys.stderr)

        cmdline = block[0]
        cmd = cmdline[1:3]

        cmdargs = shlex.split(cmdline)[1:]

        if cmd == 'RH':
            assert(len(cmdargs) == 5)

            src, name, category, date, year = cmdargs

            if DEBUG >= 3: print('.RH, name="%s", cat="%s", date="%s"' % (name, category, date), file=sys.stderr)

            # should only be one per file. if not, last in wins
            ret['name'] = name
            ret['category'] = category
            ret['date'] = date

        elif cmd == 'RZ':
            assert(len(cmdargs) == 2)

            title, desc = cmdargs

            intro = " ".join(block[1:]) if (len(block) > 1) else ''

            if DEBUG >= 3: print('.RZ, title="%s", desc="%s", intro="%s"' % (title, desc, intro), file=sys.stderr)

            # should only be one per file. if not, last in wins
            ret['title'] = title
            ret['description'] = desc

            if intro:
                ret['introduction'] = intro

        elif cmd == 'SH':
            sh = { 'name': cmdargs[0] if len(cmdargs) > 0 else '' }
            
            if len(block) > 1:
                sh['comments'] = block[1:]

            # may be at the start of an ingredients section, or a section of its own
            # not sure what 'within' an ingredients section would mean...
            #
            if len(ret['sections'][-1]['ingredient_sets'][-1]['ingredients']) > 0:
                ret['sections'].append(sh)
            else:
                assert('name' not in  ret['sections'][-1])
                assert('comments' not in ret['sections'][-1])

                ret['sections'][-1]['name'] = sh['name']

                if 'comments' in sh:
                    ret['sections'][-1]['comments'] = sh['comments']
                
        elif cmd == 'IH':
            # some recipes are broken into sections, others are not. So if we come across an ingredients header,
            # and have no sections started, start one. Otherwise, append to the last one.
            #
            if len(ret['sections']) == 0:
                ret['sections'].append( {} )
                
            if 'ingredient_sets' not in ret['sections'][-1]:
                ret['sections'][-1]['ingredient_sets'] = []

            yld = { 'us': cmdargs[0] }

            if len(cmdargs) > 1:
                yld['metric'] = cmdargs[1]

            ret['sections'][-1]['ingredient_sets'].append( { 'yield': yld, 'ingredients': [] } )
            
        elif cmd == 'IG':
            assert(len(ret['sections']) > 0)
            assert(len(cmdargs) > 1)

            if  'ingredient_sets' not in ret['sections'][-1]:
                ret['sections'][-1]['ingredient_sets'] = [ { 'ingredients': [] } ]

            ig = { 'us': cmdargs[0], 'ingredient': cmdargs[1] }

            if len(cmdargs) > 2:
                ig['metric'] = cmdargs[2]

            if len(block) > 1:
                ig['comments'] = block[1:]

            ret['sections'][-1]['ingredient_sets'][-1]['ingredients'].append(ig)
        
        elif cmd == 'PH':
            # if working from a named section, want to append into that. Otherwise start a new one for the procedures
            #
            if len(ret['sections']) == 0:
                ret['sections'].append( {} )
                
            if 'procedure_sets' not in ret['sections'][-1]:
                ret['sections'][-1]['procedure_sets'] = []

            name = cmdargs[0] if len(cmdargs) > 0 else ''
            
            ret['sections'][-1]['procedure_sets'].append( { 'name': name, 'procedures': [] } )

        elif cmd == 'SK':
            assert(len(ret['sections']) > 0)
            assert('procedure_sets' in ret['sections'][-1])
            assert(len(cmdargs) > 0)

            step = { 'index': cmdargs[0] }
        
            if len(block) > 1:
                step['comments'] = block[1:]

            ret['sections'][-1]['procedure_sets'][-1]['procedures'].append(step)
        
        elif cmd == 'NX':
            if len(block) > 1:
                ret['notes'] = block[1:]
        
        elif cmd == 'WR':
            if len(block) > 1:
                ret['footer'] = block[1:]
        
        else:
            if DEBUG >= 3: print('skipping cmd="%s"' % cmd, file=sys.stderr)

    return ret

# "private" functions
#

def _blockIter(anIterator):
    """
    yield chunks of troff blocks as tuples of (command line, text)
    """
    ret = []

    INLINE_FORMAT_CODES = ('.TE', '.AB', '.I ', '.B ', '.SM', '.PP', '.PD', '.IP', '.RS', '.RE', '.if', '.ds', '.br', '.nf', '.fi', '.ta')

    for rawline in anIterator:
        line = rawline.strip()

        if line.startswith('.') and not line.startswith(INLINE_FORMAT_CODES):

            if ret:
                yield ret
                ret = []

        ret.append(line)

    if ret:
        yield ret
            

    
if __name__ == '__main__':
    main()
