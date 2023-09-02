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

<div class="navbar absolute z-30 bg-base-100">
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


<div>
	<div class="split master">
		<!-- <p class="text-xs pt-5 px-5" >This webapp is currently under heavy development 🏗 . You might find a lot of bugs or even dragons</p> -->
		<slot />
	</div>
	<div class="split slave">
		<Map />
	</div>
</div>

<style>
	.split {
		height: 100%;
		position: fixed;
		z-index: 1;
		top: 0;
		bottom: 0;
		overflow-x: hidden;
	}
	
	.master {
		width: 450px;
		left: 0;
		top: 4em;
	}
	
	.slave {
		width: calc(100% - 450px);
		right: 0;
	}
	
	@media only screen and (orientation:portrait) {
		.split {
			position: static;
			overflow: visible;
		}
		
		:global(#map-container) {
			height: 70%;
		}
		
		.master {
			width: 100%;
			position: absolute;
			top: calc(100% - 20%);
		}
	
		.slave {
			bottom: 0;
			width: 100%;
			height: 70%;
		}
	}
	</style>