package app.wheretopark.providers.cctv

import app.wheretopark.shared.*
import kotlinx.serialization.ExperimentalSerializationApi
import kotlinx.serialization.KSerializer
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.builtins.IntArraySerializer
import kotlinx.serialization.descriptors.SerialDescriptor
import kotlinx.serialization.encoding.Decoder
import kotlinx.serialization.encoding.Encoder

@Serializable
data class Point(val x: UInt, val y: UInt)

@Serializable
data class Configuration(
    @SerialName("parking-lots")
    val parkingLots: List<ParkingLot>
) {

    @Serializable
    data class ParkingSpot(

        val points: List<@Serializable(with = PointAsArraySerializer::class) Point>
    )

    @Serializable
    data class ParkingLot(
        val name: String,
        val address: String,
        val location: Coordinate,
        val resources: List<ParkingLotResource>,
        val features: List<ParkingLotFeature>,
        val comment: Map<LanguageCode, String> = mapOf(),
        val currency: String,
        val rules: List<ParkingLotRule>,

        @SerialName("camera-url")
        val cameraURL: String,
        val spots: List<ParkingSpot>,
    )
}

@OptIn(ExperimentalSerializationApi::class)
class PointAsArraySerializer : KSerializer<Point> {
    private val delegateSerializer = IntArraySerializer()
    override val descriptor = SerialDescriptor("Point", delegateSerializer.descriptor)

    override fun serialize(encoder: Encoder, value: Point) {
        val data = intArrayOf(value.x.toInt(), value.y.toInt())
        encoder.encodeSerializableValue(delegateSerializer, data)
    }

    override fun deserialize(decoder: Decoder): Point {
        val array = decoder.decodeSerializableValue(delegateSerializer)
        assert(array.count() == 2) { "Expected array of size 2, but got ${array.count()}" }
        return Point(array[0].toUInt(), array[1].toUInt())
    }
}