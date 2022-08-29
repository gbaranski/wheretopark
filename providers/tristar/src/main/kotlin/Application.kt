package app.wheretopark.providers.tristar

import app.wheretopark.providers.tristar.gdansk.TristarGdanskProvider
import app.wheretopark.providers.tristar.gdynia.TristarGdyniaProvider
import app.wheretopark.shared.*

abstract class Provider {
    abstract suspend fun metadatas(): Map<ParkingLotID, ParkingLotMetadata>
    abstract suspend fun states(): Map<ParkingLotID, ParkingLotState>
}


suspend fun f(storekeeperClient: StorekeeperClient, provider: Provider) {
    val metadatas = provider.metadatas()
    val states = provider.states()
    println(metadatas)
    println(states)
    storekeeperClient.postMetadatas(metadatas)
    storekeeperClient.postStates(states)
    println("Published ${metadatas.count()} metadatas and ${states.count()} states")
}

suspend fun main() {
    val storekeeperClient = StorekeeperClient()
    val gdanskProvider = TristarGdanskProvider()
    val gdyniaProvider = TristarGdyniaProvider()
    f(storekeeperClient, gdanskProvider)
    f(storekeeperClient, gdyniaProvider)
}