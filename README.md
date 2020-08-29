# cflsh
 CloudFlare Worker Shell

 - Pub-sub command and response over CloudFlare workers using [KV pairs]("https://developers.cloudflare.com/workers/runtime-apis/kv")
 - Reasonable Command Latency: 3-6 sec
 - Command history in KV buckets
 - One per command output channel.
 

### Usage

Red side client
`$ ./cflsh red`
```
✔ command: █
Setting up bucket: commandy9CgKEFRoCzxWcwRtD7LcZ
Setting up output bucket: outputy9CgKEFRoCzxWcwRtD7LcZ
Got output: 
`
LICENSE
README.md
cmd
common
`
```

Blue side client
`$ ./cflsh blue`
```
Got name of bucket: commandy9CgKEFRoCzxWcwRtD7LcZ
Got payload: ls
Executing command
OK: payload len 58
Setting up response bucket: outputy9CgKEFRoCzxWcwRtD7LcZ
```

`$ ./cflsh buckets`
```
 >> list 
    ... list of command and output buckets
 >> dump 
    ... contents of command and output buckets
```

### Note:

TODO. For now - manually put your CF API key, Email and Account ID 
in `./common/auth.go`:
```
package common

const (
        ApiKey = "82a0b14145b"
        AccountID = "aeccc94dfbf"
        EmailUser="cl@th.co"
)
```

### Diagram

![cflsh](https://github.com/dsnezhkov/cflsh/blob/master/cfl.png)
