# PlagDetector (go - plag)
[![Build Status](https://travis-ci.org/sohlich/go-plag.svg?branch=master)](https://travis-ci.org/sohlich/go-plag)

Source code plagiarism detection tool based on comparison of source code tokens.




## Usage:


Get supported languages:

GET /plugin/langs


Create assignment:

PUT /assignment
```	
{"name": "Test1","lang": "java"	} 
```


Create submission:

PUT /submission

 multipart:
 
```
submission-meta: 
{ "owner": "student", "assignmentId": "55c7b6ebe13823356f000001",  "id": "62" } 

submission-data:
content of file
```