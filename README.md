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
Configuration is minimal and can be in a yaml file or envirnoment variables which override the config file.  
    
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
```go build bfpd
./bfpd -d -i -c /path/to/config.yaml   
where
  -d use gorm debugging   
  -i initialize a database schema
  -c configuration file to use (defaults to ./config.yaml )   
  ```
