import re

s = input()
a1 = re.match(r"BA.*AB", s) is not None
a2 = re.match(r"AB.*BA", s) is not None
print("YES" if a1 or a2 else "NO")