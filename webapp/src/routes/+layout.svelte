<script lang="ts">
	import type { LayoutData } from "./$types";
    import { Text } from '@svelteuidev/core';
    import Map from '../components/Map.svelte'

    export let data: LayoutData;
	console.log({data});
</script>

<div class="split master">
    <a style="text-decoration: none;" href="/">
        <div style="padding: 20px;">
            <Text size={46} align="center" color="#313131">
                <Text root="span" inherit color="#313131">where</Text>
                <Text root="span" inherit color="#a28a2b">to</Text>
                <Text root="span" inherit>park</Text>
            </Text>
        </div>
    </a>
    <slot></slot>
</div>

<div class="split slave">
    <Map parkingLots={data.parkingLots}/>
</div>

<style>
.split {
    height: 100%;
    position: fixed;
    z-index: 1;
    top: 0;
    bottom: 0;
    overflow-x: hidden;
}

.master {
    background-color: rgb(255, 253, 246);
    width: 400px;
    left: 0;
}

.slave {
    width: calc(100% - 400px);
    right: 0;
}

@media only screen and (max-width: 400px) {
    .master {
        width: 100%;
    }

    .slave {
        width: 0;
        display: none;
    }
}
</style>