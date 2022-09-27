import app.wheretopark.shared.ParkingLotMetadata
import app.wheretopark.shared.ParkingLotStatus
import ch.poole.openinghoursparser.OpeningHoursParser
import io.leonard.OpeningHoursEvaluator
import kotlinx.datetime.toJavaZoneId
import java.time.LocalDateTime

fun ParkingLotMetadata.statusAt(dateTime: LocalDateTime): ParkingLotStatus {
    val zonedDateTime = dateTime.atZone(timezone.toJavaZoneId()).toLocalDateTime()
    val sum = rules.joinToString(",") { it.hours }
    val parser = OpeningHoursParser(sum.byteInputStream())
    val rules = parser.rules(false)
    return if (OpeningHoursEvaluator.isOpenAt(zonedDateTime, rules)) {
        if (OpeningHoursEvaluator.isOpenAt(zonedDateTime.plusHours(1), rules)) ParkingLotStatus.OPEN
        else ParkingLotStatus.CLOSES_SOON
    } else {
        if (OpeningHoursEvaluator.isOpenAt(zonedDateTime.plusHours(1), rules)) ParkingLotStatus.OPENS_SOON
        else ParkingLotStatus.CLOSED
    }
}