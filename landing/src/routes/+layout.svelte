<script>
	import '../app.css';
	import Logo from '$lib/Logo.svelte';
	import NavbarMenu from '$lib/NavbarMenu.svelte';
	import { serializeSchema } from '$lib/utils';
	import { frequentlyAskedQuestions, gregoryBaranski, organisationSchema, softwareApplication } from '$lib/schemas';
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

	{@html serializeSchema(softwareApplication)}
	{@html serializeSchema(frequentlyAskedQuestions)}
	{@html serializeSchema(organisationSchema)}
	{@html serializeSchema(gregoryBaranski)}
</svelte:head>

<div class="navbar absolute z-30 bg-base-100">
	<div class="navbar-start">
		<div class="max-md:hidden">
			<Logo />
		</div>
	</div>
	<div class="navbar-center">
		<div class="md:hidden">
			<Logo />
		</div>
		<ul class="menu menu-horizontal px-1 hidden md:flex">
			<NavbarMenu />
		</ul>
	</div>
	<div class="navbar-end">
		<select class="select max-md:hidden">
			<option selected>🇵🇱</option>
			<option>🇬🇧</option>
		</select>
		<div class="dropdown dropdown-end">
			<button tabindex="0" class="btn btn-ghost md:hidden">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-5 w-5"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
					><path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M4 6h16M4 12h8m-8 6h16"
					/></svg
				>
			</button>
			<ul
				class="menu menu-md dropdown-content mt-3 fixed z-20 p-2 shadow bg-base-100 rounded-box w-52"
			>
				<NavbarMenu languageSwitcher />
			</ul>
		</div>
	</div>
</div>

<div class="pt-[4.5rem] w-screen h-screen flex flex-col lg:flex-row">
	<slot />
</div>
