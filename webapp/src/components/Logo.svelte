<script lang="ts">
	import { isLoading } from '$lib/store';

	export let loadable: boolean;
</script>

<a href="/" class="btn btn-ghost normal-case text-xl">
	<span class="font-mono font-bold text-xl">
		<span class="text-secondary">where</span>
		{#if loadable && $isLoading}
			<span class="loading loading-bars loading-md text-primary align-top"></span>
		{:else}
			<span class="text-primary">to</span>
		{/if}
		<span class="text-secondary">park</span>
	</span>
</a>
