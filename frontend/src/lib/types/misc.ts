export type FetchFn = (
    input: RequestInfo | URL,
    init?: RequestInit | undefined
) => Promise<Response>;
