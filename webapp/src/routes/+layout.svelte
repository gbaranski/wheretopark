<script lang="ts">
	import type { LayoutData } from "./$types";
    import { Text, Title } from '@svelteuidev/core';
    import Map from '../components/Map.svelte'

    export let data: LayoutData;
</script>

<svelte:head>
    <link rel="preconnect" href="https://fonts.googleapis.com"> 
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin=""> 
    <link href="https://fonts.googleapis.com/css2?family=Josefin+Sans:wght@600&display=swap" rel="stylesheet">
</svelte:head>
<div class="split master">
    <a href="/">
        <div style="padding: 20px;">
            <Title align="center" override={{fontFamily: "Josefin Sans", fontWeight: 600, fontSize: 0 }}>
                <Text root="span" inherit override={{ color: "#313131", fontSize: 46 }}>where</Text>
                <Text root="span" inherit override={{ color: "#a28a2b", fontSize: 46 }}>to</Text>
                <Text root="span" inherit override={{ color: "#313131", fontSize: 46 }}>park</Text>
            </Title>
        </div>
    </a>
    <slot></slot>
</div>

<div class="split slave">
    <Map parkingLots={data.parkingLots}/>
</div>

<style>
a, a:hover, a:visited, a:active {
    color: inherit;
    text-decoration: none;
}

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