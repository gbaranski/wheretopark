plugins {
    kotlin("jvm")
    kotlin("plugin.serialization")
    id("com.google.cloud.tools.jib")
    application
}

val ktorVersion: String by project

application {
    mainClass.set("app.wheretopark.storekeeper.ApplicationKt")
}

dependencies {
    implementation(project(":shared"))

    implementation("io.ktor:ktor-server-netty:$ktorVersion")
    implementation("io.ktor:ktor-server-content-negotiation:$ktorVersion")
    implementation("io.ktor:ktor-serialization-kotlinx-json:$ktorVersion")
    implementation("io.ktor:ktor-server-call-logging:$ktorVersion")

    implementation("ch.qos.logback:logback-classic:1.2.11")
    implementation("io.github.crackthecodeabhi:kreds:0.8")
    implementation("org.jetbrains.kotlinx:kotlinx-datetime:0.4.0")

    testImplementation(kotlin("test"))
    testImplementation("io.ktor:ktor-server-tests:$ktorVersion")
}

jib {
    to {
        image = "registry.gbaranski.com/wheretopark-storekeeper"
    }
    from {
        platforms {
            platform {
                architecture = "amd64"
                os = "linux"
            }
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
