<script lang="ts">
	import SearchBox from '$components/SearchBox.svelte';
	import { allFeaturesString } from '$lib/parkingLot';
	import { capitalizeFirstLetter } from '$lib/utils';
</script>

<svelte:head>
	<title>wheretopark.app</title>
	<meta
		name="description"
		content="List of parking lots in Gdańsk, Sopot and Gdynia, with prices, opening hours and it's availability of parking spots."
	/>
	<meta name="keywords" content="Parking Lot, Smart City, Gdańsk, Gdynia, Sopot, Tricity" />
</svelte:head>

<div class="absolute w-full z-10 mx-auto px-5 bg-base-100 rounded-2xl">
	<SearchBox />
	<div class="divider"></div>
	<h1 class="font-mono text-lg font-bold">Search Filters</h1>
	<label class="label cursor-pointer">
		<span class="label-text font-mono">Open right now</span>
		<input type="checkbox" class="checkbox" />
	</label>
	{#each allFeaturesString as feature}
		<label class="label cursor-pointer">
			<span class="label-text font-mono">{capitalizeFirstLetter(feature.toLowerCase())}</span>
			<input type="checkbox" class="checkbox" />
		</label>
	{/each}
	<label class="label cursor-pointer">
		<span class="label-text font-mono">Min. available spaces</span>
		<input
			type="text"
			placeholder="e.g 20"
			class="input input-bordered w-1/3 max-w-xs font-mono text-sm"
		/>
	</label>
</div>
