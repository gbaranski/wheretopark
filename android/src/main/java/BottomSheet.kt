package app.wheretopark.android

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.material.MaterialTheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp

@Composable
internal fun BottomSheetHandle(
    modifier: Modifier = Modifier
) {
    Box(
        modifier
            .size(40.dp, 2.dp)
            .background(
                color = MaterialTheme.colors.onSurface.copy(alpha = 0.5f),
                shape = CircleShape
            )
    )
}