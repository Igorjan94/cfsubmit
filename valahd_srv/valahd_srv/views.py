import time
from urllib.parse import urljoin

from selenium.webdriver.support.ui import Select
from flask import request
from flask import g

from urls import INIT_SERVER_PART, INIT_CONTEST_PART, SUBMIT_PART, CHECK_VARS_PART
from valahd_srv.valahd_srv import app


def sleep():
    time.sleep(1)


def login():
    if g.browser.current_url == g.submit_url:
        return
    if g.browser.current_url != g.login_url:
        g.browser.get(g.login_url)
        sleep()
    g.browser.find_element_by_id("handle").send_keys(g.handle)
    g.browser.find_element_by_id("password").send_keys(g.password)
    g.browser.find_element_by_id("remember").click()
    sleep()
    g.browser.find_element_by_class_name("submit").submit()
    sleep()


@app.route(urljoin('/', INIT_SERVER_PART))
def init_contest():
    g.submit_url = urljoin(g.default_url, 'contest/{}/submit'.format(request.form['num']))
    return "contest {} initialized".format(g.submit_url)


@app.route(urljoin('/', INIT_CONTEST_PART))
def init_server():
    g.handle = request.form['handle']
    g.password = request.form['password']
    login()
    return "server initialized"


def cf_submit(id, lang, text):
    # select problem index
    Select(g.browser.find_element_by_name('submittedProblemIndex')).select_by_value(id)

    # select language
    Select(g.browser.find_element_by_name('programTypeId')).select_by_value(lang)

    # how we can add text to textarea
    # requires modern browser cause of String.raw js function
    # g.browser.execute_script(
    #     "editAreaLoader.setValue('sourceCodeTextarea', String.raw`{}`)".format(request.form['text']))

    # add submission text - simpler method
    g.browser.find_element_by_name('source').send_keys(text)

    # send it
    g.browser.find_element_by_class_name('submit').submit()
    sleep()


@app.route(urljoin('/', SUBMIT_PART), methods=['POST'])
def submit():
    # required form fields:
    # - id
    # - lang
    # - text

    # navigate to submit page
    if g.browser.current_url != g.submit_url:
        g.browser.get(g.submit_url)
        sleep()

    # submit solution
    cf_submit(request.form['id'], request.form['lang'], request.form['text'])

    # go back to submit page
    g.browser.get(g.submit_url)
    return "solution sent"


@app.route(urljoin('/', CHECK_VARS_PART))
def home():
    return "handle: {}<br> password: {}<br><br> submit url: {}".format(
        g.handle,
        g.password,
        g.submit_url
    )
