import sys, os, bottle

sys.path = ['/var/www/useful-cookery.com/web/'] + sys.path
os.chdir(os.path.dirname(__file__))

import useful-cookery

application = bottle.default_app()
