<script>
    import { ModalType } from '@/types/Modal.types';
    import { openModal } from '@/utils/Modal.utils';
    import {
        Button,
        Dropdown,
        DropdownItem,
        Navbar,
        NavBrand,
    } from 'flowbite-svelte';
    import { ListOutline } from 'flowbite-svelte-icons';
    import { DotsHorizontal } from 'radix-icons-svelte';

    const categories = [
        'Kitchen',
        'Living Room',
        'Bedroom',
        'Bathroom',
        'Garage',
        'Garden',
    ];
</script>

<main class="bg-primary px-[2.5vw] py-[2.5vh]">
    <Navbar class="bg-card flex rounded-lg h-[10vh] md-[2.5vh]">
        <NavBrand href="/">
            <div class="flex items-center">
                <button
                    class="flex items-center bg-transparent border-none cursor-pointer p-0"
                >
                    <ListOutline
                        class="w-6 h-6 mr-2 text-white dark:text-white"
                    />
                    <span class="font-semibold text-xl"
                        >Nazwa Naszego Systemu</span
                    >
                </button>
                <Dropdown
                    class="bg-popover absolute mt-2 w-[5hv] right-0 top-1/2 rounded-lg"
                >
                    <DropdownItem>Dashboard</DropdownItem>
                    <DropdownItem>Settings</DropdownItem>
                    <DropdownItem>Earnings</DropdownItem>
                    <DropdownItem>Sign out</DropdownItem>
                </Dropdown>
            </div>
        </NavBrand>
    </Navbar>

    <div class="bg-card rounded-lg h-[80vh] my-[2.5vw]">
        <div
            class="flex items-center justify-between bg-black w-full h-0.3 rounded-lg p-2"
        >
            <div class="flex flex-nowrap flex-row space-x-4 overflow-x-auto">
                {#each categories as category}
                    <Button class="bg-secondary rounded-lg">
                        {category}
                    </Button>
                {/each}
            </div>

            <Button
                on:click={() => openModal(ModalType.ADD_SENSOR, undefined)}
                class="bg-secondary rounded-lg"
            >
                <DotsHorizontal
                    class="w-6 h-6 mr-2 text-white dark:text-white"
                />
            </Button>
        </div>

        <div class="flex flex-wrap overflow-y-auto p-2">

            <button class="flex px-3 py-4 bg-black rounded-full">
                <p class="text-4xl">🌐</p> 
                <div class="flex flex-col px-2">
                    <p class="text-md">Termometr przy oknie</p>
                    <p class="text-sm">25🌡️</p>
                </div>
            </button>
            <!-- <h1 class="text-4xl font-bold text-white">Sensor Page</h1> -->
        </div>
    </div>
</main>
