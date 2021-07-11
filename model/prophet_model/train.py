import pandas as pd
import time
from prophet import Prophet

def train(data, duration):
    print('[PYTHON]', 'LENGTH:', len(data), 'DURATION:', pretty_timedelta(duration*60))

    metric_values = pd.DataFrame(data, columns=["ds", "y"]).apply(
        pd.to_numeric, errors="raise"
    )
    metric_values["ds"] = pd.to_datetime(metric_values["ds"], unit="s")

    model = Prophet(
        daily_seasonality=True, weekly_seasonality=True, yearly_seasonality=True
    )

    print('[PYTHON]', "BEGIN PROPHET TRAINING")
    begin = time.time()

    model.fit(metric_values)
    future = model.make_future_dataframe(
        periods=int(duration),
        freq="1MIN",
        include_history=False,
    )
    forecasted = model.predict(future)
    forecasted["timestamp"] = forecasted["ds"].map(lambda x: x.timestamp())
    forecasted = forecasted[["timestamp", "yhat", "yhat_lower", "yhat_upper"]]
    forecasted = forecasted.set_index("timestamp")

    end = time.time()
    print('[PYTHON]', 'TIME ELAPSED:', pretty_timedelta(end - begin))

    ''''
    forecasted example
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

    print('[PYTHON]', 'SENDING FORECASTED DATA:', len(forecasted), "ITEMS")
    return forecasted.to_dict('index')

def pretty_timedelta(seconds):
    seconds = int(seconds)
    days, seconds = divmod(seconds, 86400)
    hours, seconds = divmod(seconds, 3600)
    minutes, seconds = divmod(seconds, 60)
    if days > 0:
        return '{:d}d{:d}h{:d}m{:d}s'.format(days, hours, minutes, seconds)
    elif hours > 0:
        return '{:d}h{:d}m{:d}s'.format(hours, minutes, seconds)
    elif minutes > 0:
        return '{:d}m{:d}s'.format(minutes, seconds)
    else:
        return '{:d}s'.format(seconds)
