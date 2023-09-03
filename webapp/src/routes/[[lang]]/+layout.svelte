<script lang="ts">
	import Logo from '$components/Logo.svelte';
	import Map from '$components/Map.svelte';
	import NavbarMenu from '$components/NavbarMenu.svelte';
	import '../../app.css';
</script>

<svelte:head>
	<link rel="preconnect" href="https://fonts.googleapis.com" />
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="" />
	<link
		href="https://fonts.googleapis.com/css2?family=Josefin+Sans:wght@600&display=swap"
		rel="stylesheet"
	/>
	<link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet" />
	<script
		defer
		data-domain="web.wheretopark.app"
		src="https://plausible.gbaranski.com/js/plausible.js"
	></script>
</svelte:head>

<div class="navbar absolute z-30 bg-base-100 ">
	<div class="navbar-start">
		<div class="max-md:hidden">
			<Logo/>
		</div>
	</div>
	<div class="navbar-center">
		<div class="md:hidden">
			<Logo/>
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

<div class="w-screen h-screen">
	<div class="flex flex-col h-full w-full lg:flex-row">
		<div class="">
			<Map />
		</div>
		<div class="p-5 mt-[75vh] lg:overflow-y-scroll lg:mt-16 lg:w-5/12 xl:w-4/12 2xl:w-3/12">
			<slot />
		</div>
	</div>
</div>