<script lang="ts">
    import type { RuleAndType, RuleDetails, RuleInternal, RuleNotType, RuleOrType } from "@/types/rule";
    import RuleInternalBuilder from "./RuleInternalBuilder.svelte";
	import { Button } from '$lib/components/ui/button';
    import type { Sensor } from "@/types/sensor";
    import ComparisonRule from "./ComparisonRule.svelte";
	import { Symbol } from 'radix-icons-svelte';
    import { Trash, Plus, Slash} from "svelte-radix";
    import NewRule from "./NewRule.svelte";

	type Parent = RuleDetails | RuleNotType | RuleAndType | RuleOrType

	export let expanded = false;
	export let internal: RuleInternal;
	export let parent: Parent;
	export let secondParent: Parent | undefined;
	export let sensors: Sensor[];
	
	let adding = false;

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
			console.log("You can't remove first rule")
			return;
		}

		if (internal.type == "not"){
			if (isRuleDetails(parent)){
				parent.internal = internal.wrapped;
				return;
			}

			else if (parent.type === "or" || parent.type === "and") {
				parent.children = parent.children.filter((child) => {
					return child != internal
				})
				parent.children.push(internal.wrapped)
				return;
			}

			else if (parent.type === "not") {
				parent.wrapped = internal.wrapped;
				return;
			}
		}

		if (parent.type === "or" || parent.type === "and") {
			parent.children = parent.children.filter((child) => {
				return child != internal
			})
		}

		if (parent.type === "not") {
			if (!secondParent) return;
			
			else if (isRuleDetails(secondParent)){
				secondParent.internal = internal;
				return;
			}

			else if (secondParent.type === "or" || secondParent.type === "and") {
				secondParent.children = secondParent.children.filter((child) => {
					return child != internal
				})
				secondParent.children.push(internal)
				return;
			}

			else if (secondParent.type === "not") {
				secondParent.wrapped = internal;
				return;
			}
		}

		// TODO: add removing not rule
	}

	function addRule(){
		adding = true;
	}

	function negateRule(){
		if (isRuleDetails(parent)){
			parent.internal = {
				type: "not",
				wrapped: parent.internal
			}
			return;
		}

		if (parent.type === "or" || parent.type === "and") {
			parent.children = parent.children.filter((child) => {
				return child != internal
			})
			parent.children.push({
				type: "not",
				wrapped: internal
			})
		}
	}
	
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
		
		{#if internal.type != "not"} 
			<Button on:click={negateRule}>
				<Slash />
			</Button>
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
						<RuleInternalBuilder bind:internal={child} bind:parent={internal} bind:secondParent={parent} {sensors} />
					</li>
				{/each}

				{#if adding}
					<NewRule bind:open={adding} {sensors} bind:parent={internal}/>
				{:else}
					<li>
						<Button on:click={addRule}>
							<Plus />
						</Button>
					</li>
				{/if}
			</ul>
		{:else if internal.type === "not"}
			<ul>
				<li>
					<RuleInternalBuilder bind:internal={internal.wrapped} bind:parent={internal} bind:secondParent={parent} {sensors}/>
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
