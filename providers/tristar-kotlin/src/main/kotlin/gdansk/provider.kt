package app.wheretopark.providers.tristar.gdansk

import app.wheretopark.providers.shared.Provider
import app.wheretopark.shared.*
import com.charleskorn.kaml.Yaml
import io.ktor.client.*
import io.ktor.client.call.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.client.request.*
import io.ktor.serialization.kotlinx.json.*
import kotlinx.datetime.TimeZone
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.decodeFromString
import kotlinx.serialization.json.Json
import java.util.*
import kotlin.time.Duration
import kotlin.time.Duration.Companion.seconds

const val METADATA_URL =
    "https://ckan.multimediagdansk.pl/dataset/cb1e2708-aec1-4b21-9c8c-db2626ae31a6/resource/d361dff3-202b-402d-92a5-445d8ba6fd7f/download/parking-lots.json"
const val STATE_URL = "https://ckan2.multimediagdansk.pl/parkingLots"

@Serializable
data class ParkingLotConfiguration(
    val resources: List<ParkingLotResource>,
    @SerialName("total-spots")
    val totalSpots: Map<ParkingSpotType, UInt>,
    @SerialName("max-height")
    val maxHeight: Int? = null,
    val features: List<ParkingLotFeature>,
    @SerialName("payment-methods")
    val paymentMethods: List<PaymentMethod> = listOf(),
    val comment: Map<LanguageCode, String> = mapOf(),
    val rules: List<ParkingLotRule>,
)

@Serializable
data class Configuration(
    @SerialName("parking-lots")
    val parkingLots: Map<String, ParkingLotConfiguration>,
)

class TristarGdanskProvider : Provider() {
    override val name: String
        get() = "tristar/gdansk"
    override val metadataInterval: Duration
        get() = 30.seconds
    override val stateInterval: Duration
        get() = 30.seconds

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

    private val configuration =
        Yaml().decodeFromString<Configuration>(this.javaClass.getResource("/configuration-gdansk.yaml")!!.readText())

    override suspend fun metadatas(): Map<ParkingLotID, ParkingLotMetadata> {
        val vendorMetadatas = client.get(METADATA_URL).body<MetadataResponse>()
        return vendorMetadatas.parkingLots.map {
            val configuration = configuration.parkingLots[it.id] ?: return@map null
            val id = it.location.hash()
            ids[it.id] = id
            val metadata = ParkingLotMetadata(
                name = it.name,
                address = it.address,
                location = it.location,
                resources = configuration.resources,
                totalSpots = configuration.totalSpots,
                features = configuration.features,
                currency = "PLN",
                paymentMethods = configuration.paymentMethods,
                comment = configuration.comment,
                timezone = TimeZone.of("Europe/Warsaw"),
                rules = configuration.rules,
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
                availableSpots = mapOf(
                    Pair(ParkingSpotType.CAR, it.availableSpots)
                ),
            )
            id to state
        }
            .filterNotNull()
            .toMap()
    }
}