package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "gopkg.in/h2non/gentleman.v2"
  "gopkg.in/h2non/gentleman.v2/plugins/body"
  "gopkg.in/h2non/gentleman.v2/plugins/headers"
  "os"
)

type Config struct {
  AuthToken      string
  DnsType        string
  FullDomainName string
  Ttl            int
  Priority       int
  Proxied        bool
  ZoneId         string
}

type ResponseResult struct {
  Success  bool         `json:"success"`
  Errors   interface{}  `json:"errors"`
  Messages interface{}  `json:"messages"`
  Result   ResultDetail `json:"result"`
}

type ResponseResults struct {
  Success  bool           `json:"success"`
  Errors   interface{}    `json:"errors"`
  Messages interface{}    `json:"messages"`
  Result   []ResultDetail `json:"result"`
}

type ResultDetail struct {
  Id         string      `json:"id"`
  Type       string      `json:"type"`
  Name       string      `json:"name"`
  Content    string      `json:"content"`
  Proxiable  bool        `json:"proxiable"`
  Ttl        int         `json:"ttl"`
  Locked     bool        `json:"locked"`
  ZoneId     string      `json:"zone_id"`
  ZoneName   string      `json:"zone_name"`
  CreatedOn  string      `json:"created_on"`
  ModifiedOn string      `json:"modified_on"`
  Data       interface{} `json:"data"`
  Meta       interface{} `json:"meta"`
}

func main() {
  c := flag.String("c", "./config.json", "Specify the configuration file.")
  flag.Parse()
  config := LoadConfiguration(*c)
  ip := GetIp()
  _, records := GetRecords(config)
  recordCreated := false
  var recordId string
  for _, result := range records.Result {
    if result.Name == config.FullDomainName {
      recordCreated = true
      recordId = result.Id
    }
  }
  if recordCreated {
    UpdateRecord(config, recordId, ip)
  } else {
    CreateRecord(config, ip)
  }
}

func GetIp() (ip string) {
  cli := gentleman.New()
  res, err := cli.Request().Method("GET").URL("http://members.3322.org/dyndns/getip").Send()
  if err != nil {
    panic(err.Error())
  }
  ip = res.String()
  return
}

func GetRecords(config Config) (ok bool, result ResponseResults) {
  cli := gentleman.New()
  authorization := "Bearer " + config.AuthToken
  customHeaders := map[string]string{"Authorization": authorization, "Content-Type": "application/json"}
  cli.Use(headers.SetMap(customHeaders))
  res, err := cli.Request().Method("GET").URL("https://api.cloudflare.com/client/v4/zones/" + config.ZoneId + "/dns_records").Send()
  if err != nil {
    ok = false
    return
  }
  if !res.Ok {
    json.Unmarshal([]byte(res.String()), &result)
    ok = false
    return
  }

  json.Unmarshal([]byte(res.String()), &result)
  ok = true
  return
}

func CreateRecord(config Config, ip string) (ok bool, result ResponseResult) {
  cli := gentleman.New()

  // Define a custom header
  authorization := "Bearer " + config.AuthToken
  customHeaders := map[string]string{"Authorization": authorization, "Content-Type": "application/json"}
  cli.Use(headers.SetMap(customHeaders))

  data := map[string]interface{}{"type": config.DnsType, "name": config.FullDomainName, "content": ip, "ttl": config.Ttl, "priority": config.Priority, "proxied": config.Proxied}
  cli.Use(body.JSON(data))

  // Perform the request
  res, err := cli.Request().Method("POST").URL("https://api.cloudflare.com/client/v4/zones/" + config.ZoneId + "/dns_records").Send()
  if err != nil {
    ok = false
    return
  }
  if !res.Ok {
    json.Unmarshal([]byte(res.String()), &result)
    ok = false
    return
  }

  json.Unmarshal([]byte(res.String()), &result)
  ok = true
  return
}

func UpdateRecord(config Config, recordId string, ip string) (ok bool, result ResponseResult) {
  cli := gentleman.New()

  // Define a custom header
  authorization := "Bearer " + config.AuthToken
  customHeaders := map[string]string{"Authorization": authorization, "Content-Type": "application/json"}
  cli.Use(headers.SetMap(customHeaders))

  data := map[string]interface{}{"type": config.DnsType, "name": config.FullDomainName, "content": ip, "ttl": config.Ttl, "priority": config.Priority, "proxied": config.Proxied}
  cli.Use(body.JSON(data))

  // Perform the request
  res, err := cli.Request().Method("PUT").URL("https://api.cloudflare.com/client/v4/zones/" + config.ZoneId + "/dns_records/" + recordId).Send()
  if err != nil {
    ok = false
    return
  }
  if !res.Ok {
    json.Unmarshal([]byte(res.String()), &result)
    ok = false
    return
  }

  json.Unmarshal([]byte(res.String()), &result)
  ok = true
  return
}

func LoadConfiguration(file string) (conf Config) {
  configFile, err := os.Open(file)
  defer configFile.Close()
  if err != nil {
    fmt.Println(err.Error())
  }
  jsonParser := json.NewDecoder(configFile)
  jsonParser.Decode(&conf)
  return
}
