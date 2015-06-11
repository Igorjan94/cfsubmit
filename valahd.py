#!/usr/bin/env python3

import os
import requests
import sys
import time


valahd_init_addr = "http://localhost:4244/init"
valahd_submit_addr = "http://localhost:4244/submit"


def init(site, credentials):
    parts = {
        "site": site,
    }
    for (key, value) in credentials:
        parts[key] = value
    requests.post(valahd_init_addr, files=parts)
    print("Valahd server initialized. Current time is " +
          time.strftime("%H:%M:%S"))


def submit(filename):
    parts = {
        "filename": filename,
        "source":   open(sys.argv[1], "rb"),
    }
    requests.post(valahd_submit_addr, files=parts)
    print("Solution sent. Current time is " + time.strftime("%H:%M:%S"))

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("Solution filename not specified")
        sys.exit()
    if not os.path.exists(sys.argv[1]):
        print("Solution file does not exist or not enough rights to read it")
        sys.exit()
    submit(sys.argv[1])
