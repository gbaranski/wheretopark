package app.wheretopark.shared

import kotlinx.datetime.*
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertTrue

class ParkingLotStatusTest {
    private fun create(vararg rules: Pair<ParkingLotWeekdays?, ParkingLotHours?>) =
        ParkingLotMetadata(
            name = "",
            address = "",
            location = Coordinate(latitude = 0.0, longitude = 0.0),
            resources = listOf(),
            totalSpots = 0u,
            features = listOf(),
            currency = "",
            rules = rules.map { (weekdays, hours) ->
                ParkingLotRule(
                    weekdays,
                    hours,
                    pricing = listOf()
                )
            }
        )

    private fun status(weekday: DayOfWeek, hour: Int, minute: Int): ParkingLotStatus {
        val metadata = create(
            Pair(
                ParkingLotWeekdays(start = DayOfWeek.MONDAY, end = DayOfWeek.SATURDAY),
                ParkingLotHours(start = LocalTime(8, 0, 0, 0), end = LocalTime(22, 0, 0, 0))
            ),
            Pair(
                ParkingLotWeekdays(start = DayOfWeek.SUNDAY, end = DayOfWeek.SUNDAY),
                ParkingLotHours(start = LocalTime(9, 0, 0, 0), end = LocalTime(21, 0, 0, 0))
            ),
        )
        // If you don't know what happens, look at August 2022, and how weekdays perfectly matches
        val date = LocalDate(2022, 8, weekday.isoDayNumber)
        val time = LocalTime(hour, minute, 0)
        return metadata.status(LocalDateTime(date, time).toInstant(TimeZone.UTC))
    }

    @Test
    fun test() {
        assertEquals(status(DayOfWeek.MONDAY, 7, 30), ParkingLotStatus.OPENS_SOON)
        assertEquals(status(DayOfWeek.MONDAY, 8, 0), ParkingLotStatus.OPEN)
        assertEquals(status(DayOfWeek.MONDAY, 12, 0), ParkingLotStatus.OPEN)
        assertEquals(status(DayOfWeek.MONDAY, 21, 30), ParkingLotStatus.CLOSES_SOON)
        assertEquals(status(DayOfWeek.MONDAY, 22, 0), ParkingLotStatus.CLOSED)
        assertEquals(status(DayOfWeek.MONDAY, 23, 0), ParkingLotStatus.CLOSED)

        assertEquals(status(DayOfWeek.SUNDAY, 8, 30), ParkingLotStatus.OPENS_SOON)
        assertEquals(status(DayOfWeek.SUNDAY, 9, 0), ParkingLotStatus.OPEN)
        assertEquals(status(DayOfWeek.SUNDAY, 12, 0), ParkingLotStatus.OPEN)
        assertEquals(status(DayOfWeek.SUNDAY, 20, 30), ParkingLotStatus.CLOSES_SOON)
        assertEquals(status(DayOfWeek.SUNDAY, 21, 0), ParkingLotStatus.CLOSED)
        assertEquals(status(DayOfWeek.SUNDAY, 22, 0), ParkingLotStatus.CLOSED)
    }
}