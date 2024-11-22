<script lang="ts">
    import { Home, MagicWand, Exit, Enter, Checkbox } from 'svelte-radix';
    import LightSwitch from './LightSwitch.svelte';
    import { Button } from './ui/button';
    import { authToken } from '@/auth/token';
    import * as Dialog from '$lib/components/ui/dialog';
    import { Bell } from 'radix-icons-svelte';
    import Notifications from './Notifications.svelte';
    import { Avatar } from 'radix-icons-svelte';

    const LINKS = [
        {
            icon: Home,
            name: 'Home',
            url: '/',
        },
        {
            icon: Checkbox,
            name: 'Rules',
            url: '/rules',
        },
        {
            icon: MagicWand,
            name: 'Debug',
            url: '/debug',
        },
    ];

    // Notification stuff
    let notificationsOpen = false;
    // let notificationSocket = initializeNotificationSocket()

    // function initializeNotificationSocket() {
    //     const token = get(authToken);
    //     if (!token) throw new Error('auth token is required');

    //     const { set, subscribe } = writable<NotificationSocketMessage | null>();

    //     const url = new URL(`/api/v1/notifications`, window.location.href);
    //     url.protocol = url.protocol.replace('http', 'ws');
    //     url.searchParams.set('token', token.token);

    //     const socket = new WebSocket(url.toString());

    //     socket.addEventListener('message', (message) => {
    //         const { data, success } = notificationSocketMessageSchema.safeParse(
    //             JSON.parse(message.data)
    //         );
    //         if (!success) return;

    //         set(data.sort((n1, n2) => n2.date.getTime() - n1.date.getTime()););
    //     });

    //     return {
    //         subscribe,
    //         close: () => {
    //             socket.close();
    //         },
    //     };
    // };

</script>

<nav
    class="shadow-32 group fixed left-0 top-0 flex
    h-svh w-20 flex-col gap-4 overflow-x-hidden rounded-r-xl
    border-2 bg-background px-4 py-6 transition-all hover:w-64"
>
    <span>LOGO</span>
    <hr />
    <ul class="flex flex-1 flex-col gap-6 group-hover:w-[14rem]">
        {#each LINKS as link}
            <li>
                <a
                    href={link.url}
                    class="flex h-12 items-center gap-8 rounded-md px-2 hover:bg-secondary"
                >
                    <link.icon class="h-6 w-6" />
                    <span class="hidden pb-1 text-base group-hover:inline"
                        >{link.name}</span
                    >
                </a>
            </li>
        {/each}
    </ul>
    <hr />
    <div class="flex justify-between group-hover:w-[14rem]">
        <Button
            variant="outline"
            size="icon"
            class="h-11 w-11"
            on:click={() => (notificationsOpen = true)}
        >
            <Bell class="scale-150"/>
        </Button>

        <Button
            variant="outline"
            size="icon"
            class="hidden h-11 w-11 group-hover:inline-flex"
        >
            <a href="/users">
                <Avatar class="scale-150" />
            </a>
        </Button>

        <LightSwitch />
        <!-- TODO: move this to the settings -->

        {#if $authToken != undefined}
            <Button
                variant="outline"
                size="icon"
                class="hidden h-11 w-11 group-hover:inline-flex"
                on:click={() => {
                    authToken.logout();
                }}
            >
                <Exit />
            </Button>
        {:else}
            <Button
                variant="outline"
                size="icon"
                class="hidden h-11 w-11 group-hover:inline-flex"
            >
                <a href="/login">
                    <Enter />
                </a>
            </Button>
        {/if}
    </div>
</nav>

<Dialog.Root bind:open={notificationsOpen}>
    <Dialog.Portal>
        <Dialog.Overlay />
        <Dialog.Content
            class="flex max-w-none items-center justify-center px-8 py-4 md:w-fit"
        >
            <Notifications />
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>
