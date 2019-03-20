# See: https://golang.org/pkg/regexp/syntax/
#   For RE Expression matching syntax 
#
# Login and get a session token
{  "id" : "018-login",   
   "verb" : "GET", 
   "uri" : "https://abcdex1.org/login?user=test&pass=alive",
   "expected" : 200, 
   "keepBodyAs" : "sessionToken",
   "keepBodyDefault" : "D1A2C8F10A0AB0C0A1DEFAULT"
}
#END
#WAIT

# See: https://golang.org/pkg/regexp/syntax/
#   For RE Expression matching syntax 
#
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

# Test CGI Call to DEM Sub ssytem
{  "id" : "0182-DEM-CGI-CALL-GEOPOINT",   
   "verb" : "GET", 
   "uri" : "http://airsolarwater.com/dem/drains.svc?buri=gdata/dnt-rodrigues-island-50&offset=1591792&geo=-19.71542,63.34569", 
   "expected" : 200, 
   "rematch" : "-19.71542,63.34569,1.00.*-19.71736,63.34514", 
   "message" :"Checking CGI Geo Point"
}
#END


# Test CGI Call to DEM Sub system failure due to RE match failure
{  "id" : "0182-DEM-CGI-CALL-GEOPOINT-EXPECT-FAIL",   
   "verb" : "GET", 
   "uri" : "http://airsolarwater.com/dem/drains.svc?buri=gdata/dnt-rodrigues-island-50&offset=1591792&geo=-19.71542,63.34569", 
   "expected" : 200, 
   "rematch" : "-19.71542,63.34569,1.00.*-19.71736,63.XXX34514", 
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

