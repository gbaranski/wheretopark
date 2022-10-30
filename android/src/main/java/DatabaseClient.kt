package app.wheretopark.android

import io.ktor.client.*
import io.ktor.client.call.*
import io.ktor.client.plugins.*
import io.ktor.client.plugins.auth.*
import io.ktor.client.plugins.auth.providers.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.client.request.*
import io.ktor.http.*
import io.ktor.serialization.kotlinx.json.*
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json

@Serializable
data class DatabaseParkingLot (
    val id: String,
    val metadata: ParkingLotMetadata,
    val state: ParkingLotState,
)


@Serializable
data class Response<T>(
    val time: String,
    val status: String,
    val result: List<T>,
)

class DatabaseClient(
    baseURL: String,
    namespace: String,
    database: String,
    username: String,
    password: String,
) {
    private val http = HttpClient{
        defaultRequest {
            url(baseURL)
            headers {
                append("NS", namespace)
                append("DB", database)
            }
        }
        expectSuccess = true
        install(ContentNegotiation) {
            json(Json {
                prettyPrint = true
                isLenient = true
                ignoreUnknownKeys = true
            })
        }
        install(HttpRequestRetry) {
            retryOnServerErrors(maxRetries = 5)
            exponentialDelay()
        }
        install(Auth) {
            basic {
                credentials {
                    BasicAuthCredentials(username = username, password = password)
                }
            }

        }
    }

    private suspend inline fun <reified T> sql(query: String): Response<T> {
        return http.post("/sql") {
            contentType(ContentType.Text.Plain)
            setBody(query)
        }.body()
    }

    @Throws(Exception::class)
    suspend fun parkingLots(): Map<ParkingLotID, ParkingLot> {
        val parkingLots = sql<DatabaseParkingLot>("SELECT * FROM parking_lot")
        return parkingLots.result.associate {
            val id = it.id.split(":")[1]
            id to ParkingLot(
                metadata = it.metadata,
                state = it.state,
            )
        }
    }
}
