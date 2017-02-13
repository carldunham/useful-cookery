import sys
import os
import bottle

sys.path = ['/var/www/useful-cookery.com/web/'] + sys.path

import usefulcookery  # noqa

os.chdir(os.path.dirname(__file__))

application = bottle.default_app()
