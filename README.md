# SSB Tools

## SSB RPC Client

### Examples
```shell
ssb-rpc-client log --limit 100 --keys \
 | jq -c '{timestamp,ts:.value.timestamp,d:(.timestamp - .value.timestamp),type:.value.content.type,key,author:.value.author}'
```