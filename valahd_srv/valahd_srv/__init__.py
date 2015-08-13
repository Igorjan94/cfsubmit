from selenium import webdriver
from flask import Flask
from flask import g

from urls import CF_MAIN_URL, CF_LOGIN_URL

app = Flask(__name__)

# initialize global variables
ctx = app.app_context()
ctx.push()

# g.browser = webdriver.Chrome()
g.browser = webdriver.PhantomJS()
g.browser.set_window_size(10200, 10200)

g.handle = ''
g.password = ''

g.default_url = CF_MAIN_URL
g.login_url = CF_LOGIN_URL
g.submit_url = ''

g.browser.get(g.default_url)

