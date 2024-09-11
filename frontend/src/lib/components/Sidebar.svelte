<script lang="ts">
    import LightSwitch from './LightSwitch.svelte';

    const LINKS = [
        {
            name: 'Home',
            url: '/',
        },
        {
            name: 'Debug',
            url: '/debug',
        },
        {
            name: 'Login',
            url: '/login',
        },
    ];
</script>

<nav
    class="fixed left-0 top-0 flex h-svh flex-col gap-4
    rounded-r-xl border-2 bg-background px-4 py-6 shadow-xl"
>
    <span>LOGO</span>
    <ul>
        {#each LINKS as link}
            <li>
                <a href={link.url} class="font-semibold">{link.name}</a>
            </li>
        {/each}
    </ul>
    <LightSwitch />
</nav>
