package app.wheretopark.providers.cctv

import org.bytedeco.javacpp.Loader
import org.bytedeco.javacv.FFmpegFrameGrabber
import org.bytedeco.javacv.Java2DFrameConverter
import org.bytedeco.opencv.global.opencv_core
import org.opencv.core.Mat
import org.opencv.core.MatOfPoint
import org.opencv.core.MatOfPoint2f
import org.opencv.core.Point
import org.opencv.core.RotatedRect
import org.opencv.imgproc.Imgproc
import java.awt.image.BufferedImage


val converter = Java2DFrameConverter()
fun captureImage(path: String): BufferedImage {
    val grabber = FFmpegFrameGrabber(path)
    grabber.start()
    val frame = grabber.grab()!!
    grabber.stop()
    return converter.convert(frame)!!
}

fun Configuration.ParkingSpot.cvPoints(): MatOfPoint2f {
    val points = points.map {
        Point(it.x.toDouble(), it.y.toDouble())
    }
    val matOfPoints = MatOfPoint2f()
    matOfPoints.fromList(points)
    return matOfPoints

}

fun Configuration.ParkingSpot.minAreaRectangle(): RotatedRect {
    val points = this.cvPoints()
    return Imgproc.minAreaRect(points)
}

fun loadOpenCV() = Loader.load(opencv_core::class.java)

