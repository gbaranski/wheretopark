plugins {
    kotlin("jvm")
    kotlin("plugin.serialization")
    id("io.ktor.plugin") version "2.1.0"
    application
}

val ktorVersion: String by project

application {
    mainClass.set("app.wheretopark.providers.cctv.ApplicationKt")
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
    implementation("org.bytedeco:javacv:1.5.7")
    implementation("org.bytedeco:opencv-platform:4.5.5-1.5.7")
    implementation("org.bytedeco:ffmpeg-platform:5.0-1.5.7")


    testImplementation(kotlin("test"))
}

