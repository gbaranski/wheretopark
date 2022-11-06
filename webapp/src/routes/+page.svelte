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

<svelte:head>
    <title>wheretopark.app</title>
	<meta name="description" content="List of parking lots in Gdańsk, Sopot and Gdynia, with prices, opening hours and it's availability of parking spots."/>
	<meta name="keywords" content="Parking Lot, Smart City, Gdańsk, Gdynia, Sopot, Tricity"/>
</svelte:head>

<List parkingLots={data.parkingLots}/>