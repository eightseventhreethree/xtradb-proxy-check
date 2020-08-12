# XtraDB-Proxy-Check

A Go replacement for the Shell script that was provided to check health of a PXC node. This is designed to replace the HAProxy check use case, and not the shell return value use case.

[percona-clustercheck](https://github.com/olafz/percona-clustercheck)

Differences:
* No more dependency on mysql shell client, doesn't require re-write to handle XtraDB in a container.
* No more dependency on xinetd.

## Pre-requisite:
```
GRANT PROCESS ON *.* TO 'clustercheckuser'@'127.0.0.1' IDENTIFIED BY 'clustercheckpassword!'
```

## Build:
```
make build
```

## Run:
```
make run
```

## Overrides:
```
export CLUSTERCHECK_MYSQL_USERNAME='clustercheckuser'
export CLUSTERCHECK_MYSQL_PASSWORD='clustercheckpassword!'
export CLUSTERCHECK_MYSQL_HOSTNAME='127.0.0.1'
export CLUSTERCHECK_MYSQL_PORT=3306
export CLUSTERCHECK_API_PORT=9200
export CLUSTERCHECK_AVAILABLE_WHEN_DONOR=true
export CLUSTERCHECK_AVAILABLE_WHEN_READ_ONLY=false
```

## Contributing:
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License:
[GPLv2](https://www.gnu.org/licenses/old-licenses/gpl-2.0.en.html)