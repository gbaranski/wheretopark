<script lang="ts">
	import { PUBLIC_MAPBOX_ACCESS_TOKEN } from '$env/static/public';
	import { onMount } from 'svelte';
    import type { SearchBoxSuggestion } from '@mapbox/search-js-core/dist/searchbox/types'
	import { currentMap } from '$lib/store';
	import type { ChangeEventHandler, FormEventHandler } from 'svelte/elements';

	// function after input change, with applied debounce mechanism
	let searchFor: (input: string) => Promise<void>;
    let selectSuggestion: (suggestion: SearchBoxSuggestion) => Promise<void>;
	let suggestions: SearchBoxSuggestion[] = [];

	onMount(async () => {
		const { SearchBoxCore, SessionToken } = await import('@mapbox/search-js-core');
		const search = new SearchBoxCore({
			accessToken: PUBLIC_MAPBOX_ACCESS_TOKEN
		});
		const sessionToken = new SessionToken();
		searchFor = async (input: string): Promise<void> => {
			const response = await search.suggest(input, { sessionToken });
            suggestions = response.suggestions;
		};
        selectSuggestion = async (suggestion: SearchBoxSuggestion): Promise<void> => {
            const { features } = await search.retrieve(suggestion, { sessionToken });
            const coordinates = features.find(feature => feature.geometry.type === 'Point')?.geometry.coordinates!;
            $currentMap?.flyTo({
                center: [coordinates[0], coordinates[1]],
                zoom: 15
            })
            searchTerm = `${suggestion.name}, ${suggestion.place_formatted}`
            suggestions = [];
        }
	});

	let searchTerm: string;
	let timeout: NodeJS.Timeout;

    const onInputChange: FormEventHandler<HTMLInputElement> = (e) => {
        const target = e.target as HTMLInputElement;
        const searchTerm = target.value;
		if (!searchFor) return;
		clearTimeout(timeout);
		timeout = setTimeout(async () => {
            if (searchTerm.length < 3) {
				suggestions = [];
				return;
			}
            
			console.log(`searching for ${searchTerm}`);
			searchFor(searchTerm);
		}, 500);
    }
</script>

<input
	name="address"
	type="text"
	placeholder="🔍  Where'd you park today?"
	class="input input-md input-primary input-bordered w-full text-sm bg-inherit block"
	bind:value={searchTerm}
	on:input={onInputChange}
/>


{#if suggestions.length > 0}
    <ul class="menu bg-base-200 rounded-box mt-5 w-full">
        {#each suggestions as suggestion}
            <li>
                <button class="flex flex-auto flex-col items-start" on:click={() => selectSuggestion(suggestion)}>
                    <p class="font-bold text-md">
                        {suggestion.name}
                    </p>
                    <p class="text-sm">{suggestion.place_formatted}</p>
                </button>
            </li>
        {/each}
    </ul>
{/if}