import importlib


class Model:
    ADAPTERS = {
        "ARIMA": "model.adapters.arima"
    }

    def __init__(self, adapter="ARIMA"):
        if adapter not in self.ADAPTERS:
            raise ValueError("unknown model adapter: " + adapter)

        module = importlib.import_module(self.ADAPTERS[adapter])
        class_ = getattr(module, adapter)
        self.adapter = class_()

    def train_with(self, data):
        self.adapter.train_with(data)
