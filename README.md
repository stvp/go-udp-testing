go-udp-testing
==============

[![Build Status](https://travis-ci.org/stvp/go-udp-testing.png?branch=master)](https://travis-ci.org/stvp/go-udp-testing)

Provides UDP socket test helpers for Go.

[Documentation](http://godoc.org/github.com/stvp/go-udp-testing)

Example
-------

    package main

    import (
      "github.com/stvp/go-udp-testing"
      "testing"
      # ...
    )

    func TestStatsdReporting(t *testing.T) {
      # ...

      udp.SetAddr(":8125")

      udp.ShouldReceiveOnly(t, "mystat:2|g", func() {
        statsd.Gauge("mystat", 2)
      })

      udp.ShouldNotReceiveOnly(t, "mystat:1|c", func() {
        statsd.Gauge("bukkit", 2)
      })

      udp.ShouldReceive(t, "bar:2|g", func() {
        statsd.Gauge("foo", 2)
        statsd.Gauge("bar", 2)
        statsd.Gauge("baz", 2)
      })

      udp.ShouldNotReceive(t, "bar:2|g", func() {
        statsd.Gauge("foo", 2)
        statsd.Gauge("baz", 2)
      })
    }

