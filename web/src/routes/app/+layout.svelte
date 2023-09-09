<script lang="ts">
	import Map from '$components/Map.svelte';
</script>

<div class="w-full h-full lg:order-2">
	<Map />
</div>
<div class="p-5 h-1/3 lg:overflow-y-scroll lg:h-full lg:w-7/12 xl:w-5/12 2xl:w-4/12 lg:order-1">
	<slot/>
</div>