package app.wheretopark.providers.tristar.gdansk

import app.wheretopark.shared.Coordinate
import kotlinx.datetime.Instant
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class Metadata(
    val id: String,
    val name: String,
    val shortName: String,
    val address: String,
    val streetEntrance: String,
    val location: Coordinate,
)

@Serializable
data class State(
    @SerialName("parkingId")
    val id: String,
    val availableSpots: UInt,
    val lastUpdate: Instant
)


@Serializable
data class Response<P>(
    val lastUpdate: Instant,
    val parkingLots: List<P>,
)

typealias StateResponse = Response<State>
typealias MetadataResponse = Response<Metadata>
