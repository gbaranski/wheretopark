<script lang="ts">
	import { SpotType, type ParkingLot, Feature } from "$lib/types";
	import { getCategory } from "$lib/utils";
    import { Title, Text } from '@svelteuidev/core';

    export let data: {parkingLot: ParkingLot};
    const {metadata, state} = data.parkingLot;
    const features = metadata.features.map((feature) => Feature[feature as keyof typeof Feature]);
    const category = getCategory(features);
</script>

<div class="container">
    <Title size={26}>{metadata.name}</Title>
    <Text size={14} weight={"semibold"}>{category}</Text>
    <Text size={14} weight={"semibold"}>{state.availableSpots[SpotType[SpotType.CAR]]} available car spots</Text>


</div>
<style>
    .container {
        padding: 20px;

    }
</style>