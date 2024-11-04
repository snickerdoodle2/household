<script lang="ts">
    import { Home, MagicWand, Gear, Exit, Enter, Checkbox } from 'svelte-radix';
    import LightSwitch from './LightSwitch.svelte';
    import { Button } from './ui/button';
    import { authToken } from '@/auth/token';

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
        <Button variant="outline" size="icon" class="h-11 w-11">
            <a href="/settings">
                <Gear class="scale-90" />
            </a>
        </Button>
        <LightSwitch />
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
