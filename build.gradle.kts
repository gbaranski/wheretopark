plugins {
    kotlin("multiplatform") version "1.7.0" apply false
    kotlin("android") version "1.7.0" apply false
    kotlin("plugin.serialization") version "1.7.0" apply false
    id("com.android.application") version "7.2.2" apply false
    id("com.android.library") version "7.2.2" apply false
    id("com.google.cloud.tools.jib") version "3.2.1" apply false
    
}

group = "app.wheretopark"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
    maven("https://maven.pkg.jetbrains.space/public/p/kotlinx-html/maven")
}

allprojects {
    repositories {
        google()
        mavenCentral()
    }

    tasks.withType(AbstractTestTask::class) {
        testLogging {
            showStandardStreams = true
            events("passed", "failed")
        }
    }
}