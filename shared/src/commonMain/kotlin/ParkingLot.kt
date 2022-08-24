package app.wheretopark.shared

import io.ktor.http.*
import io.ktor.util.date.*
import kotlinx.datetime.Clock
import kotlinx.serialization.KSerializer
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.descriptors.PrimitiveKind
import kotlinx.serialization.descriptors.PrimitiveSerialDescriptor
import kotlinx.serialization.descriptors.SerialDescriptor
import kotlinx.serialization.encoding.Decoder
import kotlinx.serialization.encoding.Encoder
import kotlin.time.Duration
import kotlinx.datetime.DayOfWeek
import kotlinx.datetime.Instant
import kotlin.time.Duration.Companion.days
import kotlin.time.Duration.Companion.hours
import kotlin.time.Duration.Companion.seconds
import kotlin.time.ExperimentalTime

typealias ParkingLotID = String

@Serializable
data class ParkingLotLocation(
    val latitude: Double,
    val longitude: Double,
)

@Serializable
enum class ParkingLotFeature {
    UNCOVERED,
    COVERED,
    UNDERGROUND,
}

@Serializable
data class ParkingLotWeekdays(
    val start: DayOfWeek,
    val end: DayOfWeek
)

@Serializable
data class ParkingLotHours(
    val start: String,
    val end: String
)

@Serializable
data class ParkingLotPricingRule(
    @Serializable(with = DurationSerializer::class)
    val duration: Duration,
    val price: Double,
    val repeating: Boolean = false
)

object DurationSerializer : KSerializer<Duration> {
    override val descriptor = PrimitiveSerialDescriptor("Duration", PrimitiveKind.STRING)

    override fun deserialize(decoder: Decoder): Duration {
        return Duration.parse(decoder.decodeString())
    }

    override fun serialize(encoder: Encoder, value: Duration) {
        encoder.encodeString(value.toIsoString())
    }
}

@Serializable
data class ParkingLotRule(
    val weekdays: ParkingLotWeekdays? = null,
    val hours: ParkingLotHours? = null,
    val pricing: List<ParkingLotPricingRule>,
)

@Serializable
enum class ParkingLotStatus {
    OPENS_SOON,
    OPEN,
    CLOSES_SOON,
    CLOSED,
}

@Serializable
data class ParkingLotMetadata(
    val name: String,
    val address: String,
    val location: ParkingLotLocation,
    val emails: List<String>,
    @SerialName("phone-numbers")
    val phoneNumbers: List<String>,
    val websites: List<String>,
    @SerialName("total-spots")
    val totalSpots: Int,
    val features: List<ParkingLotFeature>,
    val currency: String,
    val rules: List<ParkingLotRule>,
) {
    fun status(at: Instant): ParkingLotStatus {
        return ParkingLotStatus.CLOSED
    }

    fun status(): ParkingLotStatus {
        return ParkingLotStatus.CLOSED
    }

}

@Serializable
data class ParkingLotState (
    @SerialName("last-updated")
    val lastUpdated: Instant,
    @SerialName("available-spots")
    val availableSpots: Int
)

data class ParkingLot(
    val metadata: ParkingLotMetadata,
    val state: ParkingLotState,
) {
    public companion object {
        public val galeriaBaltycka = ParkingLot(
            state = ParkingLotState(
                availableSpots = 10,
                lastUpdated = Clock.System.now().minus(10.seconds)
            ),
            metadata = ParkingLotMetadata(
            name = "Galeria Bałtycka",
            address = "ul. Dmowskiego",
            location = ParkingLotLocation(
                latitude = 54.38268,
                longitude = 18.60024,
            ),
            emails = listOf("galeria@galeriabaltycka.pl"),
            phoneNumbers = listOf("+48 58 521 85 52"),
            websites = listOf("https://www.galeriabaltycka.pl/o-centrum/dojazd-parkingi/parkingi/"),
            totalSpots = 1100,
            features = listOf(ParkingLotFeature.COVERED, ParkingLotFeature.UNCOVERED),
            currency = "PLN",
            rules = listOf(
                ParkingLotRule(
                    weekdays = ParkingLotWeekdays(start = DayOfWeek.MONDAY, end = DayOfWeek.SATURDAY),
                    hours = ParkingLotHours(start = "08:00", end = "22:00"),
                    pricing = listOf(
                        ParkingLotPricingRule(
                            duration = 1.hours,
                            price = 0.0,
                        ),
                        ParkingLotPricingRule(
                            duration = 2.hours,
                            price = 2.0,
                        ),
                        ParkingLotPricingRule(
                            duration = 3.hours,
                            price = 5.0,
                        ),
                        ParkingLotPricingRule(
                            duration = 1.days,
                            price = 25.0,
                        ),
                        ParkingLotPricingRule(
                            duration = 1.hours,
                            price = 4.0,
                            repeating = true,
                        ),
                    )
                ),
                ParkingLotRule(
                    weekdays = ParkingLotWeekdays(start = DayOfWeek.SUNDAY, end = DayOfWeek.SUNDAY),
                    hours = ParkingLotHours(start = "09:00", end = "21:00"),
                    pricing = listOf(
                        ParkingLotPricingRule(
                            duration = 1.hours,
                            price = 0.0,
                        ),
                        ParkingLotPricingRule(
                            duration = 2.hours,
                            price = 2.0,
                        ),
                        ParkingLotPricingRule(
                            duration = 3.hours,
                            price = 5.0,
                        ),
                        ParkingLotPricingRule(
                            duration = 1.days,
                            price = 25.0,
                        ),
                        ParkingLotPricingRule(
                            duration = 1.hours,
                            price = 4.0,
                            repeating = true,
                        ),
                    )
                )
            )
        ))
    }
}
