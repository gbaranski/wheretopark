package app.wheretopark.storekeeper

import app.wheretopark.shared.ParkingLotID
import app.wheretopark.shared.ParkingLotMetadata
import app.wheretopark.shared.ParkingLotState
import io.github.crackthecodeabhi.kreds.connection.Endpoint
import io.github.crackthecodeabhi.kreds.connection.newClient
import kotlinx.datetime.Instant
import kotlinx.serialization.decodeFromString
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import kotlin.time.Duration.Companion.seconds

private const val METADATA_NAMESPACE = "metadata"
private const val STATE_NAMESPACE = "state"
private val STATE_DURATION = 3600.seconds

class RedisStore(
    host: String,
    port: Int?
): Store {
    private val client = newClient(Endpoint(host, port ?: 6379))
    private val serialization = Json


    private suspend inline fun <reified T> getAll(namespace: String): Map<ParkingLotID, T> {
        val ids = client.keys("${namespace}:*")
        if (ids.isEmpty()) {
            return mapOf()
        }
        val values = client.mget(*ids.toTypedArray())
        val pairs = ids.zip(values)
        return pairs.associate { (id, value) ->
            id.removePrefix("${namespace}:") to serialization.decodeFromString(value!!)
        }
    }

    private suspend inline fun <reified T> updateMany(namespace: String, values: Map<ParkingLotID, T>, getExpireAt: (T) -> Instant) {
        updateMany(namespace, values)
        values.map { (id, value) ->
            val expireAt = getExpireAt(value)
            client.expireAt("${namespace}:${id}", expireAt.toEpochMilliseconds().toULong())
        }
    }

    private suspend inline fun <reified T> updateMany(namespace: String, values: Map<ParkingLotID, T>) {
        val pairs = values.map {(id, value) ->
            "${namespace}:${id}" to serialization.encodeToString(value)
        }
        client.mset(*pairs.toTypedArray())
    }


    override suspend fun getMetadatas(): Map<ParkingLotID, ParkingLotMetadata> = getAll(METADATA_NAMESPACE)

    override suspend fun getStates(): Map<ParkingLotID, ParkingLotState> = getAll(STATE_NAMESPACE)

    override suspend fun updateMetadatas(metadata: Map<ParkingLotID, ParkingLotMetadata>) = updateMany(METADATA_NAMESPACE, metadata)

    override suspend fun updateStates(states: Map<ParkingLotID, ParkingLotState>) = updateMany(STATE_NAMESPACE, states, getExpireAt = {
        it.lastUpdated.plus(STATE_DURATION)
    })

}