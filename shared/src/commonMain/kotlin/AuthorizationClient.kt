package app.wheretopark.shared

import io.ktor.client.*
import io.ktor.client.call.*
import io.ktor.client.plugins.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.client.request.*
import io.ktor.serialization.kotlinx.json.*
import kotlinx.serialization.KSerializer
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.descriptors.PrimitiveKind
import kotlinx.serialization.descriptors.PrimitiveSerialDescriptor
import kotlinx.serialization.descriptors.SerialDescriptor
import kotlinx.serialization.encoding.Decoder
import kotlinx.serialization.encoding.Encoder
import kotlinx.serialization.json.Json

enum class AccessType(name: String) {
    ReadMetadata("read:metadata"),
    WriteMetadata("write:metadata"),
    ReadState("read:state"),
    WriteState("write:state"),
}

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

    suspend fun token(clientID: String, clientSecret: String, accessType: Set<AccessType>) = http.get("/oauth/token") {
        url {
            parameters.append("client_id", clientID)
            parameters.append("client_secret", clientSecret)
            parameters.append("grant_type", "client_credentials")
            val scopeString = accessType.joinToString(" ") { it.name }
            parameters.append("scope", scopeString)
        }
    }.body<TokenResponse>()
}