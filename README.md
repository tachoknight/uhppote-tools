# UHPPOTE tools
This is a collection of Go-based programs that perform various tasks with the [UHPPOTE Professional Wiegand 26 Bit TCP IP Network Access Control Board Panel Controller For 4 Door 4 Reader](https://www.amazon.com/dp/B00UX02JWE/ref=asc_df_B00UX02JWE5402121/?tag=hyprod-20&creative=395033&creativeASIN=B00UX02JWE&linkCode=df0&hvadid=167135357001&hvpos=1o1&hvnetw=g&hvrand=9692092155284663930&hvpone=&hvptwo=&hvqmt=&hvdev=c&hvdvcmdl=&hvlocint=&hvlocphy=9060254&hvtargid=pla-313892092019).

## Setup
For all the programs you need to fill in following fields at the top of the file (or load from a config file or db or whatever):
* boardSerialNum
* boardIP
* boardPort

## `getsetboardtime.go`
The board has a clock but does not handle things like DST. This program will set the board's clock to the current date and time; this presumes the board's location is the same where the program runs.

### Use in a package
`GetBoardTime` and `SetBoardTime` are the two public functions.

## `user-access-management.go`
This is Go program queries, adds, and deletes user tags, as well as getting the current list of access events from the board. 

### Use in a package
`AccessRecord`, `AddUser`, `GetUser`, and `DelUser` are exportable for easy integration into your code. All the other functions are used internally and aren't really useful.
