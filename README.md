# Cloudflare DDNS

## Configuration

create a file named config.json with following fields in current directory

```ruby
{
  "ApiToken": "abcde", #create your token at https://dash.cloudflare.com/profile/api-tokens
  "DnsType": "A",
  "FullQualifiedDomainName": "my.example.com", # change to your domain
  "Ttl": 120,
  "Priority": 10,
  "Proxied": false,
  "ZoneId": "abcde", #view your Zone Id at domain overview page
  "Interval": 180, # update dns record every 180 seconds
  "GetIpApi": "http://members.3322.org/dyndns/getip" # get ip api or any other api provider that return ip address
}

```

## How to run

download the latest program from releases

```ruby
./cloudflare_ddns #if config file is not in current directory use -c to specify where the config file locate
```
