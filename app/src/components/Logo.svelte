<script>
	import { isLoading } from '$/lib/store';
	import logo from '../../../svelte/logo.png';
</script>

<a href="/" class="btn btn-ghost normal-case text-xl">
	<span class="font-mono font-bold text-xl">
		<span class="text-secondary">where</span>
		<span class="text-secondary">park</span>
	</span>
</a>

<a href="/" class="btn btn-ghost normal-case text-xl h-16">
	<img src={logo} alt="where to park icon" class="h-12 m-1" />
	<span class="logo font-bold text-2xl">
		<span class="text-primary"> where </span>

		{#if $isLoading}
			<span class="loading loading-bars loading-md text-primary align-top"></span>
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
