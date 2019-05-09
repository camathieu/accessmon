Access Mon
==========

Access Mon is able to display statistics about HTTP server logs in the [W3C Common Log File Format](https://www.w3.org/Daemon/User/Config/Logging.html)

```
$ ./accessmon --help
Usage of ./accessmon:
  -logfile string
        log file path (default "/tmp/access.log")
  -offline
        offline mode ( cat )
  -refresh duration
        screen refresh interval ( online mode only ) (default 10s)
  -threshold float
        total request per second moving average alerting threshold (default 10)
  -window duration
        total request per second moving average alerting window (default 2m0s)
```

By default access mon will seek to the end of the logfile and tail -f,
displaying every 10 seconds statistics about the last 10 seconds of received logs.
of received logs. If no logs have been received a warning message will be displayed.

If configured the program will detect and alert when the total number of requests
raise above the configured threshold ( 10 request per second by default ) for
the consecutive configured period of time ( 2 minutes by default ).

In offline mode the program will open and read the whole logfile ( cat ) and
run the alert detection algorithm.

Building from sources
=====================

Having a standard Golang environment installed simply run make
The accessmon binary will be available in the cmd directory

```
make && cd cmd && ./accessmon
```

Running linters and tests

```
make lint && make test
```

Building docker image
```
make docker
```