# Test Bing Auto suggest feature
{  "id" : "003-test-bing-auto-suggest",   
   "verb" : "GET", 
   "uri" : "https://www.bing.com/AS/Suggestions?pt=page.home&mkt=en-us&qry=quant&cp=1&cvid=C4C0C8857D9B4D3C8F145A2D93D9426E", 
   "expected" : 200, 
   "rematch" : ".*quantum.*quantify", 
   "message" :"Bing auto suggest for quant"
}
#END