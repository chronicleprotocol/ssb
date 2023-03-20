# SSB Tools

## SSB RPC Client

### Get all messages
```shell
./get_all_messages.sh SERVER_IP
```
### Get latest sequence numbers for feeds
```shell
jq -cr '.[]' ./feeds.json | ./get-latest-messages-from-feeds.sh SERVER_IP | jq -c '{key}*.value|{key,author,sequence}*.content|{key,author,sequence,type,version}'
```