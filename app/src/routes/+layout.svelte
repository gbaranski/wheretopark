<script>
	import '../app.css';
	import Logo from '$lib/components/Logo.svelte';
	import Map from '$lib/components/Map.svelte';
</script>

<svelte:head>
	<script
		defer
		data-domain="web.wheretopark.app"
		src="https://plausible.gbaranski.com/js/plausible.js"
	></script>
	<title>Where To Park</title>
	<meta
		name="description"
		content="List of parking lots in Gdańsk, Sopot and Gdynia, with prices, opening hours and it's availability of parking spots."
	/>
	<meta name="keywords" content="Parking Lot, Smart City, Gdańsk, Gdynia, Sopot, Tricity" />
</svelte:head>

<div class="w-screen h-screen flex flex-col lg:flex-row">
	<div class="w-full h-full lg:order-2">
		<Map />
	</div>
	<div
		class="px-5 py-3 h-1/3 lg:overflow-y-scroll lg:h-full lg:w-7/12 xl:w-5/12 2xl:w-4/12 lg:order-1"
	>
		<div class="pb-3 text-center">
			<Logo />
		</div>
		<slot />
	</div>
</div>
