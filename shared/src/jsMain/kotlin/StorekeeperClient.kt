@file:OptIn(ExperimentalJsExport::class, DelicateCoroutinesApi::class)

package app.wheretopark.shared

import kotlinx.coroutines.DelicateCoroutinesApi
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.promise
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

    fun parkingLots(): Promise<Map<ParkingLotID, ParkingLot>> = GlobalScope.promise {
        client.parkingLots()
    }

    fun parkingLot(id: ParkingLotID): Promise<ParkingLot?> = GlobalScope.promise {
        client.parkingLot(id)
    }
}