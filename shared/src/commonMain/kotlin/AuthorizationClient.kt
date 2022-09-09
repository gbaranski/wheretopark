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

enum class AccessType(val code: String) {
    ReadMetadata("read:metadata"),
    WriteMetadata("write:metadata"),
    ReadState("read:state"),
    WriteState("write:state"),
}

fun Set<AccessType>.encode() = this.joinToString(" ") { it.code }
fun decodeAccessScope(scope: String) = scope.split(" ").map { s ->
    AccessType.values().find{ it.code == s }!!
}.toSet()

@Serializable
data class TokenResponse(
    @SerialName("access_token")
    val accessToken: String,
    @SerialName("expires_in")
    val expiresIn: Int,
    val scope: String,
    @SerialName("token_type")
    val tokenType: String

)

class AuthorizationClient(
    private val http: HttpClient
) {
    constructor(baseURL: String = DEFAULT_STOREKEEPER_URL) : this(
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
        },
    )

    suspend fun token(clientID: String, clientSecret: String, accessScope: Set<AccessType>) = http.post("/oauth/token") {
        url {
            parameters.append("client_id", clientID)
            parameters.append("client_secret", clientSecret)
            parameters.append("grant_type", "client_credentials")
            parameters.append("scope", accessScope.encode())
        }
    }.body<TokenResponse>()
}