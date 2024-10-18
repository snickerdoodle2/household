<script lang="ts">
    import type {
    NewRule,
        RuleAndType,
        RuleDetails,
        RuleInternal,
        RuleNotType,
        RuleOrType,
    } from '@/types/rule';
    import RuleInternalBuilder from './RuleInternalBuilder.svelte';
    import { Button } from '$lib/components/ui/button';
    import type { Sensor } from '@/types/sensor';
    import ComparisonRule from './ComparisonCondition.svelte';
    import { Symbol } from 'radix-icons-svelte';
    import { Trash, Plus, Slash } from 'svelte-radix';
    import ConditionBuilder from './ConditionBuilder.svelte';

    type Parent = RuleDetails | NewRule | RuleNotType | RuleAndType | RuleOrType;

    export let expanded = false;
    export let internal: RuleInternal | {};
    export let parent: Parent;
    export let secondParent: Parent | undefined;
    export let sensors: Sensor[];
    export let editingDisabled: boolean = false;
    
    let adding = false;

    $: background = isRule(internal) && (internal.type === 'lt' || internal.type === 'gt') ? '' : 'bg-foreground';
    let isFirstRule = isRootRule(parent);

    function toggleExpand() {
        expanded = !expanded;
    }

    function isRootRule(
        parentInput: RuleInternal | RuleDetails | NewRule
    ): parentInput is RuleDetails | NewRule {
        return Object.hasOwn(parentInput, 'description');
    }

    function isRule(internal: RuleInternal | {}): internal is RuleInternal {
        return Object.keys(internal).length !== 0;
    }

    function deleteRule() {
        if(!isRule(internal)) return;

        if (isRootRule(parent)) {
            internal = {};
            return;
        }

        if (internal.type == 'not') {
            if (isRootRule(parent)) {
                parent.internal = internal.wrapped;
                return;
            } else if (parent.type === 'or' || parent.type === 'and') {
                parent.children = parent.children.filter((child) => {
                    return child != internal;
                });
                parent.children.push(internal.wrapped);
                return;
            } else if (parent.type === 'not') {
                parent.wrapped = internal.wrapped;
                return;
            }
        }

        if (parent.type === 'or' || parent.type === 'and') {
            parent.children = parent.children.filter((child) => {
                return child != internal;
            });
        }

        if (parent.type === 'not') {
            if (!secondParent) return;
            else if (isRootRule(secondParent)) {
                secondParent.internal = internal;
                return;
            } else if (
                secondParent.type === 'or' ||
                secondParent.type === 'and'
            ) {
                secondParent.children = secondParent.children.filter(
                    (child) => {
                        return child != internal;
                    }
                );
                secondParent.children.push(internal);
                return;
            } else if (secondParent.type === 'not') {
                secondParent.wrapped = internal;
                return;
            }
        }
    }

    function addRule() {
        adding = true;
    }

    function negateRule() {
        if(!isRule(internal)) return;

        if (isRootRule(parent)) {
            parent.internal = {
                type: 'not',
                wrapped: parent.internal,
            };
            return;
        }

        if (parent.type === 'or' || parent.type === 'and') {
            parent.children = parent.children.filter((child) => {
                return child != internal;
            });
            parent.children.push({
                type: 'not',
                wrapped: internal,
            });
        }
    }
</script>

<div class="w-full min-w-[35rem]">
    {#if isRule(internal)}
        <!-- Main view (AND, OR, ...) -->
        <div class="flex inline-flex {background} rounded">
            {#if internal.type === 'lt' || internal.type === 'gt'}
                <ComparisonRule {internal} {sensors} bind:editingDisabled={editingDisabled}>
                    {#if !editingDisabled}
                        <Button on:click={negateRule} variant="outline" size="icon" >
                            <Slash class="w-4"/>
                        </Button>
                        <Button on:click={deleteRule} variant="outline" size="icon">
                            <Trash class="w-4"/>
                        </Button>
                    {/if}
                </ComparisonRule>
            {:else}
                {#if internal.type === 'and' || internal.type === 'or'}
                    <div class="flex">
                        <Button on:click={toggleExpand} size="sm">{internal.type.toUpperCase()}</Button>
                        {#if !editingDisabled}
                            <Button
                                on:click={() => {
                                    if (!isRule(internal)) return
                                    internal.type = internal.type === 'and' ? 'or' : 'and';
                                }}
                                size="sm"
                            >
                                <Symbol />
                            </Button>
                        {/if}
                    </div>
                {:else}
                    <Button on:click={toggleExpand} size="sm">{'NOT'}</Button>
                {/if}

                {#if !editingDisabled}
                    {#if internal.type != 'not' && (parent && !isRootRule(parent) && parent.type !=  'not')}
                        <Button on:click={negateRule} size="sm">
                            <Slash class="w-4"/>
                        </Button>
                    {/if}

                    <Button on:click={deleteRule} size="sm">
                        <Trash class="w-4"/>
                    </Button>
                {/if}
            {/if}
        </div>

        <!-- Expanded block -->
        {#if expanded || isFirstRule || internal.type === 'not'}
            {#if internal.type === 'and' || internal.type === 'or'}
                <ul>
                    {#each internal.children as child}
                        <li>
                            <RuleInternalBuilder
                                bind:internal={child}
                                bind:parent={internal}
                                bind:secondParent={parent}
                                {sensors}
                                bind:editingDisabled={editingDisabled}
                            />
                        </li>
                    {/each}

                    {#if adding}
                        <ConditionBuilder
                            bind:open={adding}
                            {sensors}
                            bind:parent={internal}
                        />
                    {:else}
                        {#if !editingDisabled}
                            <li>
                                <Button on:click={addRule} variant="outline" size="sm">
                                    <Plus class="w-4"/>
                                </Button>
                            </li>
                        {/if}
                    {/if}
                </ul>
            {:else if internal.type === 'not'}
                <ul>
                    <li>
                        <RuleInternalBuilder
                            bind:internal={internal.wrapped}
                            bind:parent={internal}
                            bind:secondParent={parent}
                            {sensors}
                            bind:editingDisabled={editingDisabled}
                        />
                    </li>
                </ul>
            {/if}
        {/if}


    <!-- The internal is empty (first rule) -->
    {:else}
        {#if !editingDisabled}
            {#if adding}
                <ConditionBuilder
                    bind:open={adding}
                    {sensors}
                    bind:parent={parent}
                />
            {:else}
                <Button on:click={addRule} variant="outline" size="sm">
                    <Plus class="w-4"/>
                </Button>
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
