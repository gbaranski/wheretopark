package app.wheretopark.providers.tristar.gdynia

import kotlinx.datetime.Instant
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class Location(val a: Int)

@Serializable
data class Metadata(
    val id: UInt,
    val code: String,
    val name: String,
    val address: String,
    val streetEntrance: String,
    // Did they really misspell pricing with "princing"?
    @SerialName("princing")
    val pricing: String,
    val location: Location,
    val last_update: Instant,
)

@Serializable
data class MetadataResponse(
    @SerialName("parkings")
    val parkingLots: List<Metadata>
)

@Serializable
data class State(
    val id: UInt,
    val parkingID: UInt,
    val capacity: UInt,
    val freePlaces: UInt,
    val insertTime: Instant
)

typealias StateResponse = List<State>