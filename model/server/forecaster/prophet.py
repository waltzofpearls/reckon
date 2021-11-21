from prophet import Prophet
import model.api.forecast_pb2 as serializer
import model.api.forecast_pb2_grpc as grpc_lib
import pandas as pd
import time

class Forecaster(grpc_lib.ForecastServicer):
    def Prophet(self, request, context):
        self.logger.info('data received in Prophet model', extra={
            'length': len(request.values),
            'duration': self.pretty_timedelta(request.duration*60)
        })

        data = []
        for value in request.values:
            data.append([value.timestamp, value.value])
        metric_values = pd.DataFrame(data, columns=["ds", "y"]).apply(
            pd.to_numeric, errors="raise"
        )
        metric_values["ds"] = pd.to_datetime(metric_values["ds"], unit="s")

        model = Prophet(
            daily_seasonality=True, weekly_seasonality=True, yearly_seasonality=True
        )

        self.logger.info('begin Prophet model training')
        begin = time.time()

        model.fit(metric_values)
        future = model.make_future_dataframe(
            periods=int(request.duration),
            freq="1MIN",
            include_history=False
        )
        forecasted = model.predict(future)
        forecasted["timestamp"] = forecasted["ds"].map(lambda x: x.timestamp())
        forecasted = forecasted[["timestamp", "yhat", "yhat_lower", "yhat_upper"]]
        forecasted = forecasted.set_index("timestamp")

        end = time.time()
        self.logger.info('Prophet model training completed', extra={
            'time_elapsed': self.pretty_timedelta(end - begin),
            'forecasted_items': len(forecasted)
        })

        ''''
        here is an example of the dict returned from forecasted.to_dict('index')
        {
            1626066880.0: {'yhat': 54.25773692364051, 'yhat_lower': -185.86933939411077, 'yhat_upper': 292.58353401454053},
            1626066940.0: {'yhat': 54.268292717184536, 'yhat_lower': -185.78517860185912, 'yhat_upper': 292.3456783844498},
            1626067000.0: {'yhat': 54.278932934537536, 'yhat_lower': -186.22052382613808, 'yhat_upper': 292.784502130985},
            1626067060.0: {'yhat': 54.28965884353986, 'yhat_lower': -186.0651051496073, 'yhat_upper': 292.7570446167972},
            1626067120.0: {'yhat': 54.300471710615334, 'yhat_lower': -186.18541588535416, 'yhat_upper': 293.01597966352915},
            1626067180.0: {'yhat': 54.31137280015447, 'yhat_lower': -186.2887686717037, 'yhat_upper': 293.3697261469165},
            1626067240.0: {'yhat': 54.32236337443156, 'yhat_lower': -186.3326835298265, 'yhat_upper': 292.97249639455214},
            1626067300.0: {'yhat': 54.33344469267559, 'yhat_lower': -186.68846226172354, 'yhat_upper': 293.0033622182267},
        }
        '''

        values = []
        for timestamp, forecast in forecasted.to_dict('index').items():
            values.append(serializer.Forecasted(
                timestamp=timestamp,
                yhat=forecast['yhat'],
                yhatLower=forecast['yhat_lower'],
                yhatUpper=forecast['yhat_upper']
            ))
        return serializer.ProphetReply(values=values)
