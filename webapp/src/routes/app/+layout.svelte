<script lang="ts">
	import Map from '$components/Map.svelte';
</script>

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