# Test simple retrieve of a record that should exist.
{  "id" : "0182772-airsolarwater-dem-index-contains-solomon",   
   "verb" : "GET", 
   "uri" : "http://airsolarwater.com/dem/", 
   "expected" : 200, 
   "rematch" : ".*Solomon.*Nendo.*", 
   "message" :"Read after write existing record"
}
#END
