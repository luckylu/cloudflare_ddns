package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "gopkg.in/h2non/gentleman.v2"
  "gopkg.in/h2non/gentleman.v2/plugins/body"
  "gopkg.in/h2non/gentleman.v2/plugins/headers"
  "os"
  "time"
)

type Config struct {
  ApiToken                string
  DnsType                 string
  FullQualifiedDomainName string
  Ttl                     int
  Priority                int
  Proxied                 bool
  ZoneId                  string
  Interval                int64
  GetIpApi                string
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
  ip := GetIp(config)
  _, records := GetRecords(config)
  recordCreated := false
  var recordId string
  for _, result := range records.Result {
    if result.Name == config.FullQualifiedDomainName {
      recordCreated = true
      recordId = result.Id
    }
  }
  // done := make(chan bool, 1)
  if recordCreated {
    LoopUpdateRecord(config, recordId)
  } else {
    _, createResult := CreateRecord(config, ip)
    LoopUpdateRecord(config, createResult.Result.Id)
  }
  // <-done
}

func LoopUpdateRecord(config Config, recordId string) {
  // ticker := time.NewTicker(time.Duration(config.Interval) * time.Second)
  // defer ticker.Stop()
  for {
    UpdateRecord(config, recordId)
    time.Sleep(time.Duration(config.Interval) * time.Second)
  }
}

func GetIp(config Config) (ip string) {
  defer func() {
    if err := recover(); err != nil {
      fmt.Println("err get ip:", err)
      ip = GetIp(config)
    }
  }()
  cli := gentleman.New()
  res, err := cli.Request().Method("GET").URL(config.GetIpApi).Send()
  if err != nil {
    panic(err.Error())
  }
  ip = res.String()
  return
}

func GetRecords(config Config) (ok bool, result ResponseResults) {
  cli := gentleman.New()
  authorization := "Bearer " + config.ApiToken
  customHeaders := map[string]string{"Authorization": authorization, "Content-Type": "application/json"}
  cli.Use(headers.SetMap(customHeaders))
  res, err := cli.Request().Method("GET").URL("https://api.cloudflare.com/client/v4/zones/" + config.ZoneId + "/dns_records").Send()
  if err != nil {
    ok = false
    return
  }
  if !res.Ok {
    // json.Unmarshal([]byte(res.String()), &result)
    // ok = false
    panic(res.String())
  }

  json.Unmarshal([]byte(res.String()), &result)
  ok = true
  return
}

func CreateRecord(config Config, ip string) (ok bool, result ResponseResult) {
  cli := gentleman.New()

  // Define a custom header
  authorization := "Bearer " + config.ApiToken
  customHeaders := map[string]string{"Authorization": authorization, "Content-Type": "application/json"}
  cli.Use(headers.SetMap(customHeaders))

  data := map[string]interface{}{"type": config.DnsType, "name": config.FullQualifiedDomainName, "content": ip, "ttl": config.Ttl, "priority": config.Priority, "proxied": config.Proxied}
  cli.Use(body.JSON(data))

  // Perform the request
  res, err := cli.Request().Method("POST").URL("https://api.cloudflare.com/client/v4/zones/" + config.ZoneId + "/dns_records").Send()
  if err != nil {
    ok = false
    return
  }
  if !res.Ok {
    // json.Unmarshal([]byte(res.String()), &result)
    // ok = false
    // return
    panic(res.String())
  }

  json.Unmarshal([]byte(res.String()), &result)
  ok = true
  return
}

func UpdateRecord(config Config, recordId string) (ok bool, result ResponseResult) {
  ip := GetIp(config)
  fmt.Println("update record", ip)
  fmt.Println(time.Now())
  cli := gentleman.New()

  // Define a custom header
  authorization := "Bearer " + config.ApiToken
  customHeaders := map[string]string{"Authorization": authorization, "Content-Type": "application/json"}
  cli.Use(headers.SetMap(customHeaders))

  data := map[string]interface{}{"type": config.DnsType, "name": config.FullQualifiedDomainName, "content": ip, "ttl": config.Ttl, "priority": config.Priority, "proxied": config.Proxied}
  cli.Use(body.JSON(data))

  // Perform the request
  res, err := cli.Request().Method("PUT").URL("https://api.cloudflare.com/client/v4/zones/" + config.ZoneId + "/dns_records/" + recordId).Send()
  fmt.Println(res)
  if err != nil {
    ok = false
    return
  }
  if !res.Ok {
    // json.Unmarshal([]byte(res.String()), &result)
    // ok = false
    // return
    panic(res.String())
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
