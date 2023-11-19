# This go app updates cloudflare "A" record every hour for DDNS

It fetches the Public IP from ip-api.com and uses cloudflare API to update to the new IP only if the IP has changed

## Create a .env in root dir and include the following
```
CLOUDFLARE_API_KEY=<Cloudflare API Key>
ZONE_ID=<Zone ID can be obtained from cloudflare dashboard>
ZONE_NAME=<Your website domain e.g. google.com>
```

## docker compose
```docker compose up```
