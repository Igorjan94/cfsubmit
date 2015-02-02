cfsubmit
========

Send your solution to codeforces from command line

### Dependencies ###
golang  
Set up ```$GOPATH``` environment variable and make sure ```$GOPATH/bin``` is in your ```$PATH```  

If you don't want to install golang, use [gobuild.io](http://gobuild.io/) or similar service.

### Installation ###
```go get github.com/cnt0/cfsubmit```  
```go install github.com/cnt0/cfsubmit/...```

(Yes, this ```/...``` is mandatory)

### Setup ###
Place ```cfsubmit_settings.json``` into your working directory and open it.

Every time you log into codeforces you have to update ```X-User``` and ```CSRF-Token``` variables.  
```X-User```: look for "X-User" in your browser cookies  
```CSRF-Token```: open codeforces main page source and look for element ```<meta name="X-Csrf-Token"```  

If you prefer ```codeforces.com``` over ```codeforces.ru```, set ```CF-Domain``` to ```com```.  
  
Edit ```Ext-ID``` to set up your favorite codeforces compiler for each file extension. For example, you can choose between Java 6, 7 and 8 for ```.java``` files.

### Usage ###
```cfsubmit 123a.cpp```

cfmanage
========

Organize your submission files

### CLI Options ###
- ```-a```   arhive old submissions into one folder per contest; dominates -z flag
- ```-z```   arhive old submissions into one gzip file per contest
- ```-c```   create empty templates for contest; existing files will be rewritten
- ```-cnt``` how many templates will be created (at most 26, default: 5)
- ```-t```   which file will be used as base template (default: empty)

