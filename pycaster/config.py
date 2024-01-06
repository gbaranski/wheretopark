import os

pycaster_path = os.environ.get('PYCASTER_DATA', os.path.join(os.path.expanduser("~"), ".local/share/wheretopark/pycaster"))
models_path = os.path.join(pycaster_path, "models")