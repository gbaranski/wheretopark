<script>
	import { isLoading } from '$lib/store';
</script>

<a href="/" class="btn btn-ghost normal-case text-xl h-16" data-sveltekit-noscroll>
	<img src="/icon@64.webp" alt="where to park icon" class="m-1" width="40px" height="40px" />
	<span class="logo font-bold text-2xl">
		<span class="text-primary"> where </span>

		{#if $isLoading}
			<span class="loading loading-bars loading-md text-success align-top"></span>
		{:else}
			<span class="text-success text-opacity-80"> to </span>
		{/if}
		<span class="text-error text-opacity-80"> park </span>
	</span>
</a>

<style>
	@import url('https://fonts.googleapis.com/css2?family=Josefin+Sans:wght@700&display=swap');

	.logo {
		font-family: 'Josefin Sans', sans-serif;
	}
</style>
