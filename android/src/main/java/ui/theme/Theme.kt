package app.wheretopark.android.ui.theme

import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.material.MaterialTheme
import androidx.compose.material.darkColors
import androidx.compose.material.lightColors
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color

private val DarkColorPalette = darkColors(
    primary = Color(0xFF_F7B500),
    primaryVariant = Color(0xFF_FFE74C),
    secondary = Color(0xFF_39393D),
    secondaryVariant = Color(0xFF_98989E),
    background = Color(0xFF_1C1C1E),
    surface = Color(0xFF1C1C1E)
)

private val LightColorPalette = lightColors(
    primary = Color(0xFF_F7B500),
    primaryVariant = Color(0xFF_FFE74C),
    secondary = Color(0xFF_E9E9E9),
    secondaryVariant = Color(0xFF_8A8A8E)

    /* Other default colors to override
    background = Color.White,
    surface = Color.White,
    onPrimary = Color.White,
    onSecondary = Color.Black,
    onBackground = Color.Black,
    onSurface = Color.Black,
    */
)

@Composable
fun WheretoparkTheme(darkTheme: Boolean = isSystemInDarkTheme(), content: @Composable () -> Unit) {
    val colors = if (darkTheme) {
        DarkColorPalette
    } else {
        LightColorPalette
    }

    MaterialTheme(
        colors = colors,
        typography = Typography,
        shapes = Shapes,
        content = content
    )
}