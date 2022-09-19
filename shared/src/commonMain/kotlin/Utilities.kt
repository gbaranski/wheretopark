@file:OptIn(ExperimentalJsExport::class)

package app.wheretopark.shared

import io.ktor.http.*
import kotlinx.serialization.Serializable
import kotlin.js.ExperimentalJsExport
import kotlin.js.JsExport
import kotlin.math.*

typealias LanguageCode = String

const val BASE_WEBAPP_URL = "https://web.wheretopark.app"

fun getShareURL(id: ParkingLotID): String {
    val url = URLBuilder(BASE_WEBAPP_URL).appendPathSegments("parking-lot", id)
    return url.buildString()
}

private val BITS = intArrayOf(16, 8, 4, 2, 1)

// note: no a,i,l, and o
private val BASE32_CHARS = charArrayOf(
    '0',
    '1',
    '2',
    '3',
    '4',
    '5',
    '6',
    '7',
    '8',
    '9',
    'b',
    'c',
    'd',
    'e',
    'f',
    'g',
    'h',
    'j',
    'k',
    'm',
    'n',
    'p',
    'q',
    'r',
    's',
    't',
    'u',
    'v',
    'w',
    'x',
    'y',
    'z'
)

fun toRadians(v: Double) = v * PI / 180.0
private const val earthRadiusMeters: Double = 6371000.0


@Serializable
@JsExport
data class Coordinate(
    val latitude: Double,
    val longitude: Double,
) {
    /**
     * Haversine formula. Giving great-circle distances between two points on a sphere from their longitudes and latitudes.
     * It is a special case of a more general formula in spherical trigonometry, the law of haversines, relating the
     * sides and angles of spherical "triangles".
     *
     * https://rosettacode.org/wiki/Haversine_formula#Java
     *
     * @return Distance in meters
     */
    fun distanceTo(to: Coordinate): Double {
        val dLat = toRadians(to.latitude - this.latitude);
        val dLon = toRadians(to.latitude - this.latitude);
        val originLat = toRadians(this.latitude);
        val destinationLat = toRadians(to.latitude);

        val a = sin(dLat / 2).pow(2) + sin(dLon / 2).pow(2) * cos(originLat) * cos(destinationLat);

        val c = 2 * asin(sqrt(a));
        return earthRadiusMeters * c;
    }

    fun hash(length: Int = 12): String {
        if (length < 1 || length > 12) {
            throw IllegalArgumentException("length must be between 1 and 12")
        }
        val latInterval = doubleArrayOf(-90.0, 90.0)
        val lonInterval = doubleArrayOf(-180.0, 180.0)

        val geohash = StringBuilder()
        var isEven = true
        var bit = 0
        var ch = 0

        while (geohash.length < length) {
            if (isEven) {
                val mid = (lonInterval[0] + lonInterval[1]) / 2
                if (longitude > mid) {
                    ch = ch or BITS[bit]
                    lonInterval[0] = mid
                } else {
                    lonInterval[1] = mid
                }
            } else {
                val mid: Double = (latInterval[0] + latInterval[1]) / 2
                if (latitude > mid) {
                    ch = ch or BITS[bit]
                    latInterval[0] = mid
                } else {
                    latInterval[1] = mid
                }
            }

            isEven = !isEven

            if (bit < 4) {
                bit++
            } else {
                geohash.append(BASE32_CHARS[ch])
                bit = 0
                ch = 0
            }
        }
        return geohash.toString()
    }
}