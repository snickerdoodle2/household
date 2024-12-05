<script lang="ts">
    import * as Card from '$lib/components/ui/card';
    import { AppWebsocket } from '$lib/helpers/socket.svelte';
    import Notification from './Notification.svelte';

    const ws = new AppWebsocket();
</script>

<Card.Root class="w-[600px] border-none shadow-none">
    <Card.Header class="text-3xl">
        <Card.Title>Notifications</Card.Title>
    </Card.Header>
    <Card.Content class="flex max-h-[600px] flex-col gap-3 overflow-y-auto">
        {#each ws.notifications as notification}
            <Notification {notification} />
        {/each}
    </Card.Content>
</Card.Root>
