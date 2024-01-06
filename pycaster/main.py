import config
import os
from model import ForecastingModel
from fastapi import FastAPI, File, UploadFile, Response
from pydantic import BaseModel
from io import StringIO
import pandas as pd
from datetime import datetime

app = FastAPI()

models: dict[str, ForecastingModel] = dict()

@app.post("/forecast/{parking_lot_id}")
async def forecast(parking_lot_id: str, start: datetime | None, end: datetime | None, file: UploadFile = File(...)):
    # Read the contents of the uploaded file
    content = await file.read()

    # Use StringIO to convert bytes to a file-like string object
    string_io = StringIO(content.decode("utf-8"))

    timeseries = pd.read_csv(string_io)
    dataset = timeseries.rename(columns={'date': 'ds', 'occupancy': 'y'})
    dataset['ds'] = pd.to_datetime(dataset['ds'])
    print(dataset)
    if parking_lot_id not in models:
        model = ForecastingModel(
            id=parking_lot_id,
            dataset=dataset,
            path=os.path.join(config.models_path, parking_lot_id + ".json")
        )
        models[parking_lot_id] = model

    model = models[parking_lot_id]

    if start is None:
        start = datetime.now().replace(hour=10, minute=0, second=0, microsecond=0, tzinfo=None)
    else:
        start = start.replace(tzinfo=None)
        
    if end is None:
        end = datetime.now().replace(hour=20, minute=0, second=0, microsecond=0, tzinfo=None)
    else:
        end = end.replace(tzinfo=None)
    
    dates = pd.date_range(start=start, end=end, freq='15min')
    predictions = model.predict(dates)

    predictions = pd.DataFrame(predictions[['ds', 'yhat']])
    predictions = predictions.rename(columns={'ds': 'date', 'yhat': 'occupancy'})

    stream = StringIO()
    predictions.to_csv(stream, index=False)
    response = Response(stream.getvalue(), media_type="text/csv")

    return response
