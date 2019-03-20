## ACTIONS: 

* Add Ability to Read list of files on command line rather than single file.
* Add ability to specify a directory on command line rather than single file.  When specified process every file found in directory.
* Allow specification of a default file extension.  Ignore files in directory that do not have the specified extension.
* Add option to save header value 
* Add option to limit transactions to X per second.  Slow down the next request processed until RPS drops below the requested level. 
* Add option to keep URI timing by base URI upto the ? or # as a set to show at end of report. Would want the fastest, slowest, average and number that exceed SLA.
* Add option to read the body contents from a file relative to the location of the current script file.
* Add the option for a #include which will include a  named in the input script.   
* Modify input to use YAML or similar spec rather than JSON which requires more editing than desired. 


## Under Consideration

- Add "saveas" parameter to test spec so results of HTTP call are saved to local file as if fetched by Curl.  If present then treated as a relative file name relative to output file.  
- Add ability to run for a while at a given concurrency level and then increase concurrency to find the sweet spot for the server for the current set of data.
- Consider output format that uses JSON to make parsing easy.
- Modify logging output to use atomic output to avoid mixing lines in threads and to reduce sync calls.  Current version could easily mix logs on same line when heavily multi-threaded.
- Add option to run several test clients simultaneously using os spawn.  This may help ensure client is not blocked in context switching overhead.   May need a way to collate output from multiple spawned processes to get cumulative throughput.



## DONE:

- DONE:JOE:2019-03-19: Add ability to save results of a command into a named variable and then allow that variable to be interpolated into the JSON string before parsing. (Required for OIDC)
- DONE:JOE:2019-03-19: Add ability to parse environment variables specified on command line and interpolate into the script.
- DONE: JOE:2017-10-20: Add ReNoMatch to ensure a given string is not in the result string from the service. 