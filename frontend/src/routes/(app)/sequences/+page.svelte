<script lang="ts">
    import type { PageData } from './$types';
    import * as Card from '$lib/components/ui/card/index.js';
    import Button from '@/components/ui/button/button.svelte';
    import { Plus } from 'radix-icons-svelte';
    import { goto } from '$app/navigation';

    export let data: PageData;
</script>

<div class="flex h-full items-start">
    {#await data.sequences then sequences}
        <div
            class="grid flex-1 grid-cols-1 gap-8 py-20 sm:grid-cols-1 lg:grid-cols-2"
        >
            {#each sequences as sequence}
                <div class="lg:min-w-[32rem]">
                    <a href={`/sequences/${sequence.id}`}>
                        <Card.Root class="hover:bg-accent">
                            <Card.Header>
                                <Card.Title>
                                    {sequence.name}
                                </Card.Title>
                                <Card.Description>
                                    {sequence.description}
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
                    on:click={() => goto(`/sequences/create`)}
                >
                    <Plus />
                </Button>
            </div>
        </div>
    {/await}
</div>
