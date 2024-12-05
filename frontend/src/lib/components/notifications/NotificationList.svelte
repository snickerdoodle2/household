<script lang="ts">
    import * as Card from '$lib/components/ui/card';
    import { AppWebsocket } from '$lib/helpers/socket.svelte';
    import Notification from './Notification.svelte';

    const ws = new AppWebsocket();

    const onReadClick = (id: string) => {
        console.log(id);
    };
</script>

<Card.Root class="w-[600px] border-none shadow-none">
    <Card.Header class="text-3xl">
        <Card.Title>Notifications</Card.Title>
    </Card.Header>
    <Card.Content class="flex max-h-[600px] flex-col gap-3 overflow-y-auto">
        {#each ws.notifications as notification (notification.id)}
            <Notification
                {notification}
                onclick={() => onReadClick(notification.id)}
            />
        {/each}
    </Card.Content>
</Card.Root>
