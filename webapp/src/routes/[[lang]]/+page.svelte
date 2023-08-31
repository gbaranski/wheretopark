<script lang="ts">
	import Map from '$components/Map.svelte';
	import { Feature, features } from '$lib/types';
	import { capitalizeFirstLetter } from '$lib/utils';

	const setCurrentLocation = () => {};
</script>

<svelte:head>
	<title>wheretopark.app</title>
	<meta
		name="description"
		content="List of parking lots in Gdańsk, Sopot and Gdynia, with prices, opening hours and it's availability of parking spots."
	/>
	<meta name="keywords" content="Parking Lot, Smart City, Gdańsk, Gdynia, Sopot, Tricity" />
</svelte:head>

<div class="absolute w-full z-10 mx-auto p-5 bg-base-100 rounded-2xl">
	<div class="flex items-center">
		<input
			name="address"
			type="text"
			placeholder="🔍  Where'd you park today?"
			class="input input-md input-primary input-bordered w-11/12 text-sm bg-inherit"
		/>
		<button class="btn btn-md btn-outline btn-secondary ml-2" on:click={setCurrentLocation}>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="1.5"
				stroke="currentColor"
				class="w-6 h-6"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M15 10.5a3 3 0 11-6 0 3 3 0 016 0z"
				/>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1115 0z"
				/>
			</svg>
		</button>
		<dialog id="filtersModal" class="modal">
			<form method="dialog" class="modal-box">
				<h3 class="font-bold text-lg">Hello!</h3>
				<p class="py-4">Press ESC key or click the button below to close</p>
				<div class="modal-action">
					<!-- if there is a button in form, it will close the modal -->
					<button class="btn">Close</button>
				</div>
			</form>
		</dialog>
	</div>
	
	<div class="divider"></div>
	<h1 class="font-mono text-lg font-bold">Search Filters</h1>
	<label class="label cursor-pointer">
		<span class="label-text font-mono">Open right now</span>
		<input type="checkbox" class="checkbox" />
	</label>
	{#each features as feature}
		<label class="label cursor-pointer">
			<span class="label-text font-mono">{capitalizeFirstLetter(feature.toLowerCase())}</span>
			<input type="checkbox" class="checkbox" />
		</label>
	{/each}
	<label class="label cursor-pointer">
		<span class="label-text font-mono">Min. available spaces</span>
		<input type="text" placeholder="e.g 20" class="input input-bordered w-1/3 max-w-xs font-mono text-sm" />
	</label>
</div>
