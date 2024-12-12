<script lang="ts">
    import type { PageData } from './$types';
    import * as Card from '$lib/components/ui/card/index.js';
    import Button from '@/components/ui/button/button.svelte';
    import { Plus } from 'radix-icons-svelte';

    type Props = {
        data: PageData;
    };

    let { data }: Props = $props();
</script>

<div class="flex flex-col h-full w-full px-4 py-4 gap-6">
    <div class="flex justify-end">
        {#if data.currentUser.role === 'admin'}
            <Button variant="outline" size="icon" href="/sequences/create">
                <Plus class="w-6 h-6" />
            </Button>
        {/if}
    </div>
    <div class="flex-1 flex flex-col gap-8 items-center">
        {#await data.sequences then sequences}
            {#each sequences as sequence}
                <div class="lg:min-w-[48rem]">
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
        {/await}
    </div>
</div>
