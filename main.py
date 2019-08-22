from model.data import Data
from model.model import Model


def main():
    data = Data()
    df = data.fetch()
    model = Model("ARIMA")
    model.train_with(df)


if __name__ == "__main__":
    main()
