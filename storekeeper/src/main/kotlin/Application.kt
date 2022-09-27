package app.wheretopark.storekeeper

import app.wheretopark.shared.*
import com.auth0.jwt.JWT
import com.auth0.jwt.algorithms.Algorithm
import io.ktor.http.*
import io.ktor.serialization.kotlinx.json.*
import io.ktor.server.application.*
import io.ktor.server.auth.*
import io.ktor.server.auth.jwt.*
import io.ktor.server.engine.*
import io.ktor.server.netty.*
import io.ktor.server.plugins.*
import io.ktor.server.plugins.autohead.*
import io.ktor.server.plugins.callloging.*
import io.ktor.server.plugins.contentnegotiation.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import statusAt
import java.net.URI
import java.time.LocalDateTime
import java.time.ZoneId

data class Config(
    val jwtSecret: String,
    val storeURI: URI,
    val port: Int,
)

fun loadConfig() = Config(
    jwtSecret = System.getenv("JWT_SECRET")!!,
    storeURI = URI(System.getenv("STORE_URI") ?: "memory:/"),
    port = System.getenv("PORT")?.toInt() ?: 8080,
)


fun main() {
    val config = loadConfig()
    val store = when (config.storeURI.scheme) {
        "memory" -> {
            MemoryStore()
        }

        "redis" -> {
            RedisStore(config.storeURI.host, if (config.storeURI.port == -1) 6379 else config.storeURI.port)
        }

        else -> {
            throw IllegalArgumentException("Unknown store scheme: ${config.storeURI.scheme}")
        }
    }
    embeddedServer(Netty, port = config.port) {
        configure(store, config)
    }.start(wait = true)
}

// TODO: Find a cleaner way to do it
suspend fun validateScope(call: ApplicationCall, accessType: AccessType): Boolean {
    val principal = call.principal<JWTPrincipal>()
    val scope = principal?.getClaim("scope", String::class) ?: ""
    val accessScope = decodeAccessScope(scope)
    return if (accessScope.contains(accessType)) {
        true
    } else {
        call.respond(HttpStatusCode.Unauthorized, "missing access ${accessType.code} in scope $scope")
        false
    }

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
                JWT.require(Algorithm.HMAC512(config.jwtSecret)).build()
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
            route("/parking-lot") {
                get {
                    if (!validateScope(call, AccessType.ReadMetadata)) return@get
                    if (!validateScope(call, AccessType.ReadState)) return@get
                    val metadatas = store.getMetadatas()
                    val states = store.getStates()
                    val parkingLots = metadatas.map { (id, metadata) ->
                        states[id]?.let {
                            id to ParkingLot(metadata, it)
                        }
                    }.filterNotNull().toMap()
                    call.respond(parkingLots)
                }
                post {
                    if (!validateScope(call, AccessType.WriteMetadata)) return@post
                    if (!validateScope(call, AccessType.WriteState)) return@post
                    val parkingLots = call.receive<Map<ParkingLotID, ParkingLot>>()
                    val (metadatas, states) = parkingLots.split()
                    store.updateMetadatas(metadatas)
                    store.updateStates(states)
                    call.respondText("updated ${parkingLots.count()} parking lots", status = HttpStatusCode.Accepted)
                }
                route("/metadata") {
                    get {
                        if (!validateScope(call, AccessType.ReadMetadata)) return@get
                        call.respond(store.getMetadatas())
                    }
                    post {
                        if (!validateScope(call, AccessType.WriteMetadata)) return@post
                        val newMetadatas = call.receive<Map<ParkingLotID, ParkingLotMetadata>>()
                        store.updateMetadatas(newMetadatas)
                        call.respondText("updated ${newMetadatas.count()} metadatas", status = HttpStatusCode.Accepted)
                    }
                }
                route("/state") {
                    get {
                        if (!validateScope(call, AccessType.ReadState)) return@get
                        call.respond(store.getStates())
                    }
                    post {
                        if (!validateScope(call, AccessType.WriteState)) return@post
                        val newStates = call.receive<Map<ParkingLotID, ParkingLotState>>()
                        store.updateStates(newStates)
                        call.respondText("updated ${newStates.count()} states", status = HttpStatusCode.Accepted)
                    }
                }
                get("/status") {
                    if (!validateScope(call, AccessType.ReadStatus)) return@get
                    val metadatas = store.getMetadatas()
                    val now = LocalDateTime.now(ZoneId.of("UTC"))
                    val statuses = metadatas.map { (id, metadata) ->
                        id to metadata.statusAt(now)
                    }.toMap()
                    call.respond(statuses)
                }
                route("/{id}") {
                    get {
                        if (!validateScope(call, AccessType.ReadMetadata)) return@get
                        if (!validateScope(call, AccessType.ReadState)) return@get
                        val id = call.parameters["id"] ?: throw BadRequestException("Missing ID")
                        val metadata = store.getMetadata(id) ?: throw NotFoundException("Metadata not found for $id")
                        val state = store.getState(id) ?: throw NotFoundException("State not found for $id")
                        call.respond(ParkingLot(metadata, state))
                    }
                    get("/metadata") {
                        if (!validateScope(call, AccessType.ReadMetadata)) return@get
                        val id = call.parameters["id"] ?: throw BadRequestException("Missing ID")
                        val metadata = store.getMetadata(id) ?: throw NotFoundException("Metadata not found for $id")
                        call.respond(metadata)
                    }
                    get("/state") {
                        if (!validateScope(call, AccessType.ReadState)) return@get
                        val id = call.parameters["id"] ?: throw BadRequestException("Missing ID")
                        val state = store.getState(id) ?: throw NotFoundException("State not found for $id")
                        call.respond(state)
                    }
                    get("/status") {
                        if (!validateScope(call, AccessType.ReadStatus)) return@get
                        val id = call.parameters["id"] ?: throw BadRequestException("Missing ID")
                        val metadata = store.getMetadata(id) ?: throw NotFoundException("Metadata not found for $id")
                        val now = LocalDateTime.now(ZoneId.of("UTC"))
                        val status = metadata.statusAt(now)
                        call.respond(status)
                    }
                }
            }
        }
    }
}