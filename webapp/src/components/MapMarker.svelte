<script lang="ts">
	import { SpotType, type ParkingLot } from "$lib/types";
	import { markerColor, parkingLotStatus } from "$lib/utils";

    export let parkingLot: ParkingLot;
    const status = parkingLotStatus(parkingLot, SpotType.CAR)[0];
    const color = markerColor(parkingLot.state.availableSpots["CAR"], parkingLot.metadata.totalSpots["CAR"], status);
</script>

<div class="flex flex-auto justify-center flex-col">
    <svg xmlns="http://www.w3.org/2000/svg" fill={color} viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-8 h-8">
        <path stroke-linecap="round" stroke-linejoin="round" d="M15 10.5a3 3 0 11-6 0 3 3 0 016 0z" />
        <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1115 0z" />
    </svg>
    <!-- <span class="text-center text-xs font-light">
        {parkingLot.metadata.name}
    </span> -->
</div>