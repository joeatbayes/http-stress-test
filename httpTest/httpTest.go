package main

// httpTest.go
// See:  httpTest.md
// See: ../data/test_airsolarwater_dem_index.txt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
        "io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	s "strings"
	"time"
    "sync"
	"github.com/joeatbayes/goutil/jutil"
)

type MTReader struct {
	perf         *jutil.PerfMeasure
	reqPending   int
	logFile      *os.File
	pargs        *jutil.ParsedCommandArgs
	inPaths      []string
	processExt   string // file name extension to use when processing directories.
	linesChan    chan TestSpec
	isDone       chan bool
	maxRecPerSec float64
	reqMade      int
	start        float64
	writeLock    sync.Mutex
}

func makeMTReader(outFiName string) *MTReader {
	r := MTReader{}
	r.perf = jutil.MakePerfMeasure(25000)
	logFiName := outFiName
	var logFile, sferr = os.Create(logFiName)
	if sferr != nil {
		fmt.Println("Can not open log file ", logFiName, " sferr=", sferr)
	}
	r.logFile = logFile
	r.reqMade = 0
	r.start = jutil.Nowms()
	r.maxRecPerSec = -1.0
	return &r
}

func (r *MTReader) elapSec() float64 {
	return jutil.CalcElapSec(r.start)
}

func (r *MTReader) requestsPerSec() float64 {
	elap := r.elapSec()
	if elap == 0.0 {
		return 0.0
	} else {
		return float64(r.reqMade) / elap
	}
}

func (r *MTReader) done() {
	defer r.logFile.Close()
}

type TestSpec struct {
	Id              string
	Verb            string
	Uri             string
	Headers         map[string]string
	Expected        int
	Rematch         string
	ReNoMatch       string
	Message         string
	Body            string
	KeepBodyAs      string
	KeepBodyDefault string
}

func keepLines(s string, n int) string {
	result := strings.Join(strings.Split(s, "\n")[:n], "\n")
	return strings.Replace(result, "\r", "", -1)
}

func (u *MTReader) procLine(spec *TestSpec) {
        //time.Sleep(2 * time.Second)
	u.reqPending += 1
	u.perf.NumReq += 1
	if u.maxRecPerSec != -1 {
	  fmt.Println("maxRecPerSec=", u.maxRecPerSec, " requestsPerSec=", u.requestsPerSec())
	  for u.maxRecPerSec > 0 && u.requestsPerSec() > u.maxRecPerSec && u.reqMade > 0 {
		// slow down to comply with maxReqPerSec directive
		// whenever maxReqPerSec has been specified as a
		// positive number.
		//fmt.Println("maxRecPerSec=", u.maxRecPerSec, " requestsPerSec=", u.requestsPerSec())
		time.Sleep(500.0)
		// TODO: Compute exact time we need to sleep to reach
		// the RPS and then delay by that amount rather than
		// polling.
	  }
    }
	u.reqMade += 1
	startms := jutil.Nowms()
	//fmt.Println("L45: spec=", spec)
	//fmt.Println("L49: spec.Rematch=", spec.Rematch)
	uri := spec.Uri
	// ===== JOE THIS IS SUSPICOUS See:http://networkbit.ch/golang-http-client/
	// req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	hc := http.Client{}
	reqStat := true
	errMsg := ""
	//fmt.Println("L49: spec id=", spec.Id, " keepBodyAs=", spec.KeepBodyAs, "default=", spec.KeepBodyDefault)
    //fmt.Println("L117: uri=", uri)
	req, err := http.NewRequest(spec.Verb, u.pargs.Interpolate(uri), bytes.NewBuffer([]byte(u.pargs.Interpolate(spec.Body))))
	//fmt.Println("L50: req=", req, " err=", err)
	if err != nil {
		u.perf.NumFail += 1
		serr := fmt.Sprintln(u.logFile, "FAIL: L74: id=", spec.Id, " message=", spec.Message, " error opening uri=", uri, " err=", err)
		fmt.Println(serr)
		fmt.Fprintln(u.logFile, serr)
		u.reqPending--
		reqStat = false
		return
	}

	//  Ieterate over set headers
	if spec.Headers != nil {
		for key, val := range spec.Headers {
			//fmt.Println("L94: header key=", key, " val=", val)
			sval := u.pargs.Interpolate(val)
			skey := u.pargs.Interpolate(key)
			req.Header.Set(skey, sval)
			//fmt.Println("L105: sval=", sval)
			//fmt.Println("L96: header key=", key, " val=", val, " skey=", skey, " val=", sval)
		}
		if spec.KeepBodyAs > " " {
			// modify save name to make it easier to use without
			// having to match the case in interpolation
			spec.KeepBodyAs = s.TrimPrefix(s.ToLower(spec.KeepBodyAs), " ")
		}
	}
	req.Header.Set("Connection", "keep-alive") // use "close" when you want to disable keep alive
	//req.Header.Set("Connection", "close") // use "close" when you want to disable keep alive
	//req.Close = true // When this is true it prevents http from using keep-alive.
    req.Close = false // set to false when http keep alive is desired.
	resp, err := hc.Do(req)
	//fmt.Println("L146: reps=", resp, "err=", err)
	if err != nil {
		fmt.Fprintln(u.logFile, "FAIL: L85: id=", spec.Id, " message=", spec.Message, "err=", err)
                if (err == io.EOF) {
                   fmt.Println("L157: Error premature closed connection")
                   http.DefaultTransport.(*http.Transport).CloseIdleConnections()
                   time.Sleep(2 * time.Second)
		}
                fmt.Println("FAIL: L85: id=", spec.Id, " message=", spec.Message, "err=", err)
		reqStat = false
		if spec.KeepBodyAs > " " && spec.KeepBodyDefault > " " {
			//fmt.Println("L106: keep default as failure keepBodyAs=", spec.KeepBodyAs, " default=", spec.KeepBodyDefault)
			//fmt.Println("L108: u.pargs=", u.pargs)
			//fmt.Println("L107: u.pargs.NamedStr=", u.pargs.NamedStr)
			u.pargs.NamedStr[spec.KeepBodyAs] = spec.KeepBodyDefault
		}
     	u.perf.NumFail += 1
		u.reqPending--
		return
	}
	
	//fmt.Println("L166: resp=", resp, " status=", resp.StatusCode, "err=", err)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)
	if resp.StatusCode != spec.Expected {
		errMsg = fmt.Sprintln("\tL:97 Expected StatusCode=", spec.Expected, " got=", resp.StatusCode)
		reqStat = false
	}
	//resp.Body.Close()


	//fmt.Println("L173: bodyStr=", bodyStr)
	// Add Logic to check the RE Pattern
	// match for body result
	if spec.Rematch != "" && spec.Rematch > " " {
		match, _ := regexp.MatchString(spec.Rematch, bodyStr)
		//fmt.Println("L86 match=", match, "merr=", merr)
		if match != true {
			errMsg = fmt.Sprintln(" L107:failed rematch=", spec.Rematch)
			reqStat = false
		}
	}
	
	// If KeepAs is specified in the input spec then save it for
	// latter interpolation.  If the call fails then save the default
	// if specified.
    if spec.KeepBodyDefault > " " {
      u.writeLock.Lock()
	  if resp.StatusCode != spec.Expected  {
	 	u.pargs.NamedStr[spec.KeepBodyAs] = spec.KeepBodyDefault
	  } else {
		u.pargs.NamedStr[spec.KeepBodyAs] = bodyStr
      }
      u.writeLock.Unlock()
    }
    
	// Add Logic to check the RE Pattern
	// match for body result
	if spec.ReNoMatch != "" && spec.ReNoMatch > " " {
		match, _ := regexp.MatchString(spec.ReNoMatch, bodyStr)
		//fmt.Println("118 match=", match, "ReNoMatch=", spec.ReNoMatch)
		if match == true {
			errMsg = fmt.Sprintln(" L120:FAIL ReNoMatch pattern found in record reNoMatch=", spec.ReNoMatch, " match=", match)
			reqStat = false
		}
	}
    
	u.perf.NumSinceStatPrint += 1
	u.perf.CheckAndPrintStat(u.logFile)

	//fmt.Println("L82: body len=", len(string(body)))
	//	fmt.Println("L82: body =", string(body))
	//len(body)
	//defer jutil.TimeTrack(now, "finished id=" + id)
	endms := jutil.Nowms()
	elapms := endms - startms
	if reqStat == true {
		u.perf.NumSuc += 1
		tbuff := fmt.Sprintf("SUCESS: L125: elap=%4.3fms\trps=%4.1f\tid=%s\tmessage=%s", elapms, u.requestsPerSec(), spec.Id, spec.Message)
		fmt.Println(tbuff)
		fmt.Fprintln(u.logFile, tbuff)
	} else {
		u.perf.NumFail += 1
		tbuff := fmt.Sprintf("FAIL: L128: elap=%4.3fms\t id=%v \tmessage=%s \terr Msg=%s  \tverb=%s uri=%s\n", elapms, spec.Id, spec.Message, errMsg, spec.Verb, spec.Uri)
		fmt.Fprintln(u.logFile, tbuff)
		fmt.Printf(tbuff)
	}
	u.reqPending--
	time.Sleep(10) 
        //time.Sleep(2 * time.Second)
    
}

var (
	server *http.Server
	client *http.Client
)

func PrintHelp() {
	fmt.Println(
		`httpTest -in=InputFileName -out=OutputFileName -MaxThread=5 -ENV=TST
	   -in defaults to data/sample.txt
	    name of input paramter file or directory
        If named resource is directory will process all 
	    files in that directory.    Multiple inputs
	    can be specified seprated by ;.  Each input set
	    will be processed in order
		
	   -out defaults  httpTest.log.txt 
	      Output log file where results will be written to allow
		  for secondary analysis. 
	
	   -MaxThread defaults to 20
	
	   -ext = file extension to include for processing 
	      when processing a full directory.  Defaults to
		  txt.
		
	   -mrps = Maximum number of requests to make
	      per second.   This is useful when wanting to stress
		  a server to a specific level regaurdless of the 
		  number of threads.  Defaults to no limit when 
		  not set. 
		
	   -env =  variable used for interpolation.  An arbitrary
	      set of variables can be defined in this fashion 
		  and will be available for use with interpolation
		
       -company = variable used for interpolation
	-`)
}

// Process a single input file containing various test cases.
func (u *MTReader) processFile(inFiName string) {
	// Add the rows to the Queue
	// so we can process them in parralell
	// It is blocked at MaxQueueSize by the
	// channel size.
	start := time.Now().UTC()
	inFile, err := os.Open(inFiName)
	if err != nil {
		fmt.Println("error opening input file ", inFiName, " err=", err)
		os.Exit(3)
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	var b bytes.Buffer
	for scanner.Scan() {
		aline := scanner.Text()
		aline = strings.TrimSpace(aline)
		if len(aline) < 1 {
			continue
		} else if aline[0] == '#' {
			if strings.HasPrefix(aline, "#END") {
				recStr := strings.TrimSpace(b.String())
				//fmt.Println("L228: RecStr=", recStr)
				if len(recStr) > 0 {
					b.Reset()
					spec := TestSpec{}
					err := json.Unmarshal([]byte(recStr), &spec)
					//fmt.Println("L150: spec=", &spec, "err=", err)
					if err != nil {
						fmt.Println("L222: FAIL: to parse err=", err, "str=", recStr)
					} else {
						u.linesChan <- spec
					}
				}
			} else if strings.HasPrefix(aline, "#WAIT") {
				// Add a Pause
				time.Sleep(9500)
				fmt.Fprintln(u.logFile, "L208: waiting queue=", len(u.linesChan), "reqPending=", u.reqPending)
				for len(u.linesChan) > 0 || u.reqPending > 0 {
					u.logFile.Sync()
					fmt.Fprintln(u.logFile, "L209: waiting queue=", len(u.linesChan), "reqPending=", u.reqPending)
					u.logFile.Sync()
					time.Sleep(5500)
					continue
				}
			} else {
				continue
			}
		} else {
			// Add current line to buffer
			b.Write([]byte(aline))
		}
	}

	u.logFile.Sync()
	jutil.TimeTrack(u.logFile, start, "Finished Queing "+inFiName+"\n")

}

// Process a directory processing any files that match
// -ext settings.
func (u *MTReader) processDir(inFiName string) {
	fmt.Println("L295: processDir ", inFiName)
	globPath := inFiName + "/*" + u.processExt
	//fmt.Println("L297: globPath=", globPath)
	files, err := filepath.Glob(globPath)
	if err != nil {
		fmt.Println("L289: ERROR processsing dir", inFiName, " err=", err)
	} else {
		//fmt.Println("L300: files=", files)
		for _, fiPath := range files {
			//fmt.Println("L301:  fiPath=", fiPath)
			u.processFile(fiPath)
		}
	}
}

// Process a path and determine whether it needs to
// be processed as a file or a extension.
func (u *MTReader) processPath(inFiName string) {
	fi, err := os.Stat(inFiName)
	switch {
	case os.IsNotExist(err):
		fmt.Println("L294: Error file ", inFiName, " does not exists")
	case err != nil:
		fmt.Println("L295: Error processing file ", inFiName, " err=", err)
	case fi.IsDir():
		u.processDir(inFiName)
	default:
		u.processFile(inFiName)
	}
}

func main() {
	startms := jutil.Nowms()
	const procs = 2
	const DefMaxWorkerThread = 20 //5 //150 //5 //15 // 50 // 350
	const MaxQueueSize = 3000
	const DefInFiName = "data/sample.tst"
	const DefOutFiName = "httpTest.log.txt"
	parms := jutil.ParseCommandLine(os.Args)
	if parms.Exists("help") {
		PrintHelp()
		return
	}
	MaxWorkerThread := parms.Ival("maxthread", int(DefMaxWorkerThread))
	fmt.Println(parms.String())
        transport := http.DefaultTransport.(*http.Transport) 
        transport.MaxIdleConnsPerHost = MaxWorkerThread + 10
        transport.MaxIdleConns = MaxWorkerThread + 10
        transport.MaxConnsPerHost = MaxWorkerThread + 10
        transport.DisableKeepAlives = false
        
	inPathStr := parms.Sval("in", DefInFiName)
	outFiName := parms.Sval("out", DefOutFiName)
	u := makeMTReader(outFiName)
	u.inPaths = strings.Split(inPathStr, ";")
	u.processExt = parms.Sval("ext", "tst")
	u.maxRecPerSec = parms.F64val("mrps", float64(-1.0)) // maxrecpersec
	fmt.Fprintln(u.logFile, "GenericHTTPTestClient.go")
	fmt.Println("OutFileName=", outFiName)
	fmt.Println("MaxWorkerThread=", MaxWorkerThread)
	fmt.Println("u.maxRecPerSec=", u.maxRecPerSec)

	start := time.Now().UTC()
	u.pargs = parms
	u.linesChan = make(chan TestSpec, MaxQueueSize)
	u.isDone = make(chan bool)

	// Spin up 100 worker threads to post
	// content to the server.

	for i := 0; i < MaxWorkerThread; i++ {
		go func() {
			for {
				spec, more := <-u.linesChan
				if more {
					u.procLine(&spec)
					//fmt.Println("L128 spec=", spec)
				} else {
					u.isDone <- true
					fmt.Println("L413: isDone()")
					return
				}
			}
		}()
	}

	// Process all the files listed in the -in parameter
	for _, path := range u.inPaths {
		path = strings.TrimSpace(path)
		u.processPath(path)
	}
	fmt.Println("L425: All files read")
	close(u.linesChan)
	jutil.Elap("L374: finished processing all input", startms, jutil.Nowms())
	time.Sleep(1)
	<-u.isDone // wait until queue has been marked as finished.
	jutil.TimeTrack(u.logFile, start, "Finished Read test records\n")
    finishWaitStart := jutil.Nowms()
	for u.reqPending > 0 {
          time.Sleep(2 * time.Second)
          fmt.Println("L434: Finished wait u.reqPending=", u.reqPending)
          jutil.Elap("L435: finish wait for", finishWaitStart, jutil.Nowms())
          if (jutil.Nowms() - finishWaitStart > 60000) {
            fmt.Println("L436: abort waited to long")
            break
          }
	}
	if u.reqPending <= 0 {
	  jutil.Elap("L382: Queue is empty", startms, jutil.Nowms())
    }
	u.perf.PrintStat(u.logFile)
	u.logFile.Sync()
	defer u.logFile.Close()
}
