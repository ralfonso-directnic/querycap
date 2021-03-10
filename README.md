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

## Prometheus

Prometheus is available on host:promport/metrics

```
qry_db_table{db="",table=""} 3
qry_db_table{db="Example",table="Test"} 24
qry_db_table{db="Example",table="Test2"} 13
qry_db_table{db="Watching",table="Entry"} 4
# HELP qry_total The count of Total Queries
# TYPE qry_total counter
qry_total 44
```