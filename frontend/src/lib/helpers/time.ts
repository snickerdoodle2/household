export function convertMsToTime(ms: number): {
    hours: number;
    minutes: number;
    seconds: number;
} {
    const seconds = Math.floor(Math.max(0, ms / 1000) % 60);
    const minutes = Math.floor(Math.max(0, ms / (1000 * 60)) % 60);
    const hours = Math.floor(Math.max(0, ms / (1000 * 60 * 60)));

    return {
        hours,
        minutes,
        seconds,
    };
}

export function convertTimeToMs(time: {
    hours: number;
    minutes: number;
    seconds: number;
}): number {
    const { hours, minutes, seconds } = time;

    const ms = hours * 60 * 60 * 1000 + minutes * 60 * 1000 + seconds * 1000;

    return ms;
}

type DurationObject = {
    hours: number;
    minutes: number;
    seconds: number;
    milliseconds: number;
    microseconds: number;
    nanoseconds: number;
};

export function parseDuration(durationStr: string): DurationObject {
    const regex = /(-?\d+(?:\.\d+)?)(h|m|s|ms|µs|us|ns)/g;

    const result: DurationObject = {
        hours: 0,
        minutes: 0,
        seconds: 0,
        milliseconds: 0,
        microseconds: 0,
        nanoseconds: 0,
    };

    let match;
    while ((match = regex.exec(durationStr)) !== null) {
        const value = parseFloat(match[1]);
        const unit = match[2];

        switch (unit) {
            case 'h':
                result.hours += value;
                break;
            case 'm':
                result.minutes += value;
                break;
            case 's':
                result.seconds += value;
                break;
            case 'ms':
                result.milliseconds += value;
                break;
            case 'µs':
            case 'us':
                result.microseconds += value;
                break;
            case 'ns':
                result.nanoseconds += value;
                break;
        }
    }

    return result;
}
