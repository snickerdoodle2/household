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

export function generateDurationString(startDate: Date, endDate: Date): string {
    // Calculate the difference in milliseconds
    const diffMs = endDate.getTime() - startDate.getTime();

    if (diffMs < 0) {
        throw new Error('End date must be greater than start date');
    }

    // Duration in various units
    const units: Record<string, number> = {
        ns: 1e-6, // nanoseconds
        us: 1e-3, // microseconds
        ms: 1, // milliseconds
        s: 1000, // seconds
        m: 60000, // minutes
        h: 3600000, // hours
        d: 86400000, // days (24 hours)
    };

    // Iterate over the units and calculate the appropriate duration
    for (const unit in units) {
        const value = diffMs / units[unit];
        if (value >= 1) {
            // Format the duration to the nearest whole or decimal value
            const formattedValue =
                value % 1 === 0 ? Math.floor(value) : value.toFixed(1);
            return `${formattedValue}${unit}`;
        }
    }

    // If no suitable unit found (e.g., microsecond precision)
    return `${diffMs}ms`;
}
