# bfpd-api
A small web server providing a REST API for the Branded Food Products Database.  
### Installation
Clone this repo into your [go workspace](https://golang.org/doc/code.html)  
Install supporting packages as needed:     
*[gin framework](https://github.com/gin-gonic/gin)   
*[go gorm](http://jinzhu.me/gorm/)    
*[gin-jwt](https://github.com/appleboy/gin-jwt)    
*[bcrypt](https://godoc.org/golang.org/x/crypto/bcrypt)       
*[gopkg.in/yaml.v2](http://gopkg.in/yaml.v2)   
### Configuration
Configuration is minimal and can be in a YAML file or envirnoment variables which override the config file.  
    
YAML    
```
db: <schema name>
user: <dbuser>
pw: <dbpassword>
```
Environment   
```
BFPD_DB=Schema_name
BFPD_USER=Database_user
BFPD_PW=Database_user_password
```
### Running
```
go build bfpd.go
./bfpd -d -i -c /path/to/config.yaml   
where
  -d use gorm debugging   
  -i initialize a database schema
  -c configuration file to use (defaults to ./config.yaml )  
  -p TCP port to run server (defaults to 8080)
  ```
### Usage
Authentication is required for POST, PUT and DELETE.  Use the login handler to obtain a token which then must be sent in the Authorization header as shown in the examples below.     
>Authenticate and obtain JWT token:
``` 
curl -X POST -H "Content-type:application/json" -d '{"password":"your-password","username":"your-user-name"}' http://localhost:8080/ndb/api/v1/login 
```
or if you prefer http:
```
http -v --json POST localhost:8080/ndb/api/v1/login username=your-password password=your-user-name
```
>Add foods to the database:   
```
curl -X POST -H "Content-type:application/json" -H "Authorization:Bearer <your jwt token>" \
-d '{"ndbno":"45001535","name":"STEAK HOUSE STEAK SAUCE, UPC: 5051379020064","manu":{"name":"FRESH & EASY"},"fg":{"cd":"4500"}, \
"ingredients":{"desc":"TOMATO PUREE, ONION PUREE. SUGAR, MOLASSES, DISTILLED VINEGAR, HORSERADISH, SALT, SOYBEAN OIL, ORANGE JUICE CONCENTRATE, LEMON JUICE CONCENTRATE, ANCHOVY PASTE(ANCHOVY OLIVE OIL, ACETIC ACID), ROASTED GARLIC PUREE(GARLIC, WATER, NATURAL FLAVOR, CITRIC ACID), NATURAL FLAVOR, JALAPENO PUREE(JALAPENO CHILE, DISTILLED VINEGAR, SALT), SOY SAUCE(WATER, WHEAT, SOYBEAN, SALT), CHILI POWDER, MUSTARD FLOUR, XANTHAN GUM, BLACK PEPPER, CARAMEL COLORING, CLOVE POWDER.","updated":"01/02/2014"}, \
"measures":[{"seq":1,"unit":"Tbsp","amt":1.0,"weight":15.0}], \
"nutrients":[{"nutno":203,"value":0.0,"dp":0,"source":{"code":"8"},"deriv":{"code":"LCBF"}},{"nutno":204,"value":0.0,"dp":0,"source":{"code":"8"},"deriv":{"code":"LCBF"}},{"nutno":205,"value":20.0,"dp":0,"source":{"code":"8"},"deriv":{"code":"LCBF"}},{"nutno":208,"value":100.0,"dp":0,"source":{"code":"8"},"deriv":{"code":"LCBF"}},{"nutno":269,"value":20.0,"dp":0,"source":{"code":"8"},"deriv":{"code":"LCBF"}},{"nutno":291,"value":0.0,"dp":0,"source":{"code":"8"},"deriv":{"code":"LCBF"}},{"nutno":301,"value":133.0,"dp":0,"source":{"code":"8"},"deriv":{"code":"LCBF"}},{"nutno":303,"value":4.8,"dp":0,"source":{"code":"8"},"deriv":{"code":"LCBF"}},{"nutno":307,"value":833.0,"dp":0,"source":{"code":"8"},"deriv":{"code":"LCBF"}},{"nutno":318,"value":2000.0,"dp":0,"source":{"code":"8"},"deriv":{"code":"LCBF"}},{"nutno":401,"value":16.0,"dp":0,"source":{"code":"8"},"deriv":{"code":"LCBF"}},{"nutno":601,"value":0.0,"dp":0,"source":{"code":"8"},"deriv":{"code":"LCBF"}},{"nutno":605,"value":0.0,"dp":0,"source":{"code":"8"},"deriv":{"code":"LCBF"}},{"nutno":606,"value":0.0,"dp":0,"source":{"code":"8"},"deriv":{"code":"LCBF"}}]' http://localhost:8000/ndb/api/v1/food
```
>Delete a food by database id:
```
curl -X DELETE -H "Content-type:application/json" -H "Authorization:Bearer <your jwt token>" \
http://localhost:8080/ndb/api/v1/food/<food-db-id>
```
or
 ```
 http DELETE --json localhost:8080/ndb/api/v1/food "Authorization:Bearer <your jwt token>" id=<food-db-id>
 ```
>Fetch food by ndbno:
curl -X GET http://localhost:8000/ndb/api/v1/ndb/45001535
```
>Fetch a list of foods:
```
http GET localhost:8080/ndb/api/v1/food/ max=50 offset=50
```
>Add a nutrient
```
curl -X POST -H "Content-type:application/json" -H "Authorization:Bearer <your jwt token>" \
-d '{"desc":"Total lipid (fat)","nutno":204,"Decimalpoint":2,"Tagname":"FAT","Srnutorder":800,"Unit":{"Unit":"g"}}' \ http://localhost:8080/ndb/api/v1/nutrient
```
