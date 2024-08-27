import { Duration } from '@/utils/Duration';

// Server stuff
export const SERVER_URL = 'http://localhost:8080';
export const getWSUrl = () => {
    return (SERVER_URL.length > 0 ? SERVER_URL : location.origin).replace(/^http/, 'ws');
};

// Sync configuration
export const SENSOR_VALUE_INTERVAL = Duration.ofSeconds(5);
