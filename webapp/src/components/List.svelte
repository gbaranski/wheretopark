<script lang="ts">
	import { parkingLots as parkingLotsStore } from "$lib/store";
	import { SpotType } from "$lib/types";
    $: parkingLots = $parkingLotsStore;
</script>

<ul class="menu menu-xs bg-base-200 w-56 rounded-box">
    {#each Object.entries(parkingLots) as [id, parkingLot]}
        <li>
            <a href={`/parking-lot/${id}`} style="text-decoration: none;">
                <h1>{parkingLot.metadata.name}</h1>
                <h6>{parkingLot.state.availableSpots[SpotType[SpotType.CAR]]} available car spots</h6>
            </a>
        </li>
    {/each}
</ul>


<style>
    a {
        background-color: inherit;
    }

    a:hover {
        background-color: rgb(245, 243, 236);
    }
</style>