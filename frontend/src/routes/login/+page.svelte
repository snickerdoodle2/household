<script lang="ts">
import { run, preventDefault } from 'svelte/legacy';

import { Input } from '$lib/components/ui/input';
import { Button } from '$lib/components/ui/button';
import { loginSchema } from '@/types/login';
import { authToken } from '@/auth/token';
import { goto } from '$app/navigation';

// eslint-disable-next-line @typescript-eslint/no-unsafe-function-type
const debounce = (callback: Function, ...args: unknown[]) => {
    clearTimeout(timeout);
    timeout = window.setTimeout(() => callback(args), 300);
};

const validate = () => {
    const { error, success } = loginSchema.safeParse({
        username,
        password,
    });
    if (success) return;
    errors = Object.fromEntries(
        error.issues.map((e) => [e.path[0], e.message])
    );
    console.log(errors);
};

let timeout: number;
let username = $state('');
let password = $state('');
let errors: Record<string, string> = {};

run(() => {
    debounce(validate, username, password);
});

const handleLogin = async () => {
    const { data, success } = loginSchema.safeParse({ username, password });
    if (!success) return;
    const res = await authToken.login(data);
    if (!res) {
        username = '';
        password = '';
        goto('/');
    }
};
</script>

<main class="flex h-svh items-center justify-center">
    <form class="flex flex-col gap-3" onsubmit={preventDefault(handleLogin)}>
        <Input placeholder="Username" name="username" bind:value={username} />
        <Input
            placeholder="Password"
            name="password"
            type="password"
            bind:value={password}
        />
        <Button type="submit">Login</Button>
    </form>
</main>
