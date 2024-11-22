import type { Result } from '@/types/result';
// import type { FetchFn } from '@/types/misc';
import { type User } from '@/types/user';

export const getAllUsers = async () // fetch: FetchFn
: Promise<Result<User[], string>> => {
    // const res = await authFetch(`/api/v1/user`, {}, fetch);
    // const data = await res.json();
    // if (!res.ok) {
    //     return {
    //         isError: true,
    //         error: data.error,
    //     };
    // }

    // const parsed = z.object({ data: userSchema.array() }).safeParse(data);
    // if (!parsed.success) {
    //     return {
    //         isError: true,
    //         error: 'Error while parsing the data! (getAllUsers)',
    //     };
    // }

    // return {
    //     isError: false,
    //     data: parsed.data.data,
    // };

    return {
        isError: false,
        data: [
            {
                id: '3e7cbb34-86a4-4f9e-8327-1e7a2f7a50d4',
                username: 'johndoe',
                name: 'John Doe',
                created_at: new Date('2024-01-15T08:32:22Z'),
            },
            {
                id: 'ff62b7bc-9d57-4a69-97be-cd087c68d8e3',
                username: 'janedoe',
                name: 'Jane Doe',
                created_at: new Date('2023-11-05T14:22:18Z'),
            },
            {
                id: '7fcb31b5-28d2-4c10-bf3f-b87354efcd91',
                username: 'sam_smith',
                name: 'Sam Smith',
                created_at: new Date('2022-08-27T19:45:33Z'),
            },
            {
                id: '6454fc76-2a6f-4ae5-8a5b-bf23981231cf',
                username: 'alexjones',
                name: 'Alex Jones',
                created_at: new Date('2023-03-09T11:08:47Z'),
            },
            {
                id: 'd9d77b26-ef1d-4fc1-8efb-56d5b73f0d8a',
                username: 'lindab',
                name: 'Linda Brown',
                created_at: new Date('2023-07-21T09:12:55Z'),
            },
        ],
    };
};
