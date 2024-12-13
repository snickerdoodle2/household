<script lang="ts">
    import * as Card from '$lib/components/ui/card';
    import { AppWebsocket } from '$lib/helpers/socket.svelte';
    import { Button } from '$lib/components/ui/button';
    import Notification from './Notification.svelte';

    const ws = new AppWebsocket();
</script>

<Card.Root class="w-[600px] border-none shadow-none">
    <Card.Header class="text-3xl flex flex-row justify-between items-center">
        <Card.Title>Notifications</Card.Title>
        <Button
            on:click={() => {
                ws.markAllAsRead();
            }}>Mark all as read</Button
        >
    </Card.Header>
    <Card.Content class="flex max-h-[600px] flex-col gap-3 overflow-y-auto">
        {#each ws.notifications as notification (notification.id)}
            <Notification
                {notification}
                onclick={() => ws.markNotificationAsRead(notification.id)}
            />
        {/each}
    </Card.Content>
</Card.Root>
