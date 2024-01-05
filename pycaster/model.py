import os
import pandas as pd
from prophet import Prophet
from prophet.serialize import model_to_json, model_from_json
from matplotlib import pyplot as plt  # Import pyplot

class ForecastingModel:
    def __init__(self, id: str, path: str, dataset: pd.DataFrame):
        if os.path.exists(path):
            with open(path, 'r') as fin:
                self.m = model_from_json(fin.read()) # Read model
                print(f'loaded model for {id} from {path}')
        else:
            self.m = Prophet()
            self.m.fit(dataset)
            with open(path, 'w') as fout:
                fout.write(model_to_json(self.m))  # Save model
                print(f'saved newly trained model for {id} to {path}')

    def predict(self, dates: pd.DatetimeIndex):
        future = pd.DataFrame(dates, columns=['ds'])
        fcst = self.m.predict(future)
        # print(fcst[['ds', 'yhat', 'yhat_lower', 'yhat_upper']].tail())
        # fig = self.m.plot(fcst)
        # plt.show()
        return fcst
