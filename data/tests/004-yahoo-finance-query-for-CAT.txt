{  "id" : "#004-yahoo-finance-query-for-CAT",   
   "verb" : "GET", 
   "uri" : "https://query1.finance.yahoo.com/v8/finance/chart/CAT?region=US&lang=en-US&includePrePost=false&interval=2m&range=1d&corsDomain=finance.yahoo.com&.tsrc=finance", 
   "expected" : 200, 
   "rematch" : ".*symbol.*CAT.*timezone",    
   "message" :"query yahoo finance for CAT"
}
#END
