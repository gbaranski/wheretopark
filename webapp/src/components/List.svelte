<script lang="ts">
	import { parkingLots as parkingLotsStore } from "$lib/store";
	import { SpotType } from "$lib/types";
    import { Card, Title, Text } from '@svelteuidev/core';
    $: parkingLots = $parkingLotsStore;
</script>

{#each Object.entries(parkingLots) as [id, parkingLot]}
    <a href={`/parking-lot/${id}`} style="text-decoration: none;">
        <Card class="tile" p="sm" override={{backgroundColor: 'inherit'}}>
            <Title size={24}>{parkingLot.metadata.name}</Title>
            <Text color="gray" size={14} weight={"semibold"}>{parkingLot.state.availableSpots[SpotType[SpotType.CAR]]} available car spots</Text>
        </Card>
    </a>
{/each}


<style>
    a {
        background-color: inherit;
    }

    a:hover {
        background-color: rgb(245, 243, 236);
    }
</style>