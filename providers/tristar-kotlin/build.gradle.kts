plugins {
    kotlin("jvm")
    kotlin("plugin.serialization")
    id("com.google.cloud.tools.jib")
    application
}

val ktorVersion: String by project

application {
    mainClass.set("app.wheretopark.providers.tristar.ApplicationKt")
}

dependencies {
    implementation(project(":shared"))
    implementation(project(":providers:shared"))

    implementation("ch.qos.logback:logback-classic:1.2.11")
    implementation("io.ktor:ktor-client-core:$ktorVersion")
    implementation("io.ktor:ktor-client-cio:$ktorVersion")
    implementation("io.ktor:ktor-client-content-negotiation:$ktorVersion")
    implementation("io.ktor:ktor-serialization-kotlinx-json:$ktorVersion")
    implementation("org.jetbrains.kotlinx:kotlinx-datetime:0.4.0")
    implementation("com.charleskorn.kaml:kaml:0.47.0") // Get the latest version number from https://github.com/charleskorn/kaml/releases/latest

    testImplementation(kotlin("test"))
}

jib {
    to {
        image = "ghcr.io/wheretopark/providers-tristar"
    }
    from {
        platforms {
//            platform {
//                architecture = "amd64"
//                os = "linux"
//            }
            platform {
                architecture = "arm64"
                os = "linux"
            }
        }
    }
    container {
        creationTime = "2022-09-06T14:20:32+0000"
    }
}
