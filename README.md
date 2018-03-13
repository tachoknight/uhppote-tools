# UHPPOTE tools
This is a collection of Go-based programs that perform various tasks with the [UHPPOTE Professional Wiegand 26 Bit TCP IP Network Access Control Board Panel Controller For 4 Door 4 Reader](https://www.amazon.com/dp/B00UX02JWE/ref=asc_df_B00UX02JWE5402121/?tag=hyprod-20&creative=395033&creativeASIN=B00UX02JWE&linkCode=df0&hvadid=167135357001&hvpos=1o1&hvnetw=g&hvrand=9692092155284663930&hvpone=&hvptwo=&hvqmt=&hvdev=c&hvdvcmdl=&hvlocint=&hvlocphy=9060254&hvtargid=pla-313892092019). 

Note that these programs contain their own `main()` function as they were developed to run independently. I mention it because some editors see a collection of Go files in a directory as a single program and might complain about duplicate defintions.

## Setup
For all the programs you need to fill in following fields at the top of the file (or load from a config file or db or whatever):

* boardSerialNum
* boardIP
* boardPort

## `get-set-board-time.go`
The board has a clock but does not handle things like DST. This program will set the board's clock to the current date and time; this presumes the board's location is the same where the program runs.
### Use in a package
`GetBoardTime` and `SetBoardTime` are the two public functions that will return and set the time, respectively. 

## `user-access-management.go`
This Go program queries, adds, and deletes user tags. It has an additional function for converting the scanned tag serial number into the format that the board seems to really want it.
### Use in a package
`AddUser`, `GetUser`, and `DelUser` are exportable for easy integration into your code. All the other functions are used internally and aren't really useful.

## `get-access-list.go`
Gets the current list of access events from the board and formats it for output to stdout. You can modify it to write to a database or whatever.
### Use in a package
`GetAccessListCount()` and `GetAccessList()` are the exportable functions. Note that `GetAccessList()` takes a parameter, `count`, which returns that number of events. 
