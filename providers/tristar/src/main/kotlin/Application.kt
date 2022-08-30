package app.wheretopark.providers.tristar

import app.wheretopark.providers.tristar.gdansk.TristarGdanskProvider
import app.wheretopark.providers.tristar.gdynia.TristarGdyniaProvider
import app.wheretopark.shared.*
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.launch
import kotlinx.coroutines.time.delay
import java.time.Duration

abstract class Provider {
    abstract val name: String
    abstract suspend fun metadatas(): Map<ParkingLotID, ParkingLotMetadata>
    abstract suspend fun states(): Map<ParkingLotID, ParkingLotState>
}


val METADATA_INTERVAL = Duration.ofMinutes(60)!!
val STATE_INTERVAL = Duration.ofSeconds(30)!!

suspend fun metadata(storekeeperClient: StorekeeperClient, provider: Provider) {
    val metadatas = provider.metadatas()
    storekeeperClient.postMetadatas(metadatas)
    println("tristar/${provider.name} published ${metadatas.count()} metadatas")
}

suspend fun state(storekeeperClient: StorekeeperClient, provider: Provider) {
    val states = provider.states()
    storekeeperClient.postStates(states)
    println("tristar/${provider.name} published ${states.count()} states")
}

suspend fun runEvery(delay: Duration, action: suspend () -> Unit) {
    while(true) {
        action()
        delay(delay)
    }
}

suspend fun run(storekeeperClient: StorekeeperClient, provider: Provider) = coroutineScope {
    metadata(storekeeperClient, provider)
    launch {
        delay(METADATA_INTERVAL)
        runEvery(METADATA_INTERVAL) { metadata(storekeeperClient, provider) }
    }
    launch {
        runEvery(STATE_INTERVAL) { state(storekeeperClient, provider) }
    }
}

suspend fun main(): Unit = coroutineScope {
    val storekeeperClient = StorekeeperClient()
    val gdanskProvider = TristarGdanskProvider()
    val gdyniaProvider = TristarGdyniaProvider()
    launch {
        run(storekeeperClient, gdanskProvider)
    }
    launch {
        run(storekeeperClient, gdyniaProvider)
    }
}