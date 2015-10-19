# RRD-FIXPEAKS

* Homepage: https://github.com/jbuchbinder/rrd-fixpeaks
* Twitter: [@jbuchbinder](https://twitter.com/jbuchbinder)
* [![Build
Status](https://secure.travis-ci.org/jbuchbinder/rrd-fixpeaks.png)](http://travis-ci.org/jbuchbinder/rrd-fixpeaks)
* [![Gobuild Download](http://gobuild.io/badge/github.com/jbuchbinder/rrd-fixpeaks/downloads.svg)](http://gobuild.io/github.com/jbuchbinder/rrd-fixpeaks)

## USAGE

```
rrd-fixpeaks -threshold=80 RRDFILE.rrd
rrd-fixpeaks -absabove=1e10 RRDFILE.rrd

Usage of ./rrd-fixpeaks:
  -absabove float
    	If not -1, every value above this will be removed (default -1)
  -dryrun
    	Dry run flag (don't write)
  -mindiff float
    	Minimum difference above average
  -multiplier float
    	Factor which max must outstrip average (default 2)
  -rrdtool string
    	Path to rrdtool executable (default "rrdtool")
  -threshold float
    	Threshold percentage above avg above which values should be clipped
```

## BUILDING

This requires **rrdtool** be installed and executable on your PATH. It
was tested with 1.4.7, so your mileage may vary.

```
go build
```

