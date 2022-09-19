package app.wheretopark.providers.cctv

//import org.opencv.core.Mat
//import org.opencv.highgui.HighGui.imshow
//import org.opencv.videoio.VideoCapture
import app.wheretopark.providers.shared.Provider
import app.wheretopark.shared.ParkingLotID
import app.wheretopark.shared.ParkingLotMetadata
import app.wheretopark.shared.ParkingLotState
import app.wheretopark.shared.ParkingSpotType
import com.charleskorn.kaml.Yaml
import kotlinx.serialization.decodeFromString
import org.bytedeco.javacv.CanvasFrame
import org.bytedeco.javacv.FFmpegFrameGrabber
import org.bytedeco.javacv.Java2DFrameConverter
import org.bytedeco.librealsense.frame
import org.bytedeco.opencv.global.opencv_imgproc.*
import org.opencv.core.*
import org.opencv.imgcodecs.Imgcodecs.imencode
import java.awt.Graphics2D
import java.awt.Rectangle
import java.awt.geom.AffineTransform
import java.awt.geom.GeneralPath
import java.awt.image.BufferedImage
import java.io.ByteArrayInputStream
import java.io.File
import javax.imageio.ImageIO
import javax.swing.WindowConstants
import kotlin.math.min
import kotlin.time.Duration
import kotlin.time.Duration.Companion.minutes
import kotlin.time.Duration.Companion.seconds


class CCTVProvider : Provider() {
    private val configuration =
        Yaml().decodeFromString<Configuration>(
            object {}.javaClass.getResource("/configuration.yaml")!!.readText()
        )

    override val name: String
        get() = "cctv"
    override val metadataInterval: Duration
        get() = 30.minutes
    override val stateInterval: Duration
        get() = 30.seconds

    override suspend fun metadatas() = configuration.parkingLots.associate {
        val id = it.location.hash()
        id to ParkingLotMetadata(
            name = it.name,
            address = it.address,
            location = it.location,
            resources = it.resources,
            totalSpots = mapOf(
                ParkingSpotType.CAR to it.spots.count().toUInt()
            ),
            features = it.features,
            comment = it.comment,
            currency = it.currency,
            rules = it.rules,
        )
    }

    override suspend fun states(): Map<ParkingLotID, ParkingLotState> {
//        val canvas = CanvasFrame("Visualisation", 1.0)
//            canvas.showImage(image)
        configuration.parkingLots[0].let {
            val image = captureImage(it.cameraURL)
            println("got image from ${it.cameraURL}")
            val predictions = it.spots.mapIndexed { index, spot ->
                val minAreaRectangle = spot.minAreaRectangle()
                print("min area rect x=${minAreaRectangle.size.width} y=${minAreaRectangle.size.height}")
//                val boundingRectangle = boundingRect()
















//                val clip = GeneralPath()
//                spot.points.forEachIndexed { index, point ->
//                    if (index == 0) clip.moveTo(point.x.toDouble(), point.y.toDouble())
//                    else clip.lineTo(point.x.toDouble(), point.y.toDouble())
//                }
//                clip.closePath()
//
//                val bounds: Rectangle = clip.bounds
//                val img = BufferedImage(bounds.width, bounds.height, BufferedImage.TYPE_INT_ARGB)
//                val g2d: Graphics2D = img.createGraphics()
//                clip.transform(AffineTransform.getTranslateInstance(-65, -123))
//                g2d.clip = clip
//                g2d.translate(-65, -123)
//                g2d.drawImage(image, 0, 0, null)
//                g2d.dispose()
//
//                ImageIO.write(img, "png", File(""))
//
//
//
//                val points = Mat(spot.points.count(), 2)
//                // (1) Crop the bounding rect
//                val boundingRectangle = boundingRect(points)
//                val cropped = Mat(frame(boundingRectangle).address())
//                val matOfByte = MatOfByte()
//                imencode(".jpg", cropped, matOfByte);
//                val byteArray = matOfByte.toArray();
//                val input = ByteArrayInputStream(byteArray)
//                val bufImage = ImageIO.read(input)
//                canvas.showImage(bufImage)
            }
        }
        return mapOf()
    }

}

suspend fun main() {
    loadOpenCV()

    val provider = CCTVProvider()
    provider.states()
}