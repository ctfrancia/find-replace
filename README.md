# find-replace

## Description

This is a simple CLI tool that allows you to search for a string in a directory and replace it with another string.
This was made because a bash script was not working for me and as a result I wanted to write it in Go.


#### Important Note
This tool works by opening the `go.mod` file and then taking the value of the `module` field.
Then goes through the rest of the files searching for that value and replacing it with the new value after the `-rs` flag.

## Installation
- download this repo
- run `go install .`

## Usage
- run `find-replace -rs "<VALUE>" -d /path/to/directory`

## Example

`$ find-replace -rs "foo"`
