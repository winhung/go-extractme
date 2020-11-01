# go-extractme
![Go](https://github.com/winhung/go-extractme/workflows/Go/badge.svg)

Search and extract keywords to another file

## Quick start
Usage example: ./go-extractme -ct=tf2json -of="dev qe sit uat prd" -if=convertme.tf -verify=true -od=secrets -rf=variables.tf
### Compiled binaries
[LINK HERE](http://google.com)

## Problem statement
One day, i was given a task to extract sensitive data like username, passwords, API keys etc. from Terraform files into JSON files. Each of JSON file represented the environment they were going to be used in ( eg. Development(dev.json) , QA(qa.json), Production(prd.json) ). These JSON files then had to be encrypted with [git-crypt](https://github.com/AGWA/git-crypt) and the Terraform files that used to read those data, now had to be amended to read from the respective JSON files instead.

There were quite a number of variables and the JSON files could increase one day which meant the following problems
* too much manual work
* mistakes were bound to happen with manual work

Hence, the idea to come up with such a tool came up !

If this is something you are facing then this tool might be able to help you.

## Available CLI parameters
| CLI parameter | Mandatory or Optional ?| Default value (if any) | Explanation  | Usage example |
| :------------ | :--------------------- | :--------------------- | :------------| :------------ |
| ct | Mandatory | N/A | Conversion type. Convert the input file type to an output file. Refer to [Supported conversions](#anchor-supportconv) for a list of supported conversion types and their enum values | -ct=tf2json |
| of | Mandatory | N/A | Output filename(s). Name(s) of the output file(s) | Single file: -of=dev Multiple files: -of"dev qa prd" |
| if | Mandatory | N/A | Input filename. Name of the input file to read and extract the value(s) from | -if=convertme.tf |
| od | Optional | output | Output directory. The name of the directory where the files will be output to. The folder should have already been created | -od=secrets |
| rf | Optional | N/A | File to amend such that it is updated with new location of extracted values | -rf=parameters.tf |
| rfko | Optional | true |  If true, will keep the amended file as it is by creating a copy of it with the amended values. If false, will overwrite the file specified in 'rf' with updated values. Used with 'rf' | -rfko=true |
| verify | Optional | false | If true, a check will be done on the input files (from -if) and the output files (from -rf) to ensure that they contain the correct values.


## Supported conversions <sup id="anchor-supportconv" />
| Conversion types | CLI parameter for -ct flag |
| :--------------- | :-------------------------- |
| Terraform to JSON | tf2json |



## Notes before using this tool
Please ensure that all comments are removed.
The file that you want to read, contains only the variables that needs to be converted into JSON file.
