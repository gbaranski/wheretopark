package app.wheretopark.providers.tristar

import app.wheretopark.providers.shared.startMany
import app.wheretopark.providers.tristar.gdansk.TristarGdanskProvider
import app.wheretopark.providers.tristar.gdynia.TristarGdyniaProvider
suspend fun main() = startMany(TristarGdanskProvider(), TristarGdyniaProvider())