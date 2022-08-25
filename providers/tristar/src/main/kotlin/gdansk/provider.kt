package app.wheretopark.providers.tristar.gdansk

import app.wheretopark.providers.tristar.Provider
import app.wheretopark.shared.*
import com.charleskorn.kaml.Yaml
import io.ktor.client.*
import io.ktor.client.call.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.client.request.*
import io.ktor.serialization.kotlinx.json.*
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.decodeFromString
import kotlinx.serialization.json.Json
import java.time.DayOfWeek

const val METADATA_URL = "https://ckan.multimediagdansk.pl/dataset/cb1e2708-aec1-4b21-9c8c-db2626ae31a6/resource/d361dff3-202b-402d-92a5-445d8ba6fd7f/download/parking-lots.json"
const val STATE_URL = "https://ckan2.multimediagdansk.pl/parkingLots"

@Serializable
data class ParkingLotConfigurationRule(
    val weekdays: String?,
    val hours: String?,
    val pricing: List<ParkingLotPricingRule>
)

@Serializable
data class ParkingLotConfiguration(
    val emails: List<String>,
    @SerialName("phone-numbers")
    val phoneNumbers: List<String>,
    val websites: List<String>,
    @SerialName("total-spots")
    val totalSpots: UInt,
    val features: List<ParkingLotFeature>,
    val rules: List<ParkingLotConfigurationRule>
)

@Serializable
data class Configuration(
    @SerialName("parking-lots")
    val parkingLots: Map<ParkingLotID, ParkingLotConfiguration>,
)

class TristarGdanskProvider: Provider() {
    private val client = HttpClient {
        install(ContentNegotiation) {
            json(Json {
                prettyPrint = true
                isLenient = true
                ignoreUnknownKeys = true
            })
        }
    }
    // Mapping of vendor and storekeeper ID's
    private val ids = mutableMapOf<String, ParkingLotID>()

    private val configuration = Yaml().decodeFromString<Configuration>(this.javaClass.getResource("/configuration-gdansk.yaml")!!.readText())

    override suspend fun metadatas(): Map<ParkingLotID, ParkingLotMetadata> {
        val vendorMetadatas = client.get(METADATA_URL).body<MetadataResponse>()
        return vendorMetadatas.parkingLots.map { it ->
            val configuration = configuration.parkingLots[it.id] ?: return@map null
            val id = it.location.hash()
            ids[it.id] = id
            val metadata = ParkingLotMetadata(
                name = it.name,
                address = it.address,
                location = it.location,
                emails = configuration.emails,
                phoneNumbers = configuration.phoneNumbers,
                websites = configuration.websites,
                totalSpots = configuration.totalSpots,
                features = configuration.features,
                currency = "PLN",
                rules = configuration.rules.map {
                    val weekdaysSplit = it.weekdays?.split("-")
                    val hoursSplit = it.hours?.split("-")
                    ParkingLotRule(
                        weekdays = if (weekdaysSplit != null) {
                            assert(weekdaysSplit.count() == 2)
                            ParkingLotWeekdays(
                                start = DayOfWeek.valueOf(weekdaysSplit[0].uppercase()),
                                end = DayOfWeek.valueOf(weekdaysSplit[1].uppercase()),
                            )
                        } else { null },
                        hours = if (hoursSplit != null) {
                            assert(hoursSplit.count() == 2)
                            ParkingLotHours(
                                start = hoursSplit[0],
                                end = hoursSplit[1],
                            )
                        } else { null },
                        pricing = it.pricing,
                    )
                }
            )
            id to metadata
        }
        .filterNotNull()
        .toMap()
    }

    override suspend fun states(): Map<ParkingLotID, ParkingLotState> {
        val vendorStates = client.get(STATE_URL).body<StateResponse>()
        return vendorStates.parkingLots.map {
            val id = ids[it.id] ?: return@map null
            val state = ParkingLotState(
                lastUpdated = it.lastUpdate,
                availableSpots = it.availableSpots,
            )
            id to state
        }
        .filterNotNull()
        .toMap()
    }
}