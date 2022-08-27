package app.wheretopark.shared

import io.ktor.http.*
import kotlinx.datetime.*
import kotlinx.serialization.KSerializer
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.descriptors.PrimitiveKind
import kotlinx.serialization.descriptors.PrimitiveSerialDescriptor
import kotlinx.serialization.descriptors.SerialDescriptor
import kotlinx.serialization.encoding.Decoder
import kotlinx.serialization.encoding.Encoder
import kotlin.time.Duration
import kotlin.time.Duration.Companion.days
import kotlin.time.Duration.Companion.hours
import kotlin.time.Duration.Companion.seconds

typealias ParkingLotID = String

@Serializable
enum class ParkingLotFeature {
    UNCOVERED,
    COVERED,
    UNDERGROUND,
}


@Serializable(with = ParkingLotWeekdaysSerializer::class)
data class ParkingLotWeekdays(
    val start: DayOfWeek,
    val end: DayOfWeek
)

object ParkingLotWeekdaysSerializer : KSerializer<ParkingLotWeekdays> {
    override val descriptor = PrimitiveSerialDescriptor("ParkingLotWeekdays", PrimitiveKind.STRING)
    override fun serialize(encoder: Encoder, value: ParkingLotWeekdays) {
        val string = "${value.start}-${value.end}"
        encoder.encodeString(string)
    }
    override fun deserialize(decoder: Decoder): ParkingLotWeekdays {
        val string = decoder.decodeString()
        val split = string.split('-', limit = 1)
        val start = DayOfWeek.valueOf(split[0].uppercase())
        val end = DayOfWeek.valueOf(split[1].uppercase())
        return ParkingLotWeekdays(start, end)
    }
}

@Serializable(with = ParkingLotHoursSerializer::class)
data class ParkingLotHours(
    val start: LocalTime,
    val end: LocalTime
)

object ParkingLotHoursSerializer : KSerializer<ParkingLotHours> {
    override val descriptor = PrimitiveSerialDescriptor("ParkingLotHours", PrimitiveKind.STRING)
    override fun serialize(encoder: Encoder, value: ParkingLotHours) {
        val string = "${value.start}-${value.end}"
        encoder.encodeString(string)
    }
    override fun deserialize(decoder: Decoder): ParkingLotHours {
        val string = decoder.decodeString()
        val split = string.split('-', limit = 1)
        val start = LocalTime.parse(split[0])
        val end = LocalTime.parse(split[1])
        return ParkingLotHours(start, end)
    }
}


@Serializable
data class ParkingLotPricingRule(
    @Serializable(with = DurationSerializer::class)
    val duration: Duration,
    val price: Double,
    val repeating: Boolean = false
)

object DurationSerializer : KSerializer<Duration> {
    override val descriptor = PrimitiveSerialDescriptor("Duration", PrimitiveKind.STRING)
    override fun serialize(encoder: Encoder, value: Duration) = encoder.encodeString(value.toIsoString())
    override fun deserialize(decoder: Decoder) = Duration.parse(decoder.decodeString())
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

@Serializable(with = ParkingLotResourceSerializer::class)
data class ParkingLotResource(val url: Url): Comparable<ParkingLotResource> {
    constructor(s: String) : this(Url(s))

    fun label() =
        when(url.protocol.name) {
            "http", "https" -> "Website"
            "mailto" -> "E-Mail address"
            "tel" -> "Phone number"
            else -> "Unknown"
        }

    override fun toString() = url.toString()
    override fun compareTo(other: ParkingLotResource): Int {
        return (url == other.url).compareTo(false)
    }
}

object ParkingLotResourceSerializer : KSerializer<ParkingLotResource> {
    override val descriptor: SerialDescriptor = PrimitiveSerialDescriptor("ParkingLotResource", PrimitiveKind.STRING)

    override fun serialize(encoder: Encoder, value: ParkingLotResource) = UrlSerializer.serialize(encoder, value.url)

    override fun deserialize(decoder: Decoder) = ParkingLotResource(UrlSerializer.deserialize(decoder))
}


@Serializable
data class ParkingLotMetadata(
    val name: String,
    val address: String,
    val location: Coordinate,
    val resources: List<ParkingLotResource>,
    val totalSpots: UInt,
    val features: List<ParkingLotFeature>,
    val currency: String,
    val rules: List<ParkingLotRule>,
) {
    fun status(at: Instant): ParkingLotStatus {
        val dateTime = at.toLocalDateTime(TimeZone.UTC)
        val weekday = dateTime.dayOfWeek
        val rule = rules.sortedBy { it.weekdays != null }.find{
            weekday >= (it.weekdays?.start ?: DayOfWeek.MONDAY) && weekday <= (it.weekdays?.end ?: DayOfWeek.SUNDAY)
        } ?: return ParkingLotStatus.CLOSED
        if(rule.hours == null) return ParkingLotStatus.OPEN
        val startDateTime = rule.hours.start.atDate(dateTime.date)
        val endDateTime = if (rule.hours.end < rule.hours.start) {
            rule.hours.end.atDate(dateTime.date.plus(DatePeriod(days = 1)))
        } else {
            rule.hours.end.atDate(dateTime.date)
        }
        val toStart = dateTime.toInstant(TimeZone.UTC).periodUntil(startDateTime.toInstant(TimeZone.UTC), TimeZone.UTC)
        val toEnd = dateTime.toInstant(TimeZone.UTC).periodUntil(endDateTime.toInstant(TimeZone.UTC), TimeZone.UTC)
        val isOpen = dateTime >= startDateTime && dateTime < endDateTime
        return if (isOpen) {
            if (toEnd.hours == 0) ParkingLotStatus.CLOSES_SOON else ParkingLotStatus.OPEN
        } else if (toStart.hours == 0) {
            ParkingLotStatus.OPENS_SOON
        } else {
            ParkingLotStatus.CLOSED
        }
    }

    fun status(): ParkingLotStatus = status(Clock.System.now())
}

@Serializable
data class ParkingLotState (
    @SerialName("last-updated")
    val lastUpdated: Instant,
    @SerialName("available-spots")
    val availableSpots: UInt
)

data class ParkingLot(
    val metadata: ParkingLotMetadata,
    val state: ParkingLotState,
) {
    companion object {
        val galeriaBaltycka = ParkingLot(
            state = ParkingLotState(
                availableSpots = 10u,
                lastUpdated = Clock.System.now().minus(10.seconds)
            ),
            metadata = ParkingLotMetadata(
            name = "Galeria Bałtycka",
            address = "ul. Dmowskiego",
            location = Coordinate(
                latitude = 54.38268,
                longitude = 18.60024,
            ),
            resources = listOf(
                ParkingLotResource("mailto:galeria@galeriabaltycka.pl"),
                ParkingLotResource("tel:+48-58-521-85-52"),
                ParkingLotResource("https://www.galeriabaltycka.pl/o-centrum/dojazd-parkingi/parkingi/")
            ),
            totalSpots = 1100u,
            features = listOf(ParkingLotFeature.COVERED, ParkingLotFeature.UNCOVERED),
            currency = "PLN",
            rules = listOf(
                ParkingLotRule(
                    weekdays = ParkingLotWeekdays(start = DayOfWeek.MONDAY, end = DayOfWeek.SATURDAY),
                    hours = ParkingLotHours(start = LocalTime(8,0,0), end = LocalTime(22,0,0)),
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
                    hours = ParkingLotHours(start = LocalTime(9,0,0), end = LocalTime(21,0,0)),
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
