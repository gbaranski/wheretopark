<script lang="ts">
	import type { ID, ParkingLot } from "$lib/types";
    import { Card } from '@svelteuidev/core';

    export let parkingLots: Record<ID, ParkingLot>;
</script>

{#each Object.entries(parkingLots) as [id, parkingLot]}
    <a href={`/parking-lot/${id}`} style="text-decoration: none;">
        <Card p="sm">
            <h4>{parkingLot.metadata.name}</h4>
            <p>{parkingLot.metadata.address}</p>
        </Card>
    </a>
{/each}

