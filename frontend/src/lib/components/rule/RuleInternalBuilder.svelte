<script lang="ts">
    import type { RuleInternal } from "@/types/rule";
    import RuleInternalBuilder from "./RuleInternalBuilder.svelte";
    import Sensor from "../sensor/Sensor.svelte";
	import { Button } from '$lib/components/ui/button';

	export let expanded = false;
	export let internal: RuleInternal;
	export let sensors: Sensor[];

	const name = getName();

	function toggle() {
		expanded = !expanded;
	}

	function getName(){
		if (internal.type === "and" || internal.type === "or" || internal.type === "not") {
			return internal.type.toUpperCase();
		} else {
			const sensorName = sensors.find(sensor => sensor.id === internal.sensor_id)?.name;
			return `Value of "${sensorName}" ${internal.type === "lt" ? "lower than" : "greater than"} ${internal.value}`;
		}
	}
</script>

<div>
	<Button on:click={toggle}>{name}</Button>
	{#if expanded}
		{#if internal.type === "and" || internal.type === "or"}
			<ul>
				{#each internal.children as child}
					<li>
						<RuleInternalBuilder internal={child} {sensors}/>
					</li>
				{/each}
			</ul>
		{:else if internal.type === "not"}
			<ul>
				<li>
					<RuleInternalBuilder internal={internal.wrapped} {sensors}/>
				</li>
			</ul>
		{:else}
			<p>{internal.type}</p>
			<p>{internal.sensor_id}</p>
			<p>{internal.value}</p>
		{/if}
	{/if}
</div>

<style>

	ul {
		margin: 0.6em 0 0.8em 0;
		padding: 0em 0 0 2em;
		list-style: none;
		border-left: 1px solid rgba(128, 128, 128, 0.4);
	}

	li {
		padding: 0.2em 0;
	}
</style>
