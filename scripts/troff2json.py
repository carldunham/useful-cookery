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
import re

from datetime import datetime

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


    recipe = readRecipe(opts.infile)

    if not opts.quiet:
        json.dump(recipe, sys.stdout, indent=2, separators=(',', ': '))
        print()


def readRecipe(anIterator):
    ret = {
        'name': '',
        'category': '',
        'title': '',
        'description': '',

        'sections': []
        }

    global DEBUG

    IGNORABLE_COMMANDS = ('ig', '.', 'S', 'sp')

    for block in _blockIter(anIterator):
        if DEBUG >= 4: print(block, file=sys.stderr)

        cmdline = block[0]
        cmd = cmdline[1:3]

        try:
            cmdargs = shlex.split(cmdline)[1:]
        except ValueError:
            if DEBUG >= 2: print('error parsing line [%s]' % cmdline, file=sys.stderr)

            # simple possible fix
            if not cmdline.endswith('"'):
                cmdargs = shlex.split(cmdline + '"')[1:]
            else:
                raise

        if cmd == 'RH':
            assert(len(cmdargs) == 5)

            src, name, category, date, year = cmdargs

            if DEBUG >= 3: print('.RH, name="%s", cat="%s", date="%s"' % (name, category, date), file=sys.stderr)

            # should only be one per file. if not, last in wins
            ret['name'] = name
            ret['category'] = category
            ret['date'] = datetime.strptime(date, '%d %b %y').strftime('%Y-%m-%d')

            if DEBUG >= 3: print("[%s] -> [%s]" % (date, ret['date']), file=sys.stderr)

            # assert(datetime.strptime(ret['date'], '%Y-%m-%d').strftime('%e %b %y').strip().lower() == date.lower())

        elif cmd == 'RZ':
            assert(len(cmdargs) == 2)

            title, desc = cmdargs

            if DEBUG >= 3: print('.RZ, title="%s", desc="%s"' % (title, desc), file=sys.stderr)

            # should only be one per file. if not, last in wins
            ret['title'] = title
            ret['description'] = desc

            if len(block) > 1:
                ret['introduction'] = [ _convertMultilineCodes(l) for l in _collapseLines(block[1:]) ]

        elif cmd == 'SH':
            sh = {}

            if len(cmdargs) > 0:
                sh['name'] = cmdargs[0]
            
            if len(block) > 1:
                sh['comments'] = [ _convertMultilineCodes(l) for l in _collapseLines(block[1:]) ]

            # may be at the start of an ingredients section, or a section of its own
            # not sure what 'within' an ingredients section would mean...
            #
            if ((len(ret['sections']) == 0) or
                (('ingredient_sets' in ret['sections'][-1]) and (len(ret['sections'][-1]['ingredient_sets'][-1]['ingredients']) > 0)) or
                ('name' in  ret['sections'][-1])
                ):

                ret['sections'].append(sh)
            else:
                #assert('name' not in  ret['sections'][-1])
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

            set = { 'ingredients': [] }
            
            if len(cmdargs) > 0:
                yld = { 'us': cmdargs[0].strip() }

                if len(cmdargs) > 1:
                    yld['metric'] = cmdargs[1].strip()

                set['yield'] = yld

            ret['sections'][-1]['ingredient_sets'].append(set)
            
        elif cmd == 'IG':
            assert(len(ret['sections']) > 0)
            assert(len(cmdargs) > 1)

            if  'ingredient_sets' not in ret['sections'][-1]:
                ret['sections'][-1]['ingredient_sets'] = [ { 'ingredients': [] } ]

            ig = { 'us': cmdargs[0].strip(), 'ingredient': cmdargs[1].strip() }

            if len(cmdargs) > 2:
                ig['metric'] = cmdargs[2].strip()

            if len(block) > 1:
                ig['comments'] = [ _convertMultilineCodes(l) for l in _collapseLines(block[1:]) ]

            ret['sections'][-1]['ingredient_sets'][-1]['ingredients'].append(ig)
        
        elif cmd == 'PH':
            # if working from a named section, want to append into that. Otherwise start a new one for the procedures
            #
            if len(ret['sections']) == 0:
                ret['sections'].append( {} )
                
            if 'procedure_sets' not in ret['sections'][-1]:
                ret['sections'][-1]['procedure_sets'] = []

            proc =  { 'procedures': [] } 

            if len(cmdargs) > 0:
                proc['name'] = cmdargs[0].strip()
            
            ret['sections'][-1]['procedure_sets'].append(proc)

        elif cmd == 'SK':
            assert(len(ret['sections']) > 0)
            assert(len(cmdargs) > 0)

            if  'procedure_sets' not in ret['sections'][-1]:
                ret['sections'][-1]['procedure_sets'] = [ { 'procedures': [] } ]

            step = { 'index': cmdargs[0].strip() }
        
            if len(block) > 1:
                step['comments'] = [ _convertMultilineCodes(l) for l in _collapseLines(block[1:]) ]

            ret['sections'][-1]['procedure_sets'][-1]['procedures'].append(step)
        
        elif cmd == 'NX':
            if len(block) > 1:
                ret['notes'] = [ _convertMultilineCodes(l) for l in _collapseLines(block[1:]) ]
        
        elif cmd == 'WR':
            if len(block) > 1:
                ret['footer'] = block[1:]
        
        elif cmd in IGNORABLE_COMMANDS:
            if DEBUG >= 3: print('ignoring cmd "%s"' % cmd, file=sys.stderr)

        else:
            if DEBUG >= 2: print('***unknown cmd*** "%s"' % cmd, file=sys.stderr)

    return ret

# "private" functions
#

def _blockIter(anIterator):
    """
    yield chunks of troff blocks as tuples of (command line, text)
    """
    ret = []

    CONVERSION_CODES = ('.TE', '.AB')
    NON_BREAKING_CODES = ('.SM', '.PP', '.PD', '.RS', '.RE', '.if', '.ds', '.br', '.nf', '.fi', '.ta', '..') # that last one allows for old-style email paths
    IGNORE_CODES = ('.ta', ) # for one-liners. multi-line ones (like .ig) are dealt with as ignored commands above

    FORMAT_CODE_MAP = {
        'B': ('b',),
        'I': ('i',),
        'IR': ('i',),
        'BI': ('b', 'i'),
        'IB': ('i', 'b'),
        'SM': ('small',),
        'IP': ('li',),
        }

    nextLineRE = re.compile(r'^\.(%s)(\s"?(.*?)"?)?\s*$' % '|'.join(FORMAT_CODE_MAP.keys()))

    nextop = None

    def _convertFormat(aString, anOp):
        return ''.join([ '<%s>' % c for c in FORMAT_CODE_MAP[anOp] ]) + aString + ''.join([ '</%s>' % c for c in reversed(FORMAT_CODE_MAP[anOp]) ])

    for rawline in anIterator:
        line = _convertCodes(rawline.strip())

        if line.startswith('.'):
            m = nextLineRE.match(line)

            if line.startswith(IGNORE_CODES):
                if DEBUG >= 3: print('*** ignoring line: "%s" ***' % line, file=sys.stderr)
                line = None
                
            elif line.startswith(CONVERSION_CODES):
                parts = shlex.split(line)

                # us units
                params = '"%s"' % parts[1].strip()

                if len(parts) > 2:
                    # optional metric
                    params += ', "%s"' % parts[2].strip()

                line = '{{ chooseUnits(%s) }}' % params

                if len(parts) > 3:
                    # optional appended text
                    line += parts[3]

            elif m is not None:
                op = m.group(1)
                text = m.group(3)

                if text:
                    line = _convertFormat(text, op)

                else:
                    # todo
                    if DEBUG >= 3: print('*** handling next-line format for .%s ***' % op, file=sys.stderr)
                    nextop = op
                    line = None
                
            elif not line.startswith(NON_BREAKING_CODES):

                if ret:
                    yield ret
                    ret = []
                    nextop = None

            else:
                if DEBUG >= 3: print('*** passing through line [%s] ***' % line, file=sys.stderr)
                    
        if line is not None:

            if nextop is not None:
                line = _convertFormat(line, nextop)
                if DEBUG >= 3: print('***    emitting "%s" ***' % line, file=sys.stderr)
                nextop = None
                
            ret.append(line)

    if ret:
        yield ret
            

def _convertCodes(aString):
    ret = _convertMultilineCodes(aString)

    REPLACEMENTS = {
        r'\\\(18': ' 1/8',
        r'\\\(14': ' 1/4',  # '&frac14'
        r'\\\(13': ' 1/3',
        r'\\\(12': ' 1/2',  # '&frac12'
        r'\\\(34': ' 3/4',  # '&frac34'
        r'\\\(mu': '&times;',
        r'\\\(em': '&em;',
        r'\\-': '-',       # &em;
        #r'^.I (.+)$': '<i>\\1</i>',
        #r'^.B (.+)$': '<b>\\1</b>',
        #r'^.SM (.+)$': '<small>\\1</small>',
        r'\\z\\\(aae': '&eacute;',
        r'\\z\\\(aao': '&oacute;',
        r'\\z\\\(gaa': '&agrave;',
        r"\\o'e\\\(aa": '&eacute;', # this shows up in a couple of recipes
        r'\\\*\:o': '&ouml;',
        r"\\o'o\"'": '&ouml;',
        r'``': '&ldquo;',
        r"''": '&rdquo;',
        r'\t': ' ',    # '&nbsp;',
        r'\\s\-2': '<span class="troff-s-2">',
        r'\\s0': '</span>',
        r'^\.PP\s*': '<p>',
        r'^\.br\s*': '<br>',
        r'^\.nf\s*': '{{pre}}',  # allows for more reasonable options than just <pre>, ie css white-space: pre-line
        r'^\.fi\s*': '{{/pre}}',
        r'^\.RS\s*': '<div class="troff-RS">',
        r'^\.RE\s*': '</div>',
        }

    for k,v in REPLACEMENTS.items():
        ret = re.sub(k, v, ret)

    ret = ret.strip()
    
    return ret


def _convertMultilineCodes(aString):
    ret = aString

    REPLACEMENTS = {
        r'\\fI(.+?)\\f[PR](?m)': '<i>\\1</i>',
        r'\\fB(.+?)\\f[PR](?m)': '<b>\\1</b>',
        }

    for k,v in REPLACEMENTS.items():
        ret = re.sub(k, v, ret)

    ret = ret.strip()
    
    return ret


def _collapseLines(aLineList):
    ret = []

    global DEBUG

    linebuf = ''
    pre = False

    for nextline in aLineList:

        if nextline == '{{pre}}':
            if DEBUG >= 2: print('---> pre', file=sys.stderr)
            pre = True

            if linebuf:
                ret.append(linebuf)
                linebuf = ''

        elif nextline == '{{/pre}}':
            if DEBUG >= 2: print('<--- pre', file=sys.stderr)
            pre = False

            if linebuf:
                ret.append(linebuf)
                linebuf = ''

        if pre:
            ret.append(nextline)
        elif not linebuf:
            linebuf = nextline
        else:
            linebuf += (' ' + nextline) # this leaves the {{/pre}} appended with following text, but that should be ok

    if linebuf:
        ret.append(linebuf)

    return ret

    
if __name__ == '__main__':
    main()
