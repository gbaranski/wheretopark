@file:OptIn(ExperimentalJsExport::class)

package app.wheretopark.shared

import io.ktor.http.*
import kotlinx.serialization.KSerializer
import kotlinx.serialization.Serializable
import kotlinx.serialization.builtins.ListSerializer
import kotlinx.serialization.builtins.MapSerializer
import kotlinx.serialization.builtins.serializer
import kotlinx.serialization.descriptors.PrimitiveKind
import kotlinx.serialization.descriptors.PrimitiveSerialDescriptor
import kotlinx.serialization.descriptors.SerialDescriptor
import kotlinx.serialization.encoding.Decoder
import kotlinx.serialization.encoding.Encoder
import kotlin.js.ExperimentalJsExport
import kotlin.js.JsExport

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

@Serializable
@JsExport
data class Coordinate(
    val latitude: Double,
    val longitude: Double,
) {
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