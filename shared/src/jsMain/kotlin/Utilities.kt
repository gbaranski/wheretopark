@file:OptIn(ExperimentalJsExport::class)

package app.wheretopark.shared

import kotlinx.coroutines.DelicateCoroutinesApi
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.promise
import kotlinx.datetime.Instant
import kotlinx.datetime.toDateTimePeriod
import kotlinx.datetime.toJSDate
import kotlinx.js.Record
import kotlinx.js.set
import kotlinx.serialization.decodeFromString
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import kotlin.js.Date
import kotlin.js.Promise
import kotlin.time.Duration
import kotlin.time.Duration.Companion.seconds

@JsExport
fun <K : Any, V : Any> Map<K, V>.toRecord(): Record<K, V> {
    val record = Record<K, V>()
    this.forEach { (key, value) -> record[key] = value }
    return record
}

@JsExport
fun <T: Any> List<T>.toArray(): Array<T> = this.toTypedArray()


@JsExport
fun parseParkingLot(json: String) = Json.decodeFromString<ParkingLot>(json)

@JsExport
fun parseParkingLots(json: String) = Json.decodeFromString<Map<ParkingLotID, ParkingLot>>(json).toRecord()

@JsExport
fun instantToJSDate(from: Instant) = from.toJSDate()

@JsExport
fun durationToISO(duration: Duration) = duration.toIsoString()
