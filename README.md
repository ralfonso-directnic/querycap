# Querycap

Captures mysql queries and counts the tables queried.

## Usage

```
Usage of querycap:
  -device string
    	Device to listen to (any is default) (default "any")
  -logfile string
    	Log file location (default "/var/log/querycap.log")
  -port int
    	Mysql Port (default 3306)
  -promport string
    	Prometheus Port (default "9224")
```

This was created to monitor a mysql database and identify any tables or databases that were not getting queried. 

## Install

Copy to your path and run as above

## Building

This code makes use of pcap c library, you must have that installed when building.