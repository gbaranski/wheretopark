pluginManagement {
    repositories {
        google()
        mavenCentral()
        maven(url = "https://plugins.gradle.org/m2/")
    }
}

rootProject.name = "wheretopark"

include(":shared")
include(":shared-client")
include(":android")
include(":storekeeper")

