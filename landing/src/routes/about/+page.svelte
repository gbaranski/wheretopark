<script lang="ts">
	import { gregorybaranski, stanbaranski } from '$lib/assets/team';
</script>

<div class="p-16 gap-6 flex flex-col">
	<h1 class="text-2xl font-bold">General Contact</h1>
	<p>
		Email: <a class="font-bold link link-info" href="mailto:contact@wheretopark.app"
			>contact@wheretopark.app</a
		>
	</p>
	<h1 class="text-2xl font-bold">Our Team</h1>
	<div class="flex mt-1 gap-16">
		<div class="border-2">
			<img
				src={gregorybaranski}
				alt="Gregory Baranski"
				width="320px"
				height="320px"
				class="object-contain h-[320px] w-[320px]"
			/>
			<div class="bg-slate-200 p-4 gap-4 flex flex-col h-48 space-between">
				<a href="https://gbaranski.com" target="_blank" class="link link-secondary"
					>Gregory Baranski</a
				>
				<div>
					<p>Founder @ Where To Park</p>
					<p>Software Developer</p>
					<p>Business Relations</p>
				</div>
				<a class="link link-secondary" href="mailto:gregory.baranski@wheretopark.app"
					>gregory.baranski@wheretopark.app</a
				>
			</div>
		</div>

		<div class="border-2">
			<img
				src={stanbaranski}
				alt="Stan Baranski"
				width="320px"
				height="320px"
				class="object-scale-cover h-[320px] w-[320px]"
			/>
			<div class="bg-slate-200 p-4 gap-4 flex flex-col h-48 justify-between">
				<a href="https://stan.bar" target="_blank" class="link link-secondary">Stan Baranski</a>
				<div>
					<p>Legal tutor</p>
					<p>PhD Student</p>
					<p>Legacy technology assistance</p>
				</div>
				<a class="link link-secondary" href="mailto:stan.baranski@wheretopark.app"
					>stan.baranski@wheretopark.app</a
				>
			</div>
		</div>
	</div>
</div>
