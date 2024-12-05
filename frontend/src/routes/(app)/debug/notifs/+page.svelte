<script lang="ts">
    import { Button } from '@/components/ui/button';
    import { authFetch } from '@/helpers/fetch';
    import { AppWebsocket } from '@/helpers/socket.svelte';

    const requestAll = () => {
        authFetch('/api/v1/notification/debug', {
            method: 'POST',
        });
    };

    const ws = new AppWebsocket();
</script>

<div class="flex flex-col gap-8 items-center">
    <ol>
        {#each ws.notifications as notification}
            <li>
                ({notification.level}) {notification.title} - {notification.description}
            </li>
        {/each}
    </ol>

    <Button on:click={requestAll}>Request all</Button>
</div>
