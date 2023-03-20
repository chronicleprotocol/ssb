# SSB Tools

## SSB RPC Client

### Get all messages
```shell
./get_all_messages.sh "$SSB_SERVER_ADDRESS"
```
### Get latest sequence numbers for feeds
```shell
jq -cr '.[]' ./feeds.json | ./get-latest-messages-from-feeds.sh "$SSB_SERVER_ADDRESS" | jq -c '{key}*.value|{key,author,sequence}*.content|{key,author,sequence,type,version}'
```