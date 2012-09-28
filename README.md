# RRD-FIXPEAKS

* Homepage: https://github.com/jbuchbinder/rrd-fixpeaks
* Twitter: [@jbuchbinder](https://twitter.com/jbuchbinder)
* [![Build
Status](https://secure.travis-ci.org/jbuchbinder/rrd-fixpeaks.png)](http://travis-ci.org/jbuchbinder/rrd-fixpeaks)

## USAGE

```
rrd-fixpeaks -threshold=80 RRDFILE.rrd
rrd-fixpeaks -absabove=1e10 RRDFILE.rrd
```

## BUILDING

This requires **rrdtool** be installed and executable on your PATH. It
was tested with 1.4.7, so your mileage may vary.

```
go build
```

