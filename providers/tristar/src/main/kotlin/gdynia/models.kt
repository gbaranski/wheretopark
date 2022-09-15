package app.wheretopark.providers.tristar.gdynia

import app.wheretopark.shared.Coordinate
import kotlinx.datetime.*
import kotlinx.serialization.ExperimentalSerializationApi
import kotlinx.serialization.KSerializer
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.builtins.DoubleArraySerializer
import kotlinx.serialization.descriptors.PrimitiveKind
import kotlinx.serialization.descriptors.PrimitiveSerialDescriptor
import kotlinx.serialization.descriptors.SerialDescriptor
import kotlinx.serialization.encoding.Decoder
import kotlinx.serialization.encoding.Encoder

@Serializable
data class Location(
    val type: String,
    @Serializable(with = CoordinateAsGeoJSONSerializer::class) val coordinates: Coordinate,
)

object CoordinateAsGeoJSONSerializer : KSerializer<Coordinate> {
    private val delegateSerializer = DoubleArraySerializer()

    @OptIn(ExperimentalSerializationApi::class)
    override val descriptor = SerialDescriptor("Coordinate", delegateSerializer.descriptor)

    override fun serialize(encoder: Encoder, value: Coordinate) {
        val data = doubleArrayOf(
            value.longitude, value.latitude
        )
        encoder.encodeSerializableValue(delegateSerializer, data)
    }

    override fun deserialize(decoder: Decoder): Coordinate {
        val array = decoder.decodeSerializableValue(delegateSerializer)
        assert(array.count() == 2)
        return Coordinate(latitude = array[1], longitude = array[0])
    }
}

@Serializable
data class Metadata(
    val id: UInt,
    val code: String,
    val name: String,
    val address: String,
    val streetEntrance: String,
    // Did they really misspell pricing with "princing"?
    @SerialName("princing") val pricing: String,
    val location: Location,
    @Serializable(with = InstantAsRfc3339Serializer::class) val lastUpdate: Instant,
)

@Serializable
data class MetadataResponse(
    @SerialName("parkings") val parkingLots: List<Metadata>
)

@Serializable
data class State(
    val id: UInt,
    @SerialName("parkingId") val parkingID: UInt,
    val capacity: UInt,
    val freePlaces: UInt,
    @Serializable(with = InstantAsRfc3339Serializer::class) val insertTime: Instant
)

typealias StateResponse = List<State>

object InstantAsRfc3339Serializer : KSerializer<Instant> {
    override val descriptor: SerialDescriptor = PrimitiveSerialDescriptor("Instant", PrimitiveKind.STRING)

    override fun serialize(encoder: Encoder, value: Instant) {
        val string = value.toString().replace('T', ' ')
        encoder.encodeString(string)
    }

    override fun deserialize(decoder: Decoder): Instant {
        val string = decoder.decodeString().replace(' ', 'T').plus('Z')
        return Instant.parse(string).toLocalDateTime(TimeZone.UTC).toInstant(TimeZone.of("Europe/Warsaw"))
    }
}