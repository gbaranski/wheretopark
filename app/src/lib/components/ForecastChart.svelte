<script lang="ts">
	import type { Forecast } from '$lib/forecaster';
	import type { Action } from 'svelte/action';
	import Chart, { type ChartItem } from 'chart.js/auto';
	import { getNumbersBetween } from '$lib/utils';

	export let forecast: Forecast;

	const labels = getNumbersBetween(0, 24)
		.map((n) => n.toString())
		.map((n) => n.padStart(2, '0'));

	const chart: Action<HTMLElement, Forecast> = (node, forecast: Forecast) => {
		const generateDataset = (forecast: Forecast) => {
			const hourlyPredictions = forecast.predictions.filter(({ date }) => date.getMinutes() == 0);
			const values = getNumbersBetween(0, 24).map(
				(n) => hourlyPredictions.find(({ date }) => date.getHours() == n)?.occupancy || 0
			);
			return {
				label: 'Occupancy',
				data: values,
				fill: false,
				borderColor: 'rgb(75, 192, 192)',
				tension: 0.1
			};
		};

		const chart = new Chart(node as unknown as ChartItem, {
			type: 'line',
			data: {
				labels,
				datasets: [generateDataset(forecast)],
			},
			options: {
				scales: {
					y: {
						beginAtZero: true
					}
				}
			}
		});
		return {
			update(forecast) {
                chart.data = {
                    labels,
                    datasets: [generateDataset(forecast)],
                };
                chart.update();
			},
			destroy() {
				chart.destroy();
			}
		};
	};
</script>

<canvas use:chart={forecast}></canvas>
