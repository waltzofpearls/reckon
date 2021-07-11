# reckon

Reckon is a Prometheus exporter for time series forecasting and anomaly detection. It scrapes Prometheus metrics
with a defined range of time, train a predictive model using those metrics, and expose those metrics back through
a Prometheus HTTP endpoint.
