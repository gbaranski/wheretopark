<script lang="ts">
	import Map from '$components/Map.svelte';
	import NavbarMenu from '$components/NavbarMenu.svelte';
	import '../../app.css';

	const openFiltersModal = () => {
		// @ts-ignore
		window['filtersModal'].showModal();
	};
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

<!-- <div class="split master">
        <a href="/">
            <div class="text-center p-5">
            </div>
        </a>
        <slot></slot>
    </div> -->

<div class="">
	<div class="navbar absolute z-20 bg-base-100">
		<div class="navbar-start">
			<a href="/" class="btn btn-ghost normal-case text-xl">
				<span class="font-mono">
					<span class="text-xl text-slate-300 font-bold">where</span>
					<span class="text-xl text-yellow-600 font-bold">to</span>
					<span class="text-xl text-slate-300 font-bold">park</span>
				</span>
			</a>
		</div>
		<div class="navbar-center hidden md:flex">
			<ul class="menu menu-horizontal px-1">
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
	<div
		class="absolute top-20 w-11/12 md:w-3/5 z-10 inset-x-0 mx-auto md:mx-5 p-5 bg-primary-content rounded-2xl"
	>
		<div class="flex items-center">
			<input
				name="address"
				type="text"
				placeholder="🔍  Where'd you park today?"
				class="input input-bordered w-11/12 h-9 text-sm bg-primary-content"
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
</div>

<style>
	a,
	a:hover,
	a:visited,
	a:active {
		color: inherit;
		text-decoration: none;
	}
</style>
