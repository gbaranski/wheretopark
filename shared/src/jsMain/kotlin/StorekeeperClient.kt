@file:OptIn(ExperimentalJsExport::class, DelicateCoroutinesApi::class)

package app.wheretopark.shared

import kotlinx.coroutines.DelicateCoroutinesApi
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.promise
import kotlinx.js.Record
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import kotlin.js.Promise

@JsExport
class JSStorekeeperClient(
    url: String,
    scope: Array<AccessType>,
    authorizationClient: JSAuthorizationClient
) {
    private val client = StorekeeperClient(url, authorizationClient.client, scope.toSet())

    fun parkingLots(): Promise<String> = GlobalScope.promise {
        val parkingLots = client.parkingLots()
        Json.encodeToString(parkingLots)
    }

    fun parkingLot(id: ParkingLotID): Promise<String?> = GlobalScope.promise {
        val parkingLot = client.parkingLot(id)
        parkingLot?.let { Json.encodeToString(it) }
    }
}