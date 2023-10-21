<script lang="ts">
	import '../app.css';
	import Logo from '$lib/Logo.svelte';
	import NavbarMenu from '$lib/NavbarMenu.svelte';
	import { serializeSchema } from '$lib/utils';
	import {
		frequentlyAskedQuestions,
		gregoryBaranski,
		organisationSchema,
		softwareApplication
	} from '$lib/schemas';
</script>

<svelte:head>
	<script
		data-goatcounter="https://wheretopark.goatcounter.com/count"
		async
		src="https://gc.zgo.at/count.js"
	></script>
	<title>Where To Park</title>
	<meta
		name="description"
		content="List of parking lots in Gdańsk, Sopot and Gdynia, with prices, opening hours and it's availability of parking spots."
	/>
	<meta name="keywords" content="Parking Lot, Smart City, Gdańsk, Gdynia, Sopot, Tricity" />

	{@html serializeSchema(softwareApplication)}
	{@html serializeSchema(frequentlyAskedQuestions)}
	{@html serializeSchema(organisationSchema)}
	{@html serializeSchema(gregoryBaranski)}
</svelte:head>

<navbar class="navbar z-30 px-8 md:px-32 py-8">
	<div class="navbar-start">
		<Logo />
	</div>
	<div class="navbar-end">
		<NavbarMenu />
	</div>
</navbar>

<div class="pt-16 mb-48 px-8 md:px-32 prose">
	<slot />
</div>
