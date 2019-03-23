# Add Data to MDS repository including using Custom 
# headers.    Then read back portions of that data
# to demonstrate that the data previously written 
# is available.   Demonstrates custom headers and 
# RE positive matching and RE Negative matching
# Also demonstrates advanced use of #WAIT to ensure
# CRUD oprations are complete before the reads are
# attempted
#
# See:  https://www.freeformatter.com/json-escape.html   for online JSON escape tool for post body.
# See: https://jsonlint.com/  for online tool to validate the JSON
# See: GenericHTTPTestClient.go for utility than runs this test file.


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


#WAIT
#  Pause causes the test driver to wait until all threads are idle before running next line. 


{  "id" : "0184G",   "verb" : "GET",   "uri" : "http://127.0.0.1:9601/mds/test/1817127X3",  "expected" : 200, 
   "rematch" : ".*JOHN MUIR.*", 
   "message" :"Read after write"
}
#END

#WAIT
#  Pause causes the test driver to wait until all threads are idle before running next line. 


# Test simple retrieve of a record that should exist.
{  "id" : "0184M",   "verb" : "GET",   "uri" : "http://127.0.0.1:9601/mds/test/1817127X3",  "expected" : 200, 
   "rematch" : ".*JOHN MUIR.*", 
   "message" :"Read after write existing record"
}
#END


# Test to Demonstrate failure when we get Data and shouldn't have received any.
# Should fail
{  "id" : "0184F",   "verb" : "GET",   "uri" : "http://127.0.0.1:9601/mds/test/1817127X3",  "expected" : 200, 
   "Renomatch" : ".*JOHN MUIR.*", 
   "message" :"Check re No Match functionality should fail"
}
#END


# Test to verify record we do not expect to exist fails.
# Should Suceed
{  "id" : "0185",   "verb" : "GET",   "uri" : "http://127.0.0.1:9601/mds/test/1817127X5",  "expected" : 404, 
   "Message" :"Check expected 404 on known bad key",
   "Renomatch" : "*.JOHN*."
}
#END

# Demontrate a failed test where the status code
# is not what was expected.
# Should suceed
{  "id" : "0186",   "verb" : "GET",   "uri" : "http://127.0.0.1:9601/mds/test/1817127X5",  "expected" : 519, 
   "Message" :"Check unepected response code "
}
#END



#WAIT
#  Wait until after the reads above have complete or could delete the records out from under them
# Should Succeed
{  "id" : "0187",   
   "verb" : "DELETE",   
   "uri" : "http://127.0.0.1:9601/mds/test/1817127X3",  
   "message" : "Delete a JSON record",
   "expected" : 200,
   "xyz" : 983
}
#END

#WAIT  
#  Wait for the delete to complete so we can check for the 404
#GET after delete of same ID should fail
{  "id" : "0188",
   "verb" : "GET",
   "uri" : "http://127.0.0.1:9601/mds/test/1817127X3",
   "expected" : 404,  
   "Message" :"Read after delete"
}
#END



# Second delete of same ID should fail

#WAIT
{  "id" : "0189",   
   "verb" : "DELETE",   
   "uri" : "http://127.0.0.1:9601/mds/test/1817127X3",  
   "message" : "Delete a previously deleted record",
   "expected" : 404 
}
#END
