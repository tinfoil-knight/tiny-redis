# tiny-redis

This project intends to be a rough implementation of the in-memory data store: [Redis](https://redis.io/)

Note:
- The parser implements a subset of [RESP3](https://github.com/antirez/RESP3/blob/master/spec.md) without the Attribute, Push and Streamed data types.

## Getting Started

### Pre-requisites
- [Go >=1.1](https://golang.org/)

### Running locally

```bash
make run
```
> Note: The examples assume that the TCP server is running on localhost at port 8001

You can run commands through `netcat` following the Redis protocol: 

```bash
echo -e '*1\r\n$4\r\nPING\r\n' | nc localhost 8001
```

Or you can start a redis client using the Redis CLI: `redis-cli -p 8001` and then use it in Interactive mode.

```bash
127.0.0.1:8001> SET hello 3
OK
```

### Running tests

```bash
make test
```

## Author

- ***Kunal Kundu*** [@tinfoil-knight](https://github.com/tinfoil-knight)

## Acknowledgements

- [Exotel](https://exotel.com/) for their [tech-challenge](https://exotel.com/about-us/exotel-tech-challenge/) which gave me the idea to build this.
- [Redis protocol specification](https://redis.io/topics/protocol) and [RESP3 spec](https://github.com/antirez/RESP3/blob/master/spec.md) for documenting the Redis protocol.

## Appendix
**List of Allowed Commands**

- Connection: `PING`, `ECHO`
- Keys: `DEL`, `EXISTS`
- Strings: `GET`, `SET`, `GETDEL`, `INCR`, `DECR`, `INCRBY`, `DECRBY`, `APPEND`
- Server: `SAVE`

> Note: Some commands may not support all options available in Redis 6.



