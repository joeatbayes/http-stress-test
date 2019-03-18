## ACTIONS: 

* Add Ability to Read list of files on command line rather than single file.
* Add ability to save results of a command into a named variable and then allow that variable to be interpolated into the JSON string before parsing.
* Add ability to specify a directory on command line rather than single file.  When specified process every file found in directory.
* Allow specification of a default file extension.  Ignore files in directory that do not have the specified extension.
* Add Timing Library for performance logging example.

## Under Consideration

- Add "saveas" parameter to test spec so results of HTTP call are saved to local file as if fetched by Curl.  If present then treated as a relative file name relative to output file.  
- Add ability to run for a while at a given concurrency level and then increase concurrency to find the sweet spot for the server for the current set of data.
- Consider output format that uses JSON to make parsing easy.
- Modify logging output to use atomic output to avoid mixing lines in threads and to reduce sync calls.  Current version could easily mix logs on same line when heavily multi-threaded.
- Add option to run several test clients simultaneously using os spawn.  This may help ensure client is not blocked in context switching overhead.   May need a way to collate output from multiple spawned processes to get cumulative throughput.



## DONE:

- DONE: JOE:2017-10-20: Add ReNoMatch to ensure a given string is not in the result string from the service. 