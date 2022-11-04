<script lang="ts">
	import { currentMap } from "$lib/store";
	import { onMount } from "svelte";
	import List from "../components/List.svelte";
	import type { LayoutData } from "../types/layout";
	
    export let data: LayoutData;

    currentMap.subscribe((map) => {
        if (map == null) return;
        const [longitude, latitude] = [18.64, 54.35];
		map.flyTo({
			center: [longitude, latitude],
			zoom: 10
		});
    });
</script>

<List parkingLots={data.parkingLots}/>