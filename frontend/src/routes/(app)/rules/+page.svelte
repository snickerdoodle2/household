<script lang="ts">
    import type { PageData } from './$types';
    import * as Card from '$lib/components/ui/card/index.js';
    import Button from '@/components/ui/button/button.svelte';
    import { Plus } from 'lucide-svelte';

    type Props = {
        data: PageData;
    };

    let { data }: Props = $props();
</script>

<div class="flex flex-col h-full w-full px-4 py-4 gap-6">
    <div class="flex justify-end">
        {#if data.currentUser.role === 'admin'}
            <Button variant="outline" size="icon" href="/rules/create">
                <Plus class="w-6 h-6" />
            </Button>
        {/if}
    </div>
    <div class="flex-1 flex flex-col gap-8 items-center">
        {#await data.rules then rules}
            {#each rules as rule}
                <div class="lg:min-w-[48rem]">
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
        {/await}
    </div>
</div>
