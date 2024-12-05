<script lang="ts">
    import * as Alert from '$lib/components/ui/alert';
    import * as Tooltip from '$lib/components/ui/tooltip';
    import type {
        NotificationLevel,
        Notification,
    } from '@/helpers/socket.svelte';
    import { cn } from '@/utils';
    import dayjs from 'dayjs';
    import relativeTime from 'dayjs/plugin/relativeTime';
    import { BadgeAlert, BadgeCheck, BadgeInfo, BadgeX } from 'lucide-svelte';

    type Props = {
        notification: Notification;
    };

    const { notification }: Props = $props();
    $effect(() => {
        dayjs.extend(relativeTime);
    });
</script>

{#snippet icon(level: NotificationLevel)}
    {@const common = 'w-7 h-7'}
    <Tooltip.Root>
        <Tooltip.Trigger>
            {#if level == 'info'}
                <BadgeInfo class={cn(common, 'stroke-sky-500')} />
            {:else if level == 'success'}
                <BadgeCheck class={cn(common, 'stroke-emerald-500')} />
            {:else if level == 'error'}
                <BadgeX class={cn(common, 'stroke-red-600')} />
            {:else if level == 'warning'}
                <BadgeAlert class={cn(common, 'stroke-amber-300')} />
            {/if}
        </Tooltip.Trigger>
        <Tooltip.Content>{level}</Tooltip.Content>
    </Tooltip.Root>
{/snippet}

<Alert.Root class="flex gap-4 items-center">
    {@render icon(notification.level)}
    <div class="flex-1">
        <Alert.Title class="text-base">{notification.title}</Alert.Title>
        <Alert.Description class="flex justify-between">
            <p>
                {notification.description}
            </p>
            <span>{dayjs(notification.created_at).fromNow()}</span>
        </Alert.Description>
    </div>
</Alert.Root>
