<script lang="ts">
    import { Text } from '@svelteuidev/core';
    import { SvelteUIProvider } from '@svelteuidev/core';
    import Map from '$components/Map.svelte'
</script>


<svelte:head>
    <link rel="preconnect" href="https://fonts.googleapis.com"> 
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin=""> 
    <link href="https://fonts.googleapis.com/css2?family=Josefin+Sans:wght@600&display=swap" rel="stylesheet">
	<link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
    <script defer data-domain="web.wheretopark.app" src="https://plausible.gbaranski.com/js/plausible.js"></script>
</svelte:head>

<SvelteUIProvider>
    <div class="split master">
        <a href="/">
            <div style="padding: 20px;">
                <Text root='h1' align="center" weight="bold" size={0} override={{fontFamily: "Josefin Sans"}}  >
                    <Text root="span" inherit override={{ color: "#313131", fontSize: 46 }}>where</Text>
                    <Text root="span" inherit override={{ color: "#a28a2b", fontSize: 46 }}>to</Text>
                    <Text root="span" inherit override={{ color: "#313131", fontSize: 46 }}>park</Text>
                </Text>
            </div>
        </a>
        <slot></slot>
    </div>
</SvelteUIProvider>

<div class="split slave">
    <Map/>
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
    width: 350px;
    left: 0;
}

.slave {
    width: calc(100% - 350px);
    right: 0;
}

@media only screen and (orientation:portrait) {
    .split {
        position: static;
        overflow: visible;
    }
    
    :global(#map-container) {
        height: 60%;
    }
    
    .master {
        background-color: rgb(255, 253, 246);
        width: 100%;
        position: absolute;
        top: calc(100% - 40%);
    }

    .slave {
        bottom: 0;
        width: 100%;
        height: 400px;
    }
}
</style>