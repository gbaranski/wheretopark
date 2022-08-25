package app.wheretopark.providers.tristar

import app.wheretopark.providers.tristar.gdansk.TristarGdanskProvider
import app.wheretopark.shared.*
import io.ktor.client.*
import io.ktor.client.engine.cio.*

abstract class Provider {
    abstract suspend fun metadatas(): Map<ParkingLotID, ParkingLotMetadata>
    abstract suspend fun states(): Map<ParkingLotID, ParkingLotState>
}


suspend fun main() {
    val storekeeperClient = StorekeeperClient()
    val gdanskProvider = TristarGdanskProvider()
    val metadatas = gdanskProvider.metadatas()
    val states = gdanskProvider.states()
    println(metadatas)
    println(states)
    storekeeperClient.postMetadatas(metadatas)
    storekeeperClient.postStates(states)
}