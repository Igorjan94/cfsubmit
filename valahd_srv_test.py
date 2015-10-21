#!/usr/bin/env python3

import time
import unittest

from selenium import webdriver
from selenium.webdriver.common.keys import Keys

from urls import CF_LOGIN_URL


class TestCodeforces(unittest.TestCase):

    def setUp(self):
        self.browser = webdriver.Firefox()
        self.credentials = {
            "handle": 'xxx',
            "password": "yyy",
        }

    def tearDown(self):
        self.browser.close()

    def test_login_codeforces(self):
        self.browser.get(CF_LOGIN_URL)
        self.browser.find_element_by_id("handle").send_keys(self.credentials["handle"])
        self.browser.find_element_by_id("password").send_keys(
            self.credentials["password"])
        el = self.browser.find_element_by_id("remember")
        el.click()

        el.send_keys(Keys.ENTER)
        time.sleep(5)

        self.assertIn(
            self.credentials["handle"],
            self.browser.find_element_by_xpath("//div[@id='header']").text
        )

if __name__ == "__main__":
    unittest.main()
