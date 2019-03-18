# http stress test
HTTP Test and Stress Test Utility with ability to run hundreds of threads.  Uses a simplified TEXT format based on JSON to specify URI, METHOD, BODY, Supports checking response codes and RE matching against body.

# Generic HTTP Test Client

[httpTest](src/httpTest.go) provides a data driven, multi threaded test client able to support running at many threads for a while then waiting for all concurrent threads to finish before starting the next test.  This provides basic support for read after write tests.   It also provides easily parsed output that can be used to feed test results into downstream tools.  

The input to Generic Test client is a text file containing a series of JSON strings that describe each test.   It includes a few directives such as #WAIT to indicate a desire to wait for all prior requests to finish.   Comment strings are lines starting with #WAIT



> ## Example Output

```
httpTest.go
SUCESS: L125: id= 0183 	message= saving JSON record
L209: waiting queue= 0 reqPending= 1
SUCESS: L125: id= 0184G 	message= Read after write
L209: waiting queue= 0 reqPending= 1
SUCESS: L125: id= 0184M 	message= Read after write existing record
SUCESS: L125: id= 0185 	message= Check Mismatched 404 miss on known bad key
FAIL: L128: id=0186 	message=Check unepected response code  	err Msg=	L:97 Expected StatusCode= 519  got= 404
  	verb=GET uri=http://127.0.0.1:9601/mds/test/1817127X5
FAIL: L128: id=0184F 	message=Check re No Match functionality should fail 	err Msg= L120:FAIL ReNoMatch pattern found in record reNoMatch= .*JOHN MUIR.*  match= true
  	verb=GET uri=http://127.0.0.1:9601/mds/test/1817127X3
L209: waiting queue= 0 reqPending= 4
SUCESS: L125: id= 0187 	message= Delete a JSON record
L209: waiting queue= 0 reqPending= 1
SUCESS: L125: id= 0188 	message= Read after delete
SUCESS: L125: id= 0189 	message= Delete a previously deleted record
Finished Queing
 took 0.000468 min
Finished all test records
 took 0.000501 min
numReq= 9 elapSec= 0.0350985 numSuc= 7 numFail= 2 failRate= 0 reqPerSec= 256.42121458181975
```



## Assumptions

- Unless specified otherwise assumes all requests can be completed in any order which may in fact happen since they are ran in a multi threaded fashion.     



## Command Line API

```
httpTest  -in=data/test_airsolarwater_dem_index.txt -out=test1.log.txt  -maxthread=100
```

​    

> > - Runs the test with [data/sample-1.txt](data/sample-1.txt) as the input file.    
> > - Writing basic results to test1.log.txt 
> > - Runs with 100 client threads.
> > - **Parameters**
> > - **-in** = the name of the file containing test specifications.
> > - **-out** = the name of the file to write test results and timing to.
> > - **-MaxThread** = maximum number of concurrent requests submitted from client to servers.

## File Input Format

```json
# Test simple index page contains expected content.
{  "id" : "0182772-airSolarWater-dem-index-contains-solomon",   
   "verb" : "GET", 
   "uri" : "https://airsolarwater.com/dem/", 
   "expected" : 200, 
   "rematch" : ".*Solomon.*Tuvalu.*Tesla", 
   "message" :"Air solar Water index must contain island names"
}
#END

# Test Google page contains expected text
{  "id" : "0182-google-home-contains-search",   
   "verb" : "GET", 
   "uri" : "https://google", 
   "expected" : 200, 
   "rematch" : ".*Search\<\/div\>.*Google", 
   "message" :"Air solar Water index must contain island names"
}
#END

```



- A series of lines containing JSON text which represents the specifications for the test 

- Each JSON string is terminated by a #END starting a otherwise blank line. 

- **[Sample-1](data/sample-1.txt):**

  Test ID = ID to print upon failure
  HTTP Verb = HTTP Verb to send to the server
  URI =  URI to open for this test 
  Headers = Array of Headers to send to the server
  rematch = RE pattern to match the response body against. Not match is failure.
  renomatch = RE pattern that must not be in the response data.

   expected = HTTP Response code expected from the server. 

  ```
          other response codes are treated as failure.
  ```

   body = Body string to send as Post Body the server 

- Blank rows are ignored

- Rows prefixed by # are treated as comment except when #WAIT or #END

- Rows Prefixed with "#WAIT" Cause the system to pause and wait for all previously queued requests to complete before continuing.  This can allow blocking to allow data setup calls to complete before their read equivalents to complete.

- Test Requests  can be read from file and executed in parallel threads unless blocked by #WAIT directive.

- HTTP VERBS SUPPORTED  GET,PUT,DELETE,POST

- HTTP Headers URI Encoded sent in order specified but this can not be guaranteed since it is treated internally as a map which does not guarantee ordering. 

- POST BODY IS URI Encoded in file but will be decoded prior to POSTING TO Test client.



## Build / Setup

 Download the Metadata server repository

cd  RepositoryBaseDir 

​	eg:  cd \jsoft\mdsgo

​	This is the directory where you downloaded the repository.   

```
go build src/GenericHTTPTestClient.go
```

  or  

```
makeGO.bat
```



## Some other repositories:

- [file2consul](https://github.com/joeatbayes/file2consul)  Utility to load configuration parameters managed in GO into consul or HTTP server.  Supports inheritance,  parameter interpolation and other advanced techniques to minimize manual editing required to support multiple environments. 
- [GoPackaging](https://github.com/joeatbayes/GoPackaging) - Example of how to package a library for direct use from go command line.  Also shows an example program that uses that library.   
- [DevOps Automation with LXD and LXC containers](https://lxddevops.com/) - Scripts to create images that can be booted at will.  Demonstrates layered image creation.  Scripts to launch images,  map ports and to setup layer 5 routing for images ran across many hosts.  Consider an alternative to OpenStack or Kubernetes although it could run on OpenStack.    
- [Quantized Classifier](https://bitbucket.org/joexdobs/ml-classifier-gesture-recognition) -  Advanced machine learning classifier with full examples.  Very fast and competitive for both accuracy and and recall with tensor flow for many problems.  In some instances runs over 50 times faster than tensor flow with comparable precision. 
- [Metadata server MDS](https://bitbucket.org/joexdobs/meta-data-server) A server for storing arbitrary data by ID with very fast retrieval.   Ideal for use with very large data sets when very large scaling is required.   Compare to redis, memcache, memcacheb, riak, but designed to support very high scale on reasonable memory and can handle data sets of many T with good performance.   Pure HTTP based.  Written in GO.   Optimized along line of consuming a queue to keep multiple readers updated in near real-time.   Any single server can guarantee updates but servers as set are still eventual consistency. Measured at supporting over 16K requests at 4K per second  across a set of millions of records sustaining this load for months without degradation.
- [Computer Aided Call Response Engine](https://bitbucket.org/joexdobs/computer-aided-call-response-engine) Supports non technical user definition of call scripts that can be arbitrarily complex scripts. Based on JavaScript to integrate nicely with most RIA applications.  Features script tracking,   Local data storage,  spooling update events to server,   save and restore context for different users to allow call switching. etc.   Simple text based syntax to define the call tree allows rapid deployment and maintenance by non-programmers.
- [Healthcare provider Search](https://bitbucket.org/joexdobs/healthcare-provider-search)  Search UI that provides physician locator functionality.   Scripts to parse and load with over 1.2 million records from CMS.    Allows geo-proximity filtering,  zip to city resolution, etc.   Based on elastic search with node.js middle tier server.  Demonstrates very fast RIA JavaScript which allows hundreds of records faster that many sites render a one.   Take a close look if you want to understand high performance JavaScript.
- [CSV Tables in Browser](https://github.com/joeatbayes/CSVTablesInBrowser) Render CSV files in a browser automatically with little to know code.  Renders them fast and nice looking with automatically repeated headers.   Supports sorting,   Custom formatting by column and can even allow script callback to allow custom data generation for some fields.    Can dramatically reduce work to render large sets of columnar data. 

## License:

Copyright 2018 Joseph Ellsworth  [MIT License](https://opensource.org/licenses/MIT) 