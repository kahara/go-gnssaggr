# GNSS aggregator

In order to be able to monitor my local GNSS environment, aggregate
[gpsd](https://gpsd.gitlab.io/gpsd/)
[JSON](https://gpsd.gitlab.io/gpsd/gpsd_json.html)
reports for consumption by
[Prometheus](https://prometheus.io/docs/instrumenting/writing_exporters/#metrics).

## Running

FIXME link to Docker image

## Metrics

```
gnssaggr_sky_gdop
gnssaggr_sky_hdop
gnssaggr_sky_pdop
gnssaggr_sky_tdop
gnssaggr_sky_vdop
gnssaggr_sky_xdop
gnssaggr_sky_ydop
gnssaggr_sky_satellites{gnssid="0.."}
gnssaggr_tpv_althae_iqr
gnssaggr_tpv_althae_mad
gnssaggr_tpv_althae_max
gnssaggr_tpv_althae_min
gnssaggr_tpv_althae_q1
gnssaggr_tpv_althae_q2
gnssaggr_tpv_althae_q3
gnssaggr_tpv_althae_stddev
gnssaggr_tpv_althae_variance
gnssaggr_tpv_altmsl_iqr
gnssaggr_tpv_altmsl_mad
gnssaggr_tpv_altmsl_max
gnssaggr_tpv_altmsl_min
gnssaggr_tpv_altmsl_q1
gnssaggr_tpv_altmsl_q2
gnssaggr_tpv_altmsl_q3
gnssaggr_tpv_altmsl_stddev
gnssaggr_tpv_altmsl_variance
gnssaggr_tpv_ept_iqr
gnssaggr_tpv_ept_mad
gnssaggr_tpv_ept_max
gnssaggr_tpv_ept_min
gnssaggr_tpv_ept_q1
gnssaggr_tpv_ept_q2
gnssaggr_tpv_ept_q3
gnssaggr_tpv_ept_stddev
gnssaggr_tpv_ept_variance
gnssaggr_tpv_epv_iqr
gnssaggr_tpv_epv_mad
gnssaggr_tpv_epv_max
gnssaggr_tpv_epv_min
gnssaggr_tpv_epv_q1
gnssaggr_tpv_epv_q2
gnssaggr_tpv_epv_q3
gnssaggr_tpv_epv_stddev
gnssaggr_tpv_epv_variance
gnssaggr_tpv_epx_iqr
gnssaggr_tpv_epx_mad
gnssaggr_tpv_epx_max
gnssaggr_tpv_epx_min
gnssaggr_tpv_epx_q1
gnssaggr_tpv_epx_q2
gnssaggr_tpv_epx_q3
gnssaggr_tpv_epx_stddev
gnssaggr_tpv_epx_variance
gnssaggr_tpv_epy_iqr
gnssaggr_tpv_epy_mad
gnssaggr_tpv_epy_max
gnssaggr_tpv_epy_min
gnssaggr_tpv_epy_q1
gnssaggr_tpv_epy_q2
gnssaggr_tpv_epy_q3
gnssaggr_tpv_epy_stddev
gnssaggr_tpv_epy_variance
gnssaggr_tpv_lat_iqr
gnssaggr_tpv_lat_mad
gnssaggr_tpv_lat_max
gnssaggr_tpv_lat_min
gnssaggr_tpv_lat_q1
gnssaggr_tpv_lat_q2
gnssaggr_tpv_lat_q3
gnssaggr_tpv_lat_stddev
gnssaggr_tpv_lat_variance
gnssaggr_tpv_lon_iqr
gnssaggr_tpv_lon_mad
gnssaggr_tpv_lon_max
gnssaggr_tpv_lon_min
gnssaggr_tpv_lon_q1
gnssaggr_tpv_lon_q2
gnssaggr_tpv_lon_q3
gnssaggr_tpv_lon_stddev
gnssaggr_tpv_lon_variance
```

## Rationale

In my home Kubernetes cluster, one of the nodes is running
[gpsd](https://github.com/kahara/docker-gpsd)
and has a
[cheap GNSS receiver](https://www.aliexpress.com/item/32816656706.html)
connected to it. Recording what the receiver sees and presenting that in a Grafana dashboard sounds like a nice idea.

Because it's not possible to make Prometheus scrape the GNSS metrics e.g. once or twice every second, it
Makes Sense&trade; to collect gpsd reports for each full minute (seconds 0&hellip;59; positive leap seconds are
ignored) and present  the gathered data to Prometheus as  statistics. I'm aware of what Prometheus
[says](https://prometheus.io/docs/instrumenting/writing_exporters/#drop-less-useful-statistics)
about providing stat-type metrics, but this case is a bit different as the values in the gpsd-produced reports
fluctuate, sometime wildly, during that one-minute collection period and this short term fluctuation is something
that interests me, along with any long term trends.

## The Galmon project

In my home cluster there's also another cheap receiver connected to another Raspberry Pi node running
[berthubert/galmon](https://hub.docker.com/r/berthubert/galmon),
feeding to
[galmon.eu](https://galmon.eu/observers.html).
If you can, please consider
[joining the network](https://berthub.eu/articles/posts/galmon-project/#how-to-join-in).
This as a heads-up in case anyone reading this hasn't heard of the Galmon project.
