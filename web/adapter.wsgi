import sys
import os
import bottle

pwd = os.path.dirname(__file__)
sys.path = [pwd, os.path.join(os.path.dirname(pwd), 'lib')] + sys.path

import usefulcookery  # noqa

os.chdir(os.path.dirname(__file__))

application = bottle.default_app()
