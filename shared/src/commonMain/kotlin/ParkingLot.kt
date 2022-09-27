@file:OptIn(ExperimentalJsExport::class)
@file:Suppress("NON_EXPORTABLE_TYPE")

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
    GUARDED,
}


@Serializable
@JsExport
enum class PaymentMethod {
    CASH,
    CARD,
    CONTACTLESS,
    MOBILE,
}

fun DayOfWeek.human() = name.lowercase().replaceFirstChar { it.uppercase() }

@Serializable(with = ParkingLotWeekdaysSerializer::class)
@JsExport
data class ParkingLotWeekdays(
    val start: DayOfWeek,
    val end: DayOfWeek
) {
    fun human() = "${start.human()}-${end.human()}"
}

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
) {
    fun human() = "${start}-${end}"
}

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


fun Duration.human() = this.toString()

@Serializable
@JsExport
data class ParkingLotPricingRule(
    @Serializable(with = DurationSerializer::class)
    val duration: Duration,
    val price: Double,
    val repeating: Boolean = false
) {
    fun durationComponents() = duration.components()
}

object DurationSerializer : KSerializer<Duration> {
    override val descriptor = PrimitiveSerialDescriptor("Duration", PrimitiveKind.STRING)
    override fun serialize(encoder: Encoder, value: Duration) = encoder.encodeString(value.toIsoString())
    override fun deserialize(decoder: Decoder) = Duration.parseIsoString(decoder.decodeString())
}

@Serializable
@JsExport
data class ParkingLotRule(
    // https://schema.org/openingHours
    // https://wiki.openstreetmap.org/wiki/Key:opening_hours
    val hours: String,
    // If not empty, then include only those from this list
    val includes: Set<ParkingSpotType> = emptySet(),
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
    CAR_DISABLED,
    CAR_ELECTRIC,
    MOTORCYCLE,
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

    // Max width in centimeters
    @SerialName("max-width")
    val maxWidth: Int? = null,
    // Max height in centimeters
    @SerialName("max-height")
    val maxHeight: Int? = null,

    val features: List<ParkingLotFeature>,
    @SerialName("payment-methods")
    val paymentMethods: List<PaymentMethod>,
    val comment: Map<LanguageCode, String>,
    val currency: String,
    val timezone: TimeZone,
    val rules: List<ParkingLotRule>,
)

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
                paymentMethods = listOf(
                    PaymentMethod.CASH,
                    PaymentMethod.CARD,
                    PaymentMethod.CONTACTLESS,
                ),
                comment = mapOf(
                    "pl" to
                            "Dla klientow\n" +
                            "   - kina Helios, 3 godziny parkowania bezpłatne.\n" +
                            "   - City Fit i Media Markt, 2 godziny parkowania bezpłatne. \n" +
                            "Abonament miesięczny jest dostępny, po więcej informacji: https://forumgdansk.pl/pl/przydatne-informacje/parking\n\n" +
                            "Biuro parkingu znajduje się na poziomie +2.\n" +
                            "Na parkingu znajdują się miejsca parkingowe dla osób niepełnosprawnych oraz rodzin z dziećmi.\n" +
                            "W godzinach 22:00 – 8:00 wejście na parking możliwe jest wejściem „nocnym” od strony ul. 3 Maja\n" +
                            "KASY ZNAJDUJĄ SIĘ NA KAŻDYM POZIOMIE PARKIGU PRZY WINDACH\n" +
                            "Płatności mobilne są realizowane bezdotykowo poprzez aplikację NaviPay (pobierz na Android lub IOS).\n" +
                            "W przypadku zgubienia biletu parkingowego kopię biletu można wydrukować w kasie parkingowej."
                ),
                timezone = TimeZone.of("Europe/Warsaw"),
                rules = listOf(
                    ParkingLotRule(
                        hours = "24/7",
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
                ),
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
                paymentMethods = listOf(
                    PaymentMethod.CASH,
                    PaymentMethod.CARD,
                    PaymentMethod.CONTACTLESS,
                ),
                comment = mapOf(
                    "pl" to
                            "Na dwóch najwyższych kondygnacjach budynku centrum handlowego oferujemy dwupoziomowy parking i 1100 miejsc postojowych. \n" +
                            "Wjazd do centrum handlowego odbywa się z ronda od strony ulicy Dmowskiego w Gdańsku. \n" +
                            "Komunikację między poziomami parkingowymi a poziomami handlowymi centrum handlowego zapewniają schody ruchome i windy szybkobieżne.\n" +
                            "Prosimy o zachowanie biletu parkingowego i opłacenie należności za postój w kasie automatycznej, znajdującej się przy wyjściu z parkingu.",
                    "en" to
                            "We have prepared a two-level car park with 1,100 parking spaces (including those for disabled people) for our clients.\n" +
                            "It is situated on the two top floors of the building. \n" +
                            "You can get there driving from the roundabout from the direction of Dmowskiego Street. \n" +
                            "Both levels of the car park can be reached by a spiral parking ramp. \n" +
                            "Escalators and high-speed lifts will take you from the car park decks to the Gallery's floors and back.",
                    "ru" to
                            "We have prepared a two-level car park with 1,100 parking spaces (including those for disabled people) for our clients.\n" +
                            "It is situated on the two top floors of the building. \n" +
                            "You can get there driving from the roundabout from the direction of Dmowskiego Street. \n" +
                            "Both levels of the car park can be reached by a spiral parking ramp. \n" +
                            "Escalators and high-speed lifts will take you from the car park decks to the Gallery's floors and back."
                ),
                timezone = TimeZone.of("Europe/Warsaw"),
                rules = listOf(
                    ParkingLotRule(
                        hours = "Mo-Sa 08:00-22:00; Su 09:00-21:00",
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
                )
            )
        )
    }
}
