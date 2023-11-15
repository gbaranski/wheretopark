<script lang="ts">
	import { type StatusRating, statusRatingFillColors } from "$lib/parkingLot";

    export let rating: `${StatusRating}`;
</script>

<div class="flex justify-center flex-col">
	<svg
		xmlns="http://www.w3.org/2000/svg"
		viewBox="0 0 24 24"
		stroke-width="1.5"
		stroke="currentColor"
		fill="none"
		class="w-8 h-8 {statusRatingFillColors[rating]}"
	>
		<path
			stroke-linecap="round"
			stroke-linejoin="round"
			d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1115 0z"
		/>
	</svg>
</div>
