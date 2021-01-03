package main

import (
  "fmt"
  "gopkg.in/h2non/gentleman.v2"
  "gopkg.in/h2non/gentleman.v2/plugins/headers"
  "gopkg.in/h2non/gentleman.v2/plugins/body"
  "encoding/json"
  "os"
  "flag"
)
type Config struct{
  AuthToken string
  DnsType string
  FullDomainName string
  Ttl int
  Priority int
  Proxied bool
  ZoneId string
}

type ResponseResult struct{
  Success bool `json:"success"`
  Errors interface{} `json:"errors"`
  Messages interface{} `json:"messages"`
  Result ResultDetail `json:"result"`
}

type ResponseResults struct{
  Success bool `json:"success"`
  Errors interface{} `json:"errors"`
  Messages interface{} `json:"messages"`
  Result []ResultDetail `json:"result"`
}

type ResultDetail struct{
  Id string `json:"id"`
  Type string `json:"type"`
  Name string `json:"name"`
  Content string `json:"content"`
  Proxiable bool `json:"proxiable"`
  Ttl int `json:"ttl"`
  Locked bool `json:"locked"`
  ZoneId string `json:"zone_id"`
  ZoneName string `json:"zone_name"`
  CreatedOn string `json:"created_on"`
  ModifiedOn string `json:"modified_on"`
  Data interface{} `json:"data"`
  Meta interface{} `json:"meta"`
}

func main() {
  c := flag.String("c", "./config/config.json", "Specify the configuration file.")
  flag.Parse()
  config := LoadConfiguration(*c)
  _, result := GetRecords(config)
  fmt.Println(result)
}

func GetRecords(config Config) (ok bool, result ResponseResults){
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

func CreateRecord(config Config, ip string) (ok bool, result ResponseResult){
  cli := gentleman.New()

  // Define a custom header
  authorization := "Bearer " + config.AuthToken
  customHeaders := map[string]string{"Authorization": authorization, "Content-Type": "application/json"}
  cli.Use(headers.SetMap(customHeaders))

  data := map[string]interface{}{"type": config.DnsType, "name": config.FullDomainName, "content": ip,"ttl": config.Ttl,"priority":config.Priority,"proxied":config.Proxied}
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

func UpdateRecord(config Config, recordId string, ip string) (ok bool, result ResponseResult){
  fmt.Println(config)
  cli := gentleman.New()

  // Define a custom header
  authorization := "Bearer " + config.AuthToken
  customHeaders := map[string]string{"Authorization": authorization, "Content-Type": "application/json"}
  cli.Use(headers.SetMap(customHeaders))

  data := map[string]interface{}{"type": config.DnsType, "name": config.FullDomainName, "content": ip,"ttl": config.Ttl,"priority":config.Priority,"proxied":config.Proxied}
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
