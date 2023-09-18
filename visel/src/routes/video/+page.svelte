<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageData } from './$types';
	import type { WebAnnotation } from '@recogito/annotorious';
	import { page } from '$app/stores';
	import { generateCode, type Polygon, parseCode } from '$lib';
	import {v4 as uuidv4} from 'uuid';

	export let data: PageData;

	let polygons: Polygon[] = [];
	let code: string = generateCode([]);
	let error: string | undefined = undefined;
	let updateAnnotations: ((polygons: Polygon[]) => void) | undefined = undefined;
	
	const onCodeUpdate = () => {
		console.log("on code update");
		try {
			const newPolygons = parseCode(code);
			error = undefined;
			polygons = newPolygons;
			if (updateAnnotations) updateAnnotations(polygons);
		} catch (e: any) {
			error = e.toString();
		}
	}

	$: [code] && onCodeUpdate();

	const generateYAML = () => {
		code = generateCode(polygons);
	};

	onMount(async () => {
		import('@recogito/annotorious/dist/annotorious.min.css');
		const Annotorious = await import('@recogito/annotorious');
		const anno = Annotorious.init({
			image: 'parking-lot',
			drawOnSingleClick: true
		});

		updateAnnotations = (polygons: Polygon[]) => {
			console.log("updateAnnotations()")
			anno.clearAnnotations();
			const annotations = polygons.map(
				(poly): WebAnnotation => ({
					'@context': 'http://www.w3.org/ns/anno.jsonld',
					type: 'Annotation',
					body: [
						{
							type: 'TextualBody',
							purpose: 'tagging',
							value: 'ParkingSpot'
						}
					],
					id: uuidv4(),
					target: {
						source: $page.url.toString(),
						selector: {
							type: 'SvgSelector',
							value: `<svg><polygon points="${poly
								.map(({ x, y }) => `${x},${y}`)
								.join(' ')}"></polygon></svg>`
						}
					}
				})
			);
			for (const annotation of annotations) {
				console.log('adding annotation', { annotation });
				anno.addAnnotation(annotation);
			}
		};

		const updatePolygons = () => {
			setTimeout(() => {

			const annotations: WebAnnotation[] = anno.getAnnotations();
			console.log({annotations})
			polygons = annotations.map((a) => {
				const html = a.target.selector.value;
				// take value from inside of the quotes:
				const strPoints = html.match(/"([^']+)"/)?.pop();
				if (!strPoints) throw new Error(`couldn't find points in ${html}`);
				const points = strPoints.split(' ').map((xy) => {
					const [x, y] = xy.split(',');
					return {
						x: parseInt(x),
						y: parseInt(y)
					};
				});
				return points;
			});
			}, 1000)
		};
		anno.on('createSelection', async function (selection: WebAnnotation) {
			selection.body = [
				{
					type: 'TextualBody',
					purpose: 'tagging',
					value: 'ParkingSpot'
				}
			];
			await anno.updateSelected(selection);
			anno.saveSelected();
			updatePolygons();
		});
		anno.on('createAnnotation', updatePolygons);
		anno.on('updateAnnotation', updatePolygons);
		anno.on('deleteAnnotation', updatePolygons);
		anno.setDrawingTool('polygon');
	});

	const imgSrc = `/api/capture?src=${encodeURIComponent(data.src!)}`;
</script>

<img id="parking-lot" src={imgSrc} alt="parking lot" />

{#if error}
	<span style="color: red;">{error}</span>
{/if}
<button on:click={generateYAML}>Generate YAML</button>
<textarea placeholder="YAML code" bind:value={code} id="raw" />

<style>
	#parking-lot {
		background: #fff;
		display: block;
		margin: 20px auto;
		/* zoom: 2; */
	}

	#raw {
		font-family: monospace;
		width: 100%;
		height: 300px;
	}
</style>
