<script lang="ts">
    import type { RuleInternal } from "@/types/rule";
    import RuleInternalBuilder from "./RuleInternalBuilder.svelte";
	import { Button } from '$lib/components/ui/button';
    import type { Sensor } from "@/types/sensor";
    import ComparisonRule from "./ComparisonRule.svelte";
	import { Symbol } from 'radix-icons-svelte';

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

	/** TODO:
	 * śmietniczek po prawej stronie do usuwania
	 *   - not - usuwa tylko nota
	 *   - and/or - usuwa wszystko z zawartością
	 *   - lg/lt - usuwa pojedynczą regułe 
	 
	 * zmiana typu reguły 
	 *   - OR <-> AND
	 
	 * dodawanie grupy nowych reguł (OR / AND)
	 * dodawanie pojedynczej reguły (LG / LT)
	 * zaprzeczanie reguły po przez wykrzyknik koło śmietniczka
	 * */ 

</script>

<div class="w-full">
	{#if internal.type === "lt" || internal.type === "gt"}
		<ComparisonRule internal={internal} {sensors}/>
	{:else}
		{#if internal.type === "and" || internal.type === "or"}
			<div class="flex">
				<Button on:click={toggle}>{internal.type.toUpperCase()}</Button>
				<Button on:click={() => {
					internal.type = internal.type === "and" ? "or" : "and";
				}}>
					<Symbol />
				</Button>
			</div>
			{#if expanded}
			<ul>
				{#each internal.children as child}
				<li>
					<RuleInternalBuilder internal={child} {sensors}/>
				</li>
				{/each}
			</ul>
			{/if}
		{:else}
			<Button on:click={toggle}>{"NOT"}</Button>
			{#if expanded}
				<ul>
					<li>
						<RuleInternalBuilder internal={internal.wrapped} {sensors}/>
					</li>
				</ul>
			{/if}
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
