# reckon

[![Build Status][actions-badge]][actions-url]
[![MIT licensed][mit-badge]][mit-url]

[actions-badge]: https://github.com/waltzofpearls/reckon/workflows/ci/badge.svg
[actions-url]: https://github.com/waltzofpearls/reckon/actions?query=workflow%3Aci+branch%3Amain
[mit-badge]: https://img.shields.io/badge/license-MIT-green.svg
[mit-url]: https://github.com/waltzofpearls/reckon/blob/main/LICENSE

Reckon is a prometheus exporter for time series forecasting and anomaly detection. It scrapes prometheus metrics
with a defined time range, train a predictive model with those metrics, and then expose the original and forecasted
metrics back through a prometheus HTTP endpoint.

An exmaple of original and forecasted metrics generated from [Prophet](https://facebook.github.io/prophet/) and
exposed from reckon:

```
# HELP sensehat_humidity_prophet Prophet forecasted metric value
# TYPE sensehat_humidity_prophet gauge
sensehat_humidity_prophet{column="original",instance="sensehat.rpi.topbass.studio:8000",job="sensehat_exporter"} 63.694435119628906
sensehat_humidity_prophet{column="yhat",instance="sensehat.rpi.topbass.studio:8000",job="sensehat_exporter"} 65.39311782093829
sensehat_humidity_prophet{column="yhat_lower",instance="sensehat.rpi.topbass.studio:8000",job="sensehat_exporter"} 64.17657873755101
sensehat_humidity_prophet{column="yhat_upper",instance="sensehat.rpi.topbass.studio:8000",job="sensehat_exporter"} 66.55323480537575
# HELP sensehat_temperature_prophet Prophet forecasted metric value
# TYPE sensehat_temperature_prophet gauge
sensehat_temperature_prophet{column="original",instance="sensehat.rpi.topbass.studio:8000",job="sensehat_exporter"} 28.072917938232422
sensehat_temperature_prophet{column="yhat",instance="sensehat.rpi.topbass.studio:8000",job="sensehat_exporter"} 27.972340899541923
sensehat_temperature_prophet{column="yhat_lower",instance="sensehat.rpi.topbass.studio:8000",job="sensehat_exporter"} 27.675004226891883
sensehat_temperature_prophet{column="yhat_upper",instance="sensehat.rpi.topbass.studio:8000",job="sensehat_exporter"} 28.28264431712968
```

## Try it

Gather the following info before start:

- Prometheus server address, for example, `http://prometheus.rpi.topbass.studio:9090`
- Metric names to watch and join them together with comma, for example, `sensehat_temperature,sensehat_humidity`

#### With Docker

This is the simplest method to get reckon running. You only need docker to get started.

```
PROMETHEUS={prometheus_server_address} WATCH_LIST={comma_separated_metric_names} make docker
```

#### With Go and Virtualenv

If you prefer running reckon without docker, or you would like to build and run the binary locally, you will need
Go, Python and Virtualenv. Make sure you have Go 1.16+ and Python 3.7

- Go: `brew install go` or [follow this gudie](https://golang.org/doc/install)
- Pyenv: `brew install pyenv` or [follow this guide](https://github.com/pyenv/pyenv#installation)
- Virtualenvwrapper: `pip install virtualenvwrapper` AND [follow this guide](https://virtualenvwrapper.readthedocs.io/en/latest/install.html)

```
PROMETHEUS={prometheus_server_address} WATCH_LIST={comma_separated_metric_names} make run
```
