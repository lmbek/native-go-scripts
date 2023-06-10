# This repository contains:
Some go scripts that can: 
* make screenshots and save them as images
* handle keyboard and send keys 
* move the mouse around 
* give warning/alerts 
* give yes/no dialog 
* prompt folder selection

## Note:
This only works on windows as we use windows api

## Article about the use of win32 and some of its risks
As this package uses unsafe and syscalls and depend on the windows operating system, this article should help understand some of it
https://justen.codes/breaking-all-the-rules-using-go-to-call-windows-api-2cbfd8c79724
Also the windows api documentation comes in handy
