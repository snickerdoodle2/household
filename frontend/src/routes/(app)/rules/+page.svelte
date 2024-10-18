<script lang="ts">
    import type { PageData } from './$types';
    import * as Card from '$lib/components/ui/card/index.js';
    import Button from '@/components/ui/button/button.svelte';
    import { Plus } from 'radix-icons-svelte';
    import { goto } from '$app/navigation';

    export let data: PageData;
    import ruleData from '@/components/rule/ExampleRule.json';
    import type { RuleDetails } from '@/types/rule';

    const rule = ruleData as unknown as RuleDetails;
</script>

<div class="flex h-full items-start">
    {#await data.rules then rules}
        <div
            class="grid flex-1 grid-cols-1 gap-8 py-20 sm:grid-cols-1 lg:grid-cols-2"
        >
            {#each rules as rule}
                <div class="lg:min-w-[32rem]">
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
                </div>
            {/each}

            <div class="flex items-center justify-center">
                <Button
                    variant="outline"
                    size="icon"
                    on:click={() => goto(`/rules/create`)}
                >
                    <Plus />
                </Button>
            </div>
        </div>
    {/await}
</div>
