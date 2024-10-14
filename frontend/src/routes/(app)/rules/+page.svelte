<script lang="ts">
    import type { PageData } from './$types';
    import * as Card from '$lib/components/ui/card/index.js';
    export let data: PageData;
    import ruleData from '@/components/rule/ExampleRule.json';
    import RuleInternalBuilder from '@/components/rule/RuleInternalBuilder.svelte';
    import type { RuleDetails } from '@/types/rule';

    const rule = ruleData as unknown as RuleDetails
</script>

{#await data.rules then rules}
    <div class="flex flex-col gap-2 lg:min-w-[32rem]">
        {#each rules as rule}
            <a href={`/rules/${rule.id}`}>
                <Card.Root class="hover:bg-accent">
                    <Card.Header>
                        <Card.Title>
                            {rule.name}
                        </Card.Title>
                        <Card.Description>
                            {rule.description}
                        </Card.Description>
                    </Card.Header>
                </Card.Root>
            </a>
        {/each}
    </div>
{/await}

{#await data.sensors then sensors}
    <div>
        <RuleInternalBuilder internal={rule.internal} {sensors} parent={rule} secondParent={undefined}/>
    </div>
{/await}
