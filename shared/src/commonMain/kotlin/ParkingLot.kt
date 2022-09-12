@file:OptIn(ExperimentalJsExport::class)

package app.wheretopark.shared

import kotlinx.datetime.*
import kotlinx.serialization.KSerializer
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.descriptors.PrimitiveKind
import kotlinx.serialization.descriptors.PrimitiveSerialDescriptor
import kotlinx.serialization.descriptors.SerialDescriptor
import kotlinx.serialization.encoding.Decoder
import kotlinx.serialization.encoding.Encoder
import kotlin.js.ExperimentalJsExport
import kotlin.js.JsExport
import kotlin.time.Duration
import kotlin.time.Duration.Companion.days
import kotlin.time.Duration.Companion.hours
import kotlin.time.Duration.Companion.minutes
import kotlin.time.Duration.Companion.seconds

typealias ParkingLotID = String

@Serializable
@JsExport
enum class ParkingLotFeature {
    UNCOVERED,
    COVERED,
    UNDERGROUND,
}


@Serializable(with = ParkingLotWeekdaysSerializer::class)
@JsExport
data class ParkingLotWeekdays(
    val start: DayOfWeek,
    val end: DayOfWeek
)

object ParkingLotWeekdaysSerializer : KSerializer<ParkingLotWeekdays> {
    override val descriptor = PrimitiveSerialDescriptor("ParkingLotWeekdays", PrimitiveKind.STRING)
    private const val delimiter = '-'
    override fun serialize(encoder: Encoder, value: ParkingLotWeekdays) {
        val string = "${value.start}${delimiter}${value.end}"
        encoder.encodeString(string)
    }

    override fun deserialize(decoder: Decoder): ParkingLotWeekdays {
        val string = decoder.decodeString()
        val split = string.split(delimiter, limit = 2)
        val start = DayOfWeek.valueOf(split[0])
        val end = DayOfWeek.valueOf(split[1])
        return ParkingLotWeekdays(start, end)
    }
}

@Serializable(with = ParkingLotHoursSerializer::class)
@JsExport
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
        val split = string.split('-', limit = 2)
        val start = LocalTime.parse(split[0])
        val end = LocalTime.parse(split[1])
        return ParkingLotHours(start, end)
    }
}


@Serializable
@JsExport
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
@JsExport
data class ParkingLotRule(
    val weekdays: ParkingLotWeekdays? = null,
    val hours: ParkingLotHours? = null,
    val pricing: List<ParkingLotPricingRule>,
)

@Serializable
@JsExport
enum class ParkingLotStatus {
    OPENS_SOON,
    OPEN,
    CLOSES_SOON,
    CLOSED,
}

@Serializable(with = ParkingLotResourceSerializer::class)
@JsExport
data class ParkingLotResource(val url: String) : Comparable<ParkingLotResource> {
    fun label() =
        when (url.substringBefore(':')) {
            "http", "https" -> "Website"
            "mailto" -> "E-Mail address"
            "tel" -> "Phone number"
            else -> "Unknown"
        }

    override fun toString() = url

    override fun compareTo(other: ParkingLotResource): Int {
        return url.compareTo(other.url)
    }
}

object ParkingLotResourceSerializer : KSerializer<ParkingLotResource> {
    override val descriptor: SerialDescriptor = PrimitiveSerialDescriptor("ParkingLotResource", PrimitiveKind.STRING)
    override fun serialize(encoder: Encoder, value: ParkingLotResource) = encoder.encodeString(value.url)
    override fun deserialize(decoder: Decoder) = ParkingLotResource(decoder.decodeString())
}


@JsExport
enum class ParkingSpotType {
    CAR,
    MOTORCYCLE,
    HANDICAP,
    ELECTRIC,
}

@Serializable
@JsExport
data class ParkingLotMetadata(
    val name: String,
    val address: String,
    val location: Coordinate,
    val resources: List<ParkingLotResource>,
    @SerialName("total-spots")
    val totalSpots: Map<ParkingSpotType, UInt>,
    val features: List<ParkingLotFeature>,
    val currency: String,
    val rules: List<ParkingLotRule>,
) {
    fun statusAt(at: Instant): ParkingLotStatus {
        val dateTime = at.toLocalDateTime(TimeZone.UTC)
        val weekday = dateTime.dayOfWeek
        val rule = rules.sortedBy { it.weekdays != null }.find {
            weekday >= (it.weekdays?.start ?: DayOfWeek.MONDAY) && weekday <= (it.weekdays?.end ?: DayOfWeek.SUNDAY)
        } ?: return ParkingLotStatus.CLOSED
        if (rule.hours == null) return ParkingLotStatus.OPEN
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

    fun status(): ParkingLotStatus = statusAt(Clock.System.now())
}

@Serializable
@JsExport
data class ParkingLotState(
    @SerialName("last-updated")
    val lastUpdated: Instant,
    @SerialName("available-spots")
    val availableSpots: Map<ParkingSpotType, UInt>,
)

fun Map<ParkingLotID, ParkingLot>.split(): Pair<Map<ParkingLotID, ParkingLotMetadata>, Map<ParkingLotID, ParkingLotState>> {
    val metadatas = entries.associate { it.key to it.value.metadata }
    val states = entries.associate { it.key to it.value.state }
    return Pair(metadatas, states)
}

@Serializable
@JsExport
data class ParkingLot(
    val metadata: ParkingLotMetadata,
    val state: ParkingLotState,
) {
    companion object {
        val forumGdansk = ParkingLot(
            state = ParkingLotState(
                lastUpdated = Instant.parse("2022-09-06T21:39:21Z"),
                availableSpots = mapOf(
                    ParkingSpotType.CAR to 988u,
                )
            ),
            metadata = ParkingLotMetadata(
                name = "Forum Gdańsk",
                address = "Targ Sienny 7, 80-806 Gdańsk",
                location = Coordinate(54.34941, 18.64233),
                resources = listOf(
                    ParkingLotResource("mailto:parking@forumgdansk.pl"),
                    ParkingLotResource("tel:+48-661-551-882"),
                    ParkingLotResource("https://forumgdansk.pl/pl/przydatne-informacje/parking")
                ),
                totalSpots = mapOf(
                    ParkingSpotType.CAR to 1008u,
                ),
                features = listOf(
                    ParkingLotFeature.COVERED,
                    ParkingLotFeature.UNCOVERED,
                ),
                currency = "PLN",
                rules = listOf(
                    ParkingLotRule(
                        weekdays = null,
                        hours = null,
                        pricing = listOf(
                            ParkingLotPricingRule(
                                duration = 30.minutes,
                                price = 0.0,
                            ),
                            ParkingLotPricingRule(
                                duration = 1.hours,
                                price = 3.0,
                            ),
                            ParkingLotPricingRule(
                                duration = 2.hours,
                                price = 8.0,
                            ),
                            ParkingLotPricingRule(
                                duration = 3.hours,
                                price = 13.0,
                            ),
                            ParkingLotPricingRule(
                                duration = 1.hours,
                                price = 10.0,
                                repeating = true,
                            ),
                        )
                    )
                )
            )

        )

        val galeriaBaltycka = ParkingLot(
            state = ParkingLotState(
                availableSpots = mapOf(
                    Pair(ParkingSpotType.CAR, 10u)
                ),
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
                totalSpots = mapOf(
                    Pair(ParkingSpotType.CAR, 1100u)
                ),
                features = listOf(ParkingLotFeature.COVERED, ParkingLotFeature.UNCOVERED),
                currency = "PLN",
                rules = listOf(
                    ParkingLotRule(
                        weekdays = ParkingLotWeekdays(start = DayOfWeek.MONDAY, end = DayOfWeek.SATURDAY),
                        hours = ParkingLotHours(start = LocalTime(8, 0, 0), end = LocalTime(22, 0, 0)),
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
                        hours = ParkingLotHours(start = LocalTime(9, 0, 0), end = LocalTime(21, 0, 0)),
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
            )
        )
    }
}
