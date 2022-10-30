<script lang="ts">
	import { SpotType, type ParkingLot, Feature } from "$lib/types";
	import { getCategory, googleMapsLink, resourceIcon, resourceText, timeFromNow } from "$lib/utils";
    import { Title, Text, Divider } from '@svelteuidev/core';

    export let data: {parkingLot: ParkingLot};
    const {metadata, state} = data.parkingLot;
    const features = metadata.features.map((feature) => Feature[feature as keyof typeof Feature]);
    const category = getCategory(features);
</script>

<div class="container">
    <span style="display: flex;">
        <Title root={"span"} size={26} override={{flex: 1}}>{metadata.name}</Title>
        <a style="text-align: right;" href={googleMapsLink(metadata.geometry)} target="_blank" rel="noreferrer">
            <i class="material-icons">directions</i>
        </a>
    </span>
    <Text size={14} weight={"semibold"}>{category}</Text>
    <Divider/>
    
    <div class="field">
        <i class="material-icons">place</i>
        <Text size={14} root="span" >
            {metadata.address}
        </Text>
    </div>

    <div class="field">
        <i class="material-icons">directions_car</i>
        <Text size={14} root="span" >
            {state.availableSpots[SpotType[SpotType.CAR]]} available spots
        </Text>
    </div>
    
    <div class="field">
        <i class="material-icons">schedule</i>
        <Text size={14} root="span" >
            Last updated {timeFromNow(state.lastUpdated)}
        </Text>
    </div>
    
    {#each metadata.resources as resource}
        <div class="field">
            <i class="material-icons">{resourceIcon(resource)}</i>
            <Text size={14} root="span" >
                {resourceText(resource)}
            </Text>
        </div>
    {/each}
</div>
<style>
    .container {
        padding: 20px;
    }
    
    .field {
        margin-bottom: 10px;
        
    }
    
    .field > i {
        font-size: 18px !important;
        margin-right: 10px;
        vertical-align: middle;
    }
    
    div.field > :global(*) {
        font-weight: 500;
        vertical-align: middle;
    }
    
</style>