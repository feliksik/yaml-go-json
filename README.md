# yaml-go-json

yaml is great for manual editing (using anchors is especially powerful), 
but a lot of tools assume json config (for good reasons). 

This tool converts yaml input to json. It always outputs to stdout.

Why is this useful: 

* expand yaml anchors, which are great for managing files manually
* gather yaml data from various sources and merge them. 
* output in yaml, json or json-pretty

## usage

Check usage:
```yaml-go-json --help```

Convert mydoc.yml: 

```cat mydoc.yml | ./yaml-go-json```

Convert mydoc.yml, and add some yaml from the web. This happens to be json, which is 
also valid yaml. We can then query the result with `jq`.

```cat mydoc.yml | ./yaml-go-json -s stdin:// -s http://example.com/myfile.json | jq '.somefield' ```

You can output json, json-pretty, or go back to yaml after merging. 

```cat mydoc.yml | ./yaml-go-json -s stdin:// -s http://example.com/myfile.json -f json-pretty ```

## notes

* if you have '<', '>' or '&' in your strings, it will be converted to unicode. This is per Golang's Unmarshall. If you don't want this, pipe to `jq` (or `jq -cM`)

## installing
make build

