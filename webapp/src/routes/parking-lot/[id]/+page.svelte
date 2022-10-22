<script lang="ts">
	import type { ParkingLot } from "$lib/types";

    export let data: {parkingLot: ParkingLot};
</script>

<p>{data.parkingLot.metadata.name}</p>
<p>{data.parkingLot.metadata.address}</p>