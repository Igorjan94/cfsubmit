#!/usr/bin/env python3

from selenium import webdriver
from selenium.webdriver.common.keys import Keys
import time
import unittest


class TestCodeforces(unittest.TestCase):

    def setUp(self):
        self.browser = webdriver.Firefox()
        self.creds = {
            "handle": "CountZero",
            "password": "yyy",
        }

    def tearDown(self):
        self.browser.close()

    def test_login_codeforces(self):
        self.browser.get("http://codeforces.com/enter")
        self.browser.find_element_by_id("handle").send_keys(self.creds["handle"])
        self.browser.find_element_by_id("password").send_keys(
            self.creds["password"])
        el = self.browser.find_element_by_id("remember")
        el.click()

        el.send_keys(Keys.ENTER)
        time.sleep(5)

        self.assertIn(
            self.creds["handle"],
            self.browser.find_element_by_xpath("//div[@id='header']").text
        )

if __name__ == "__main__":
    unittest.main()
