import pandas as pd


def load_file(path: str):
    dataset = pd.read_csv(path)
    dataset.rename(columns={'occupancy': 'y', 'date': 'ds'}, inplace=True)
    # # Convert 'date' column to datetime
    dataset['ds'] = pd.to_datetime(dataset['ds'])

    # Setting the 'date' as the index
    # dataset.set_index('ds', inplace=True)

    # Check for any missing values
    missing_values = dataset.isnull().sum().sum()

    # raise an error if there are missing values
    if missing_values > 0:
        raise ValueError("Missing values in the dataset")

    return dataset
