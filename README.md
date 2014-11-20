cfsubmit
========

Send your solution to codeforces from command line

### Dependencies ###
python 3.x (didn't tested it under python 2.x)  
```pip install requests```

### Setup ###
Every time you log into codeforces you have to update ```x_user``` and ```csrf_token``` variables.  
```x_user```: look for "X-User" in your browser cookies  
```csrf_token```: open codeforces main page source and look for element ```<meta name="X-Csrf-Token"```  

If you prefer ```codeforces.ru``` over ```codeforces.com```, set ```cf_domain = "ru"```.  
  
Edit ```ext_id``` variable to set up your favorite codeforces compiler for each file extension. For example, you can choose between Java 6, 7 and 8 for ```.java``` files.

### Usage ###
```python cfsubmit.py 123a.cpp```