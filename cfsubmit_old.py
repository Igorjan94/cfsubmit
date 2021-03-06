#!/usr/bin/python

import os
import sys
import itertools


# Edit these four variables according to your needs:
ext_id = {
    "cpp":   "42",
    "cs":    "29",
    "go":    "32",
    "java":  "36",
    "py":    "31",
}

# programTypeId - Language
# 1:   GNU G++ 4.9.2
# 10:  GNU GCC 4.9.2
# 12:  Haskell GHC 7.6
# 13:  Perl 5.12
# 14:  ActiveTcl 8.5
# 15:  Io-2008-01-07 (Win32)
# 17:  Pike 7.8
# 18:  Befunge
# 19:  OCaml 4
# 2:   Microsoft Visual C++ 2010
# 20:  Scala 2.11
# 22:  OpenCobol 1.0
# 23:  Java 7
# 25:  Factor
# 26:  Secret_171
# 27:  Roco
# 28:  D DMD32 Compiler v2
# 29:  MS C# .NET 4
# 3:   Delphi 7
# 31:  Python 3.4
# 32:  Go 1.2
# 33:  Ada GNAT 4
# 34:  JavaScript V8 3
# 36:  Java 8
# 38:  Mysterious Language
# 39:  FALSE
# 4:   Free Pascal 2
# 40:  PyPy 2.5.0
# 41:  PyPy 3.2.5
# 42:  GNU G++11 4.9.2
# 43:  GNU GCC C11 4.9.2
# 44:  Picat 0.9
# 45:  GNU C++ 11 ZIP
# 46:  Java 8 ZIP
# 6:   PHP 5.3
# 7:   Python 2.7
# 8:   Ruby 2
# 9:   C# Mono 2.10


def is_gym(contest_id):
    return int(contest_id) >= 100000

if len(sys.argv) < 2:
    print("Solution filename not specified")
    sys.exit()

if not os.path.exists(sys.argv[1]):
    print("Solution file does not exist or not enough rights to read it")
    sys.exit()

filename = os.path.basename(sys.argv[1])

contest_id = ''.join(itertools.takewhile(lambda c: c.isdigit(), filename))
problem_index = ''.join(itertools.takewhile(
    lambda c: c != '.', filename[len(contest_id):])).upper()
extension = filename[len(contest_id) + len(problem_index) + 1:].lower()

if (len(contest_id) == 0) or (len(problem_index) == 0):
    print("Incorrect filename format. Example: 123A.cpp")
    sys.exit()
if extension not in ext_id:
    print("Unknown extension. Please check 'ext_id' variable")
    sys.exit()

submit_addr = "http://codeforces.com/{}/{}/submit".format(
    "gym" if is_gym(contest_id) else "contest", contest_id)
