@file:OptIn(ExperimentalJsExport::class)

package app.wheretopark.shared

import io.ktor.client.*
import io.ktor.client.call.*
import io.ktor.client.plugins.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.client.request.*
import io.ktor.serialization.kotlinx.json.*
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json
import kotlin.coroutines.cancellation.CancellationException
import kotlin.js.ExperimentalJsExport
import kotlin.js.JsExport

@JsExport
enum class AccessType(val code: String) {
    ReadMetadata("read:metadata"),
    WriteMetadata("write:metadata"),
    ReadState("read:state"),
    WriteState("write:state"),
    ReadStatus("read:status"),
}

fun Set<AccessType>.encode() = this.joinToString(" ") { it.code }
fun decodeAccessScope(scope: String) = scope.split(" ").map { s ->
    AccessType.values().find { it.code == s }!!
}.toSet()

@Serializable
@JsExport
data class TokenResponse(
    @SerialName("access_token")
    val accessToken: String,
    @SerialName("expires_in")
    val expiresIn: Int,
    val scope: String,
    @SerialName("token_type")
    val tokenType: String

)

const val DEFAULT_AUTHORIZATION_URL = "https://authorization.wheretopark.app"

class AuthorizationClient(
    private val http: HttpClient,
    private val clientID: String,
    private val clientSecret: String,
) {
    constructor(baseURL: String = DEFAULT_AUTHORIZATION_URL, clientID: String, clientSecret: String) : this(
        HttpClient {
            defaultRequest {
                url(baseURL)
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
        },
        clientID,
        clientSecret
    )

    @Throws(
        RedirectResponseException::class, ClientRequestException::class, ServerResponseException::class,
        CancellationException::class
    )
    suspend fun token(scope: Set<AccessType>) = http.post("/oauth/token") {
        url {
            parameters.append("client_id", clientID)
            parameters.append("client_secret", clientSecret)
            parameters.append("grant_type", "client_credentials")
            parameters.append("scope", scope.encode())
        }
    }.body<TokenResponse>()
}