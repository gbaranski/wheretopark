package app.wheretopark.storekeeper

import app.wheretopark.shared.*
import com.auth0.jwt.JWT
import com.auth0.jwt.algorithms.Algorithm
import io.ktor.http.*
import io.ktor.serialization.kotlinx.json.json
import io.ktor.server.application.Application
import io.ktor.server.application.call
import io.ktor.server.application.install
import io.ktor.server.auth.*
import io.ktor.server.auth.jwt.*
import io.ktor.server.engine.*
import io.ktor.server.netty.*
import io.ktor.server.plugins.autohead.*
import io.ktor.server.plugins.callloging.*
import io.ktor.server.plugins.contentnegotiation.ContentNegotiation
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import java.net.URI

data class AuthorizationConfig(
    val clientID: String,
    val clientSecret: String,
    val jwtSecret: String,
)

data class StoreConfig(val uri: URI)

data class Config(
    val authorization: AuthorizationConfig,
    val store: StoreConfig,
    val port: Int,
)

fun loadConfig() = Config(
    authorization = AuthorizationConfig(
        clientID = System.getenv("AUTHORIZATION_CLIENT_ID")!!,
        clientSecret = System.getenv("AUTHORIZATION_CLIENT_SECRET")!!,
        jwtSecret = System.getenv("JWT_SECRET")!!,
    ),
    store = StoreConfig(
        uri = URI(System.getenv("STORE_URI") ?: "memory:/")
    ),
    port = System.getenv("PORT")?.toInt() ?: 8080,
)


fun main() {
    val config = loadConfig()
    val store = when (config.store.uri.scheme) {
        "memory" -> {
            MemoryStore()
        }

        "redis" -> {
            RedisStore(config.store.uri.host, if (config.store.uri.port == -1) 6379 else config.store.uri.port)
        }

        else -> {
            throw IllegalArgumentException("Unknown store scheme: ${config.store.uri.scheme}")
        }
    }
    embeddedServer(Netty, port = config.port) {
        configure(store, config)
    }.start(wait = true)
}

fun Application.configure(store: Store, config: Config) {
    install(ContentNegotiation) {
        json()
    }
    install(CallLogging)
    install(AutoHeadResponse)
    install(Authentication) {
        jwt("auth-jwt") {
            realm = "Storekeeper service"
            verifier(
                JWT.require(Algorithm.HMAC512(config.authorization.jwtSecret))
                    .withAudience(config.authorization.clientID).build()
            )
            validate { credential ->
                JWTPrincipal(credential.payload)
            }
        }
    }

    routing {
        get("/health-check") {
            call.respond("Hello, World!")
        }
        authenticate("auth-jwt") {
            get("/parking-lot/state") {
                call.respond(store.getStates())
            }
            get("/parking-lot/metadata") {
                call.respond(store.getMetadatas())
            }
        }

        authenticate("auth-jwt") {
            post("/parking-lot/state") {
                val newStates = call.receive<Map<ParkingLotID, ParkingLotState>>()
                store.updateStates(newStates)
                call.respondText("updated ${newStates.count()} states", status = HttpStatusCode.Accepted)
            }

            post("/parking-lot/metadata") {
                val newMetadatas = call.receive<Map<ParkingLotID, ParkingLotMetadata>>()
                store.updateMetadatas(newMetadatas)
                call.respondText("updated ${newMetadatas.count()} metadatas", status = HttpStatusCode.Accepted)
            }
        }
    }
}