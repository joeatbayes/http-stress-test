# http stress tester
HTTP Test case runner and Stress Test Utility.  [httpTest](httpTest/httpTest.go) provides a data driven, multi threaded test client able to support running at many threads for a while then waiting for all concurrent threads to finish before starting the next test.  This provides basic support for read after write tests.   It also provides easily parsed output that can be used to feed test results into downstream tools.  

* Easily edited text format based on JSON to specify URI, METHOD, BODY and response matching.    [sample test script](data/sample.tst)

  *  Customize tests with Interpolation from named values on command line.
  *  Interpolation for custom URI and HTTP headers
  *  Custom interpolation for URI & OIDC to allow login & session management.
  *  Save results from prior tests that can be interpolated into subsequent tests.

*  Ability to run from 1 to hundreds of concurrent requests across a wide variety complex use cases.

   *  Maximum requests per second rate limiting to control load on server.
   *  Parallel execution with #WAIT semantics to force all operations to finish before next stage of testing.  Supports testing when some tests must complete before others can be started. 

*  Success Checking

   *  Regex matching of HTTP response and body to validate good responses.
   *  Regex must not match to check proper filtering. 
   *  Response code checking 

*  Output shows test ID and status for every test ran making it easy to integrate into a CICD pipeline. *Easily parsed output to allow easy integration with CICD pipelines.*

   

## Local Build / Setup

Once you have the  [golang compiler](https://golang.org/dl/) installed.    

```
go get -u -t "github.com/joeatbayes/http-stress-test/httpTest"
```

It will create a executable in your GOPATH directory in bin/httpTest.  For windows it will be httpTest.exe.  If GOPATH is set the /tmp then the executable will be written to /tmp/bin/httpTest under Linux or /tmp/bin/httpTest.exe for windows. 

> Once you build the httpTest executable you can copy it to other computers without the compiler.   It only requires the single executable file.
>

HINT: set GOTPATH= your current working directory.  Or set it to your desired target directory and the go get command will create the executable in bin inside of that directory which is good because you may not have write privileges to the default GOPATH.

##### To Download all pre-built test cases, scripts and sourcecode

```
git clone https://github.com/joeatbayes/http-stress-test httpTest
```

You could also just save the [sample script](https://raw.githubusercontent.com/joeatbayes/http-stress-test/master/data/sample-1.tst) using your browser or curl and edit it to use the executable built above. 

## Assumptions

- Unless specified otherwise assumes all requests can be completed in any order which may in fact happen since they are ran in a multi threaded fashion.     The #WAIT directive can be used to force prior commands to be completed before subsequent test  cases are executed.

## Command Line API

```
httpTest  -in=data/sample.tst -out=test1.log.txt  -MaxThread=100 -Environment=TST
```

​    

> > - Runs the test with [data/sample-1.txt](data/sample-1.tst) as the input file.    
> > - Writing basic results to test1.log.txt 
> > - Runs with 100 client threads.
> > - **Parameters**
> > - **-in** = the name of the file containing test specifications.
> > - **-out** = the name of the file to write test results and timing to.
> > - **-MaxThread** = maximum number of concurrent requests submitted from client to servers.  Defaults to 20 threads if not specified.
> > - **-mrps** =  Maximum Requests Per Second.  Causes the system to delay processing of test cases by enough to limit load on the receiving server.    A mrps value of 0.5 submits only 1/2 request per second or 1 request every 2 seconds.   Defaults to no limit if not specified.
> > - **-Environment** = Arbitrarily named command parameter.   These can be used and interpolated into the URI,  header keys, header values and body string.   Essentially any named value can be added in the same fashion eg:  -mykey=001 where mykey can be any set of alphanumeric characters and value is any valid as a command parameter in the os shell.
> >
#### Reading a Directory full of tests
> >
> >```
> >httpTest  -in=data/dir-test -Env=JOE1
> >
> ># Run all tests in the data/tests directory in alpha order.
> >```
> >
> >By specifying a directory name instead of a file in the -in parameter the system will read all files with an extension matching that specified by -ext.   In this instant -in=data/dir-test and -ext=tst will cause the system to find all files relative to the current working directory in data/dir-test that have an extension of .tst.    The list of files should be similar to that returned by ls -l data/dir-test/*tst.   The files are  processed in sorted order so it is easy to control test execution order by naming the file with a prefix such as 000-testx1.tst 001-testzbc.tst  Since 000 sorts before 001 it will be executed first. similar to the list returns 
> >
> >#### 
> >

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

# Test CGI Call to DEM Sub ssytem
{  "id" : "0182-DEM-CGI-CALL-GEOPOINT",   
   "verb" : "GET", 
   "uri" : "http://airsolarwater.com/dem/drains.svc?buri=gdata/dnt-rodrigues-island-50&offset=1591792&geo=-19.71542,63.34569", 
   "expected" : 200, 
   "rematch" : "-19.71542,63.34569,1.00.*-19.71736,63.34514", 
   "message" :"Checking CGI Geo Point"
}
#END

# Test Google page contains expected text
{  "id" : "0182-google-home-contains-search",   
   "verb" : "GET", 
   "uri" : "https://google.com", 
   "expected" : 200, 
   "rematch" : ".*Search.*div.*Google", 
   "message" :"Google home must contain 'search' followed by 'div' followed by 'Google'"
}
#END

```

> SEE Advanced Examples including Custom Headers below
>

- A series of lines containing JSON text which represents the specifications for the test 

- Each JSON string is terminated by a #END starting a otherwise blank line. 

- **[Sample](data/sample.tst):**

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

> ## Example Output

```
GenericHTTPTestClient.go
Finished Queing
 took 0.000033 min
Finished all test records
 took 0.000067 min
FAIL: L128: elap=382.978ms       id=0182-DEM-CGI-CALL-GEOPOINT-EXPECT-FAIL      message=Checking CGI Geo Point  err Msg= L107:failed rematch= -19.71542,63.34569,1.00.*-19.71736,63.XXX34514
        verb=GET uri=http://airsolarwater.com/dem/drains.svc?buri=gdata/dnt-rodrigues-island-50&offset=1591792&geo=-19.71542,63.34569

SUCESS: L125: elap=383.976ms     id=0182-DEM-CGI-CALL-GEOPOINT message=Checking CGI Geo Point
SUCESS: L125: elap=482.708ms     id=0182-google-home-contains-search    message=Google home must contain 'search' followed by 'div' followed by 'Google'
SUCESS: L125: elap=520.611ms     id=0182772-airSolarWater-dem-index-contains-solomon    message=Air solar Water index must contain island names
numReq= 4 elapSec= 0.5245668 numSuc= 3 numFail= 1 failRate= 0 reqPerSec= 7.625339613563039
```



## Important Files

- [data/sample.tst](data/sample.tst) - Sample input data to drive some simple tests

- [actions.md](actions.md) - list of feature enhancements under consideration.  Roughly listed in order.

- [httpTest.go](src/httpTest.go) - GO source code for main driver supporting this test.

- [makego.bat](makego.bat) - windows batch file to build the httpTest executable

- [makego.sh](makego.sh) - linux shell script to build the httpTest executable

  [goutil github repository](https://github.com/joeatbayes/goutil)  [httpTest.go](httpTest.go/httpTest.go) requires code from goutil that will be automatically downloaded when building this too.

## Advanced Usage 

### Limiting Server load or Limiting Requests Per Second

```
httpTest  -in=data/dir-test -ext=tst  -mrps=0.5

# Adding the -mrps value will cause the system to delay processing
# to each request to slow the system down by enough to limit demand 
# on the server so it only places a limited load on the server.   
# A mrps value of 0.5 submits only 1/2 request per second or 1 
# request every 2 seconds.
```

### Reading Multiple Files full of tests
>
> Multiple files or directories can be processed by separating the file names specified in the -in parameter with ";".    Items will be processed in the order specified but due to the multi-threaded nature of running tests it is possible that tests from multiple directories or files could be completed out of order especially when a slow service is mixed into test cases with faster tests.  The #WAIT directive can be used to force prior tests to finish.
>
> ##### Reading multiple files 
>
> ```
> httpTest  -in=data/sample.tst;data/dir-test/002-test-google.tst 
> ```
>
> #####  Reading multiple directories
>
> ```
> httpTest  -in=data/dir-test;data/airsolarwater/dem;  -Env=JOE1
> ```
>
> ##### Reading combined set of multiple files and directories.

> ```
> httpTest  -in=data/sample.txt;data/dir-test
> ```
>
> 

### Example with Interpolation & Custom  Headers

Get A session token from one REST call and pass it as part of custom header in the next test case.

```
httpTest  -in=data/login-with-token/002-sample-login-token-passed-as-header.tst -userid=testuser -passwd=tiger1928A2 -ENV=TST  

# Demonstrate passing the userid and password in as command line parameters
# and saving the result as a named value to be used as token in subsequent
# calls.
```



```json
# Login and get a session token
{  "id" : "018-login",   
   "verb" : "GET", 
   "uri" : "https://abcdex1.{ENV}.org/login?user={userid}&pass={passwd}",
   "expected" : 200, 
   "keepBodyAs" : "sessionToken",
   "keepBodyDefault" : "D1A2C8F10A0AB0C0A1DEFAULT"
}
#END
#WAIT


# Test simple index page contains expected content.
{  "id" : "0182772-airSolarWater-dem-index-contains-solomon",   
   "verb" : "GET", 
   "uri" : "https://airsolarwater.com/dem/", 
   "headers": {
      "sessionToken": "{sessionToken}",
      "Meta-Roles": "PUBLIC",
      "ENV" : "{ENVIRONMENT}"
   },
   "expected" : 200, 
   "rematch" : ".*Solomon.*Tuvalu.*Tesla", 
   "message" :"Air solar Water index must contain island names"
}
#END

```

```
#WAIT - Forces the Logic request to finish before 
  subsequent command is processed.   This is important because
  we need the value saved using the keepBodyAs to interpolate
  into the sessionToken header in subsequent tests.
  
keepBodyAs - Causes the results from the 01-login test case to be saved 
  so it can be interpolated into future scripts.
  
keepBodyDefault - saves this string value when the call fails.

headers {} - contains a list of headers to send with the command 
  useful when custom headers are needed. 
  
{sessionToken} - in the headers array is the value sent as a
  header for the key   "sessionToken}.  In this instance the
  sessionToken was set as a named value in the prior test case
  so we are actually passing that value interpolated from the
  saved values.   The interpolation is triggered by surrounding
  the name by {}.

{ENvIRONMENT} - is passed as the data portion of the http header
  "ENv" The actual value set is looked up in saved data values 
  because the name is surrounded by {}.     This one could have
  been set using by setting the value for ENVIORNMENT on the 
  command line as shown
  httpTest -in=data/sample-1.txt -out=test1.log.txt  -MaxThread=100 -Environment=TST   
  
  Any key value could be set on the command line and used for
  interpolation.  For example if -myfavoritecolor=blue is set 
  on the commmand line then it could be interpolated into the 
  body or header as {myfavoritecolor}.  
  
  Interpolation can alsue be used in URI and header values.
  For example if -ENV=test was used on command line and 
  the uri was http://api.{ENV}.abc.org the uri interpolated
  would be transformed into http://api.test.abc.org.  
  
```



### Example With custom verb & post body

```json
{
	"id": "0183",
	"verb": "PUT",
	"uri": "http://127.0.0.1:9601/mds/test/1817127X3",
	"headers": {
		"Content-Type": "application/JSON",
		"Meta-Roles": "PUBLIC"
	},
	"expected": 200,
	"rematch": ".*sucess.*",
	"message": "saving JSON record",
	"body": "{\"patientAgeRange\": \"NA\", \"orgName\": \"JOHN MUIR PHYSICIAN NETWORK\", \"combName\": \"PORTEOUS, BRENT\", \"product\": \"NA\", \"primaryLocation\": \"NA\", \"exclude\": false, \"uniqueLocKey\": \"JOHN MUIR PHYSICIAN NETWORK..1450 TREAT BLVD..945972168\", \"loc\": {\"lat\": 37.91, \"lon\": -122.07}, \"addr\": {\"city\": \"WALNUT CREEK\", \"zip\": \"945972168\", \"county\": \"NA\", \"state\": \"CA\", \"street\": \"1450 TREAT BLVD\", \"zipPlus4\": \"NA\"}, \"medicadeId\": \"6608789813\", \"languages\": \"NA\", \"email\": \"NA\", \"fax\": \"NA\", \"npi\": \"1083040463\", \"specialty\": [\"FAMILY PRACTICE\"], \"drName\": {\"middle\": \"\", \"last\": \"PORTEOUS\", \"suffix\": \"\", \"first\": \"BRENT\"}, \"phone\": \"9252969000\", \"publicTransitAccess\": \"NA\", \"handicapAccess\": \"NA\", \"credentials\": \"\", \"acceptMedicaid\": null, \"OfficeHours\": \"NA\", \"medSchool\": \"OTHER\", \"locations\": [{\"ccn\": \"050276\", \"lbn\": \"CONTRA COSTA REGIONAL MEDICAL CENTER\"}], \"gender\": \"M\", \"gradYear\": \"2012\"}"
}
#END
```

- NOTE:  To produce the post Body is must be escaped as a safe JSON string.  This can be done easily using an [online json escaper](https://www.freeformatter.com/json-escape.html)
- 

## Local  Development Build

####  Download the Metadata server repository

```
git clone https://github.com/joeatbayes/http-stress-test httpTest
```

```
# Make your currenct directory the directory where you
# downloaded the repository
# cd  RepositoryBaseDir 
cd \jsoft\httpTest
# Use forward slash on 
# This is the directory where you downloaded the repository.   
```

​	

```
go get -u "github.com/joeatbayes/goutil/jutil"
go build src/httpTest.go
```

  or  

```
# Windows
makeGO.bat

# Linux
makeGO.sh

```

### Build to a specific location Direct from Repo

```
go get -u -t "github.com/joeatbayes/http-stress-test/httpTest"
go build -o /tmp/httpTest.exe "github.com/joeatbayes/http-stress-test/httpTest"
# Requires sucessful execution of the go get command above.
```

This will build a new executable and place it at the location specified in the -o. For windows the .exe extension is needed for Linux leave it off.  This will be a duplicate of the executable built during the go get command so it is probably better just move the one built by go get to a location the search path.

## **How httpTest Compares**

- httpTest compares to [newman portion of postman](https://github.com/postmanlabs/newman) but supplies superior multi-threaded performance or stress testing.   Newman scripts are more flexible but httpTest scripts are easier to write and provide better scalability testing.   
- Provides free and higher performance stress test than  [load runner](https://www.microfocus.com/en-us/products/loadrunner-load-testing/overview)  httpTest is better suited to rest service testing while load runner is better suited to full GUI/browser emulation.  
- httpTest provides similar functionality to [JMeter](https://jmeter.apache.org/) when JMeter is configured to read URI from CSV but provides superior thread merge control for dependent transactions and far superior control when switching between very different requests that simulate a entire ecosystem.   httpTest is easier to generate input scripts from large data sets and has a much shorter learning curve.  JMeter is more mature, requires a larger stack and consumes more resources.  JMeter provides a GUI but is harder to configure when mutli-threaded testing of complex uses cases for REST services. 

* 

#### Supporting OIDC From the Tester

* TODO: Add and Example Here

## Some of my other repositories:

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