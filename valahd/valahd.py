#!/usr/bin/env python3

from urllib.parse import urljoin

import requests

from urls import SERVER_ADDR, INIT_SERVER_PART, SUBMIT_PART, INIT_CONTEST_PART


def init_server(handle, password):
    data = {
        'handle': handle,
        'password': password,
    }
    requests.get(urljoin(SERVER_ADDR, INIT_SERVER_PART), data=data)


def init_contest(contest_id):
    data = {
        'num': contest_id,
    }
    requests.get(urljoin(SERVER_ADDR, INIT_CONTEST_PART), data=data)


def submit(problem_id, lang, text):
    data = {
        'id': problem_id,
        'lang': lang,
        'text': text,
    }
    requests.post(urljoin(SERVER_ADDR, SUBMIT_PART), data=data)


# should be used like this:
if __name__ == '__main__':
    init_server('xxx', 'yyy')
    init_contest('550')
    submit('A', '41', ''.join(open('550A.py', 'r').readlines()))
