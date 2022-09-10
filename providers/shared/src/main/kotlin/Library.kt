package app.wheretopark.providers.shared

import app.wheretopark.shared.*
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import kotlin.time.Duration

abstract class Provider {
    abstract val name: String
    abstract val metadataInterval: Duration
    abstract val stateInterval: Duration
    abstract suspend fun metadatas(): Map<ParkingLotID, ParkingLotMetadata>
    abstract suspend fun states(): Map<ParkingLotID, ParkingLotState>

    private suspend fun metadata(storekeeperClient: StorekeeperClient) {
        val metadatas = metadatas()
        storekeeperClient.postMetadatas(metadatas)
        println("$name published ${metadatas.count()} metadatas")
    }

    private suspend fun state(storekeeperClient: StorekeeperClient) {
        val states = states()
        storekeeperClient.postStates(states)
        println("$name published ${states.count()} states")
    }

    suspend fun start(storekeeperClient: StorekeeperClient) = coroutineScope {
        metadata(storekeeperClient)
        launch {
            delay(metadataInterval)
            runEvery(metadataInterval) { metadata(storekeeperClient) }
        }
        launch {
            runEvery(stateInterval) { state(storekeeperClient) }
        }
    }

//    suspend fun start() = start(StorekeeperClient(getStorekeeperURL()))
}


private fun getStorekeeperURL(): String {
    val url = System.getenv("STOREKEEPER_URL")
    println("using `$url` as Storekeeper service URL")
    return url
}

private fun getAuthorizationURL(): String {
    val url = System.getenv("AUTHORIZATION_URL")
    println("using `$url` as Authorization service URL")
    return url
}

private suspend fun runEvery(delay: Duration, action: suspend () -> Unit) {
    while (true) {
        action()
        delay(delay)
    }
}

suspend fun startMany(vararg providers: Provider) = coroutineScope {
    val clientID = System.getenv("CLIENT_ID")!!
    val clientSecret = System.getenv("CLIENT_SECRET")!!
    val authorizationClient = AuthorizationClient(getAuthorizationURL(), clientID, clientSecret)
    val storekeeperClient = StorekeeperClient(
        getStorekeeperURL(),
        authorizationClient,
        setOf(AccessType.WriteMetadata, AccessType.WriteState)
    )
    providers.forEach { launch { it.start(storekeeperClient) } }
}