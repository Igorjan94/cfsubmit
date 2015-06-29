import requests
import time


def init_server(handle, password):
    data = {
        'handle': handle,
        'password': password,
    }
    requests.get("http://localhost:5555/init_server", data=data)


def init_contest(contest_id):
    data = {
        'num': contest_id,
    }
    requests.get("http://localhost:5555/init_contest", data=data)

def submit(problem_id, lang, text):
    data = {
        'id': problem_id,
        'lang': lang,
        'text': text,
    }
    requests.post("http://localhost:5555/submit", data=data)


# should be used like this:
if __name__ == '__main__':
    init_server('xxx', 'yyy')
    init_contest('550')
    submit('A', '41', ''.join(open('550A.py', 'r').readlines()))
