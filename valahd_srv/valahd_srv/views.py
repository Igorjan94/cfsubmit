from selenium import webdriver
from selenium.webdriver.support.ui import Select
from selenium.webdriver.common.keys import Keys

from datetime import datetime
import time

from flask import request
from flask import g

from valahd_srv import app

def login():
    g.browser.get(g.default_url)
    g.browser.get(g.login_url)
    time.sleep(1)
    # if we're already logged in we'll be redirected to default_url
    if g.browser.current_url == g.login_url:
        g.browser.find_element_by_id("handle").send_keys(g.handle)
        g.browser.find_element_by_id("password").send_keys(g.password)
        g.browser.find_element_by_id("remember").click()
        g.browser.find_element_by_class_name("submit").submit()
    

@app.route('/init_contest')
def init_contest():
    g.submit_url = 'http://codeforces.com/contest/{}/submit'.format(request.form['num'])
    return "contest {} initialized".format(g.submit_url)

@app.route('/init_server')
def init_server():
    g.handle = request.form['handle']
    g.password = request.form['password']
    login()
    return "server initialized"

@app.route('/submit')
def submit():
    # required form fields:
    # - id
    # - lang
    # - text
    if g.browser.current_url != g.submit_url:
        g.browser.get(g.submit_url)
        time.sleep(1)
    # select problem index
    select = Select(g.browser.find_element_by_name('submittedProblemIndex'))
    select.select_by_value(request.form['id'])

    # select language
    select = Select(g.browser.find_element_by_name('programTypeId'))
    select.select_by_value(request.form['lang'])
    
    # paste source
    # note this javascript requires _newest_ browser
    # because of String.raw function
    g.browser.execute_script("editAreaLoader.setValue('sourceCodeTextarea', String.raw`{}`)".format(text))
    
    # submit
    g.browser.find_element_by_class_name('submit').submit()
    time.sleep(1)

    # go back
    g.browser.get(g.submit_url)

    return "solution sent"


@app.route('/check_vars')
def home():
    return "handle: {}<br> password: {}<br><br> submit url: {}".format(
        g.handle,
        g.password,
        g.submit_url
    )