<script lang="ts">
	import Map from '$components/Map.svelte';

	const openFiltersModal = () => {
		// @ts-ignore
		window['filtersModal'].showModal();
	};
</script>

<svelte:head>
	<title>wheretopark.app</title>
	<meta
		name="description"
		content="List of parking lots in Gdańsk, Sopot and Gdynia, with prices, opening hours and it's availability of parking spots."
	/>
	<meta name="keywords" content="Parking Lot, Smart City, Gdańsk, Gdynia, Sopot, Tricity" />
</svelte:head>

<div
	class="absolute max-md:bottom-10 md:top-20 w-11/12 md:w-3/5 z-10 inset-x-0 mx-auto md:mx-5 p-5 bg-base-100 rounded-2xl"
>
	<div class="flex items-center">
		<input
			name="address"
			type="text"
			placeholder="🔍  Where'd you park today?"
			class="input input-bordered w-11/12 h-9 text-sm bg-inherit"
		/>
		<button class="btn btn-sm btn-outline ml-2" on:click={openFiltersModal}>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="1.5"
				stroke="currentColor"
				class="w-4 h-4"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M10.5 6h9.75M10.5 6a1.5 1.5 0 11-3 0m3 0a1.5 1.5 0 10-3 0M3.75 6H7.5m3 12h9.75m-9.75 0a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m-3.75 0H7.5m9-6h3.75m-3.75 0a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m-9.75 0h9.75"
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
</div>

<Map />
