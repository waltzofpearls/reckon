# reckon

Reckon is a Prometheus exporter for time series forecasting and anomaly detection. It scrapes Prometheus metrics
with a defined range of time, train a predictive model using those metrics, and expose those metrics back through
a Prometheus HTTP endpoint.

By default, ARIMA is used as reckon's predictive model. To support other forecasting models and libraries, new
adapters need to be created.
