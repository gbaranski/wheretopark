<script>
	import '../app.css';
	import "@fontsource-variable/inter";
	import Logo from '$lib/components/Logo.svelte';
	import Map from '$lib/components/Map.svelte';
	import { App } from '@capacitor/app';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { Capacitor } from '@capacitor/core';

	onMount(() => {
		if (Capacitor.isNativePlatform()) {
			App.addListener('appUrlOpen', (data) => {
				goto(data.url);
			});
		}
	});
</script>

<svelte:head>
	<script
		data-goatcounter="https://wheretopark.goatcounter.com/count"
		async
		src="https://gc.zgo.at/count.js"
	></script>
</svelte:head>

<div class="h-screen pt-safe">
	<div class="absolute z-30 navbar bg-base-100 lg:hidden justify-center h-16">
		<Logo />
	</div>

	<div class="h-full flex flex-col lg:flex-row max-lg:pt-16">
		<div class="w-full h-full lg:order-2">
			<Map />
		</div>
		<div
			class="px-5 pt-3 h-1/3 lg:overflow-y-scroll lg:h-full lg:w-7/12 xl:w-5/12 2xl:w-4/12 lg:order-1"
		>
			<div class="pb-3 text-center max-lg:hidden">
				<Logo />
			</div>
			<slot />
			<div class="pb-safe" />
		</div>
	</div>
</div>

<style>
	.pt-safe {
		margin-top: env(safe-area-inset-top);
	}
	.pb-safe {
		padding-bottom: env(safe-area-inset-bottom);
	}
</style>
