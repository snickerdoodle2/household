<script lang="ts">
    import { page } from '$app/stores';
    import { Button } from '@/components/ui/button';
    import type { Snippet } from 'svelte';
    import type { LayoutData } from './$types';

    type Props = {
        children: Snippet;
        data: LayoutData;
    };

    const { children }: Props = $props();

    let links = [
        {
            label: 'Users',
            href: '/settings/users',
        },
    ];
</script>

{#snippet Link(label: string, href: string)}
    <Button
        variant="link"
        size={$page.url.pathname.startsWith(href) ? 'bold' : 'default'}
        class="text-base"
        disabled={$page.url.pathname === href}
        {href}
    >
        {label}
    </Button>
{/snippet}

<div class="w-full h-full flex flex-row pt-32 pl-40 gap-8">
    <nav>
        {#each links as link (link.href)}
            {@render Link(link.label, link.href)}
        {/each}
    </nav>
    <main class="flex-1">{@render children()}</main>
</div>
