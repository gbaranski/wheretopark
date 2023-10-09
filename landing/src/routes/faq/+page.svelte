<script>
	import { mailFromOperator } from '$lib/utils';
</script>

<div class="gap-4 flex flex-col">
	<h1 class="text-2xl font-extrabold">Frequently Asked Questions</h1>
	<div class="px-4">
		<h2 class="text-lg font-bold">When the app will be available?</h2>
		<ul class="list-disc pl-8">
			<li>
				<span class="font-bold">iOS</span> app is available on
				<a
					class="link link-info"
					href="https://apps.apple.com/us/app/where-to-park/id6444453582"
					target="_blank">App Store</a
				> since 21st of November 2022.
			</li>
			<li>
				<span class="font-bold">Web</span> app is available at
				<a href="https://web.wheretopark.app" target="_blank" class="link link-info"> https://web.wheretopark.app </a>.
			</li>
			<li>
				<span class="font-bold">Android</span> app will
				<span class="link link-secondary"> be available soon </span>.
			</li>
		</ul>
	</div>

	<div class="px-4">
		<h2 class="text-lg font-bold">What parking lots do you support</h2>
		<p class="px-4">
			Our AI-based system supports <span class="font-semibold">all kinds</span> of parking lots
		</p>
	</div>

	<div class="px-4">
		<h2 class="text-lg font-bold">What's the accuracy of your system</h2>
		<p class="px-4">
			System accuracy ranges from <span class="font-semibold">90 to 95 percent</span> accuracy in
			determining if parking spot is available.
			<br />
			We're constantly <span class="font-semibold">improving.</span>
		</p>
	</div>

	<div class="px-4">
		<h2 class="text-lg font-bold">
			I'd like to add support for my parking lot to promote it, how can I do that?
		</h2>
		<p class="px-4">
			Contact us at <a href={mailFromOperator} class="link link-info">contact@wheretopark.app</a> with
			some basic details about your parking lot.
		</p>
	</div>
</div>
