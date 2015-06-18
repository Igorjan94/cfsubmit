from selenium import webdriver
import time

from flask import Flask
from flask import g

app = Flask(__name__)

# initialize global variables
ctx = app.app_context()
ctx.push()

g.browser = webdriver.Chrome()
time.sleep(2)

g.handle = ''
g.password = ''

g.default_url = 'http://codeforces.com'
g.login_url = g.default_url + '/enter'
g.submit_url = ""

g.browser.get(g.default_url)

import valahd_srv.views
