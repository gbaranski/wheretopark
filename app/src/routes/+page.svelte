<script lang="ts">
	import Legend from '$lib/components/Legend.svelte';
import SearchBox from '$lib/components/SearchBox.svelte';
	import { allFeatures } from '$lib/parkingLot';
	import { searchFilters } from '$lib/store';
	import { capitalizeFirstLetter } from '$lib/utils';
</script>

<svelte:head>
	<title>Where To Park</title>
	<meta
		name="description"
		content="Where To Park App with real-time availability of parking lots in Poland and Scotland. Provided with opening hours, prices and other useful information."
	/>
	<meta
		name="keywords"
		content="Parking Lot, Smart City, Gdańsk, Gdynia, Sopot, Warsaw, Warszawa, Poznań, Glasgow"
	/>
</svelte:head>

<div>
	<SearchBox />
	<div class="divider"></div>
	<h1 class="font-mono text-lg font-bold">Search Filters</h1>
	<label class="label cursor-pointer">
		<span class="label-text font-mono">Open right now</span>
		<input type="checkbox" class="checkbox" bind:checked={$searchFilters.openNow} />
	</label>
	{#each allFeatures as feature}
		<label class="label cursor-pointer">
			<span class="label-text font-mono">{capitalizeFirstLetter(feature.toLowerCase())}</span>
			<input type="checkbox" class="checkbox" bind:checked={$searchFilters.hasFeatures[feature]} />
		</label>
	{/each}
	<label class="label cursor-pointer">
		<span class="label-text font-mono">Min. available spaces</span>
		<input
			type="text"
			placeholder="e.g 20"
			class="input input-bordered w-1/3 max-w-xs font-mono text-sm"
			bind:value={$searchFilters.minAvailableSpots}
		/>
	</label>
	<div class="my-8">
		<Legend />
	</div>
</div>