<script lang="ts">
    import type { RuleAndType, RuleDetails, RuleGtType, RuleInternal, RuleLtType, RuleNotType, RuleOrType } from "@/types/rule";
    import RuleInternalBuilder from "./RuleInternalBuilder.svelte";
	import { Button } from '$lib/components/ui/button';
    import type { Sensor } from "@/types/sensor";
    import ComparisonRule from "./ComparisonRule.svelte";
	import { Symbol } from 'radix-icons-svelte';
    import { Trash } from "svelte-radix";

	export let expanded = false;
	export let internal: RuleInternal;
	export let parent: RuleDetails | RuleNotType | RuleAndType | RuleOrType
	export let sensors: Sensor[];

	function toggle() {
		expanded = !expanded;
	}

	$: if ((internal.type == "and" || internal.type == "or") && internal.children){
		console.log('children changed', internal.children)
	}

	function isRuleDetails(parentInput: RuleInternal | RuleDetails): parentInput is RuleDetails{
		return Object.hasOwn(parentInput, "description") 
	}

	function deleteRule(){
		if (isRuleDetails(parent)){
			// todo errror
			return {
				isError: true,
				data: "You cannot remove first rule"
			}
		}

		if (parent.type === "or" || parent.type === "and") {
			parent.children = parent.children.filter((child) => {
				return child != internal
			})
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
	<div class="flex">
		{#if internal.type === "lt" || internal.type === "gt"}
			<ComparisonRule {internal} {sensors}/>
		{:else if internal.type === "and" || internal.type === "or"}
			<div class="flex">
				<Button on:click={toggle}>{internal.type.toUpperCase()}</Button>
				<Button on:click={() => {
					internal.type = internal.type === "and" ? "or" : "and";
				}}>
					<Symbol />
				</Button>
			</div>
		{:else}
			<Button on:click={toggle}>{"NOT"}</Button>
		{/if}
		<Button on:click={deleteRule}>
			<Trash />
		</Button>
	</div>

	{#if expanded}
		{#if internal.type === "and" || internal.type === "or"}
			<ul>
				{#each internal.children as child}
					<li>
						<RuleInternalBuilder bind:internal={child} {sensors} bind:parent={internal}/>
					</li>
				{/each}
			</ul>
		{:else if internal.type === "not"}
			<ul>
				<li>
					<RuleInternalBuilder bind:internal={internal.wrapped} {sensors} bind:parent={internal}/>
				</li>
			</ul>
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
