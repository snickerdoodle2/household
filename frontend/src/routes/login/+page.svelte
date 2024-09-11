<script lang="ts">
    import { Input } from '$lib/components/ui/input';
    import { Button } from '$lib/components/ui/button';
    import { loginSchema } from '@/types/login';
    import { authToken } from '@/auth/token';

    const debounce = (callback: Function) => {
        clearTimeout(timeout);
        timeout = setTimeout(callback, 300);
    };

    const validate = () => {
        const { error, success } = loginSchema.safeParse({
            username: +username,
            password,
        });
        if (success) return;
        errors = Object.fromEntries(
            error.issues.map((e) => [e.path[0], e.message])
        );
    };

    let timeout: ReturnType<typeof setTimeout>;
    let username = '';
    let password = '';
    let errors: { [key: string]: string } = {};

    $: {
        username;
        password;
        debounce(validate);
    }

    const handleLogin = async () => {
        const { data, success } = loginSchema.safeParse({ username, password });

        username = '';
        password = '';
        if (!success) return;
        const res = await fetch('/api/v1/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });

        if (!res.ok) {
            console.log(await res.json());
        }

        const token = (await res.json())['auth_token'];

        const err = authToken.set(token);
        if (err) {
            console.error(err);
        }
    };
</script>

<form class="flex flex-col gap-3" on:submit|preventDefault={handleLogin}>
    <Input placeholder="Username" name="username" bind:value={username} />
    <Input
        placeholder="Password"
        name="password"
        type="password"
        bind:value={password}
    />
    <Button type="submit">Login</Button>
</form>
