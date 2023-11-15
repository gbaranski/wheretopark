<script lang="ts">
	import { StatusRating } from '$lib/parkingLot';
	import MapMarker from './MapMarker.svelte';
	const elements: [StatusRating, string][] = [
		[StatusRating.Good, '<b>More than 10%</b> available spaces and open'],
		[StatusRating.Fair, '<b>Less than 10%</b> available spaces or <b>closes soon</b>'],
		[StatusRating.Poor, ' <b>No available</b> spaces or is <b>closed</b>'],
		[StatusRating.Inactive, '<b>No data</b> available since >1 hour']
	];
</script>

<div>
	{#each elements as [key, value]}
		<div class="flex flex-row items-center">
			<MapMarker rating={key} />
			<span class="text-xs text-muted">{@html value}</span>
		</div>
	{/each}
</div>
