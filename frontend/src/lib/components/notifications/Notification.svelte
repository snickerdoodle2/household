<script lang="ts">
    import * as Alert from '$lib/components/ui/alert';
    import type {
        NotificationLevel,
        Notification,
    } from '@/helpers/socket.svelte';
    import { cn } from '@/utils';
    import { BadgeAlert, BadgeCheck, BadgeInfo, BadgeX } from 'lucide-svelte';

    type Props = {
        notification: Notification;
    };

    const { notification }: Props = $props();
</script>

{#snippet icon(level: NotificationLevel)}
    {@const common = 'w-7 h-7'}
    {#if level == 'info'}
        <BadgeInfo class={cn(common, 'stroke-sky-500')} />
    {:else if level == 'success'}
        <BadgeCheck class={cn(common, 'stroke-emerald-500')} />
    {:else if level == 'error'}
        <BadgeX class={cn(common, 'stroke-red-600')} />
    {:else if level == 'warning'}
        <BadgeAlert class={cn(common, 'stroke-amber-300')} />
    {/if}
{/snippet}

<Alert.Root class="flex gap-4 items-center">
    {@render icon(notification.level)}
    <div class="">
        <Alert.Title>{notification.title}</Alert.Title>
        <Alert.Description>{notification.description}</Alert.Description>
    </div>
</Alert.Root>
