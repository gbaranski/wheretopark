<script lang="ts">
	import { SpotType, type ParkingLot } from "$lib/types";
	import { parkingLotStatus } from "$lib/utils";


    export let parkingLot: ParkingLot;
    const status = parkingLotStatus(parkingLot, SpotType.CAR)[0];
</script>
<div class="marker ${status}">
    <i class="material-icons">room</i>
    Elo
</div>

<style>
    .marker {
        display: flex;
        color: red;

    }
</style>