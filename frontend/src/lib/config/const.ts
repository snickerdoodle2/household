export const SERVER_URL = 'http://localhost:8080';

export const getWSUrl = () => {
    return (SERVER_URL.length > 0 ? SERVER_URL : location.origin).replace(
        /^http/,
        'ws'
    );
};
