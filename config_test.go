package main
import(
  "bfpd/model"
  "fmt"
  "os"
  "gopkg.in/yaml.v2"
  "testing"
)
var config = `
  db: food2
  user: tester
  pw: testpw
`
// Configuration tests
func TestConfig( t *testing.T) {
  var cs bfpd.Config
  cc := []byte(config)
  yaml.Unmarshal(cc, &cs)
  rc,msg:=chkConfig(&cs)
  if rc != true {
    t.Errorf(msg)
  }
}
func TestEnvConfig(t *testing.T) {
  os.Setenv("BFPD_USER_TEST","tester")
  os.Setenv("BFPD_DB_TEST","food2")
  os.Setenv("BFPD_PW_TEST","testpw")
  var cs bfpd.Config
  cs.User=os.Getenv("BFPD_USER_TEST")
  cs.Pw=os.Getenv("BFPD_PW_TEST")
  cs.Db=os.Getenv("BFPD_DB")
  rc,msg:=chkConfig(&cs)
  if rc != true {
    t.Errorf(msg)
  }
  os.Setenv("BFPD_USER_TEST","")
  os.Setenv("BFPD_DB_TEST","")
  os.Setenv("BFPD_PW_TEST","")
}
// test env override by invoking a fail by setting a bad password in the env
func TestEnvOverride(t *testing.T) {
  var cs bfpd.Config
  var rc bool
  var msg string
  cc := []byte(config)
  yaml.Unmarshal(cc, &cs)
  rc,msg=chkConfig(&cs)
  os.Setenv("BFPD_PW_TEST","testpww")
  cs.Pw=os.Getenv("BFPD_PW_TEST")
  rc,msg=chkConfig(&cs)
  if rc == true {
    t.Errorf("Should be false but got true %s",msg)
  } else {
    os.Setenv("BFPD_PW_TEST","testpw")
    cs.Pw=os.Getenv("BFPD_PW_TEST")
    rc,msg=chkConfig(&cs)
  }
  if rc == false {
      t.Errorf(msg)
  }
  os.Setenv("BFPD_PW_TEST","")
}
// check to see if the config matches the values we've assigned
func chkConfig(cs *bfpd.Config) (bool,string) {
  if "food2" != cs.Db {
    return false,"Wrong database name "+cs.Db
  } else if "tester" != cs.User {
    return false,"Wrong user name "+cs.User
  } else if "testpw" != cs.Pw {
    return false,"Wrong users password "+cs.Pw
  } else {
    c := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8&parseTime=True&loc=Local", cs.User, cs.Pw, cs.Db)
    if "tester:testpw@tcp(127.0.0.1:3306)/food2?charset=utf8&parseTime=True&loc=Local" != c {
      return false,"Connection string does not match"
    }
  }
  return true,""
}
