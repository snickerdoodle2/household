<script lang="ts">
import LightSwitch from "./LightSwitch.svelte";

const LINKS = [
	{
		name: "Home",
		url: "/",
	},
	{
		name: "Debug",
		url: "/debug",
	},
];
</script>

<nav
    class="h-svh fixed top-0 left-0 rounded-r-xl shadow-xl border-2
    bg-background py-6 px-4 flex flex-col gap-4"
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
