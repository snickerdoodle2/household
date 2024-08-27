/**
 * Represents a duration in various time units.
 */
export class Duration {
    /**
     * The duration in milliseconds.
     */
    private milliseconds: number;

    /**
     * Initializes a new Duration instance.
     *
     * @param {number} [milliseconds] - The duration in milliseconds. Defaults to 0.
     */
    constructor(milliseconds?: number) {
        this.milliseconds = milliseconds ?? 0;
    }

    /**
     * Adds another Duration to this Duration.
     *
     * @param {Duration} other - The other Duration.
     * @returns {Duration} - A new Duration representing the sum.
     */
    public add(other: Duration): Duration {
        return new Duration(this.milliseconds + other.toMilliseconds());
    }

    /**
     * Subtracts another Duration from this Duration.
     *
     * @param {Duration} other - The other Duration.
     * @returns {Duration} - A new Duration representing the difference.
     */
    public subtract(other: Duration): Duration {
        return new Duration(this.milliseconds - other.toMilliseconds());
    }

    /**
     * Sums up the given Durations.
     *
     * @param durations - The Durations to sum up.
     * @returns {Duration} - A new Duration representing the sum.
     */
    public static sum(...durations: Duration[]): Duration {
        const sum: number = durations.reduce(
            (acc: number, duration: Duration) => acc + duration.toMilliseconds(),
            0
        );
        return new Duration(sum);
    }

    /**
     * Creates a Duration from milliseconds.
     *
     * @param {number} milliseconds - The number of milliseconds.
     * @returns {Duration} - A new Duration representing the milliseconds.
     */
    public static ofMilliseconds(milliseconds: number): Duration {
        return new Duration(milliseconds);
    }

    /**
     * Creates a Duration from seconds.
     *
     * @param {number} seconds - The number of seconds.
     * @returns {Duration} - A new Duration representing the seconds.
     */
    public static ofSeconds(seconds: number): Duration {
        return new Duration(seconds * 1000);
    }

    /**
     * Creates a Duration from minutes.
     *
     * @param {number} minutes - The number of minutes.
     * @returns {Duration} - A new Duration representing the minutes.
     */
    public static ofMinutes(minutes: number): Duration {
        return new Duration(minutes * 60 * 1000);
    }

    /**
     * Creates a Duration with a value of 9,007,199,254,740,991
     *
     * @returns {Duration} - A new Duration representing the maximum value.
     */
    public static ofMaximum(): Duration {
        return new Duration(Number.MAX_SAFE_INTEGER);
    }

    /**
     * Converts the Duration to milliseconds.
     *
     * @returns {number} - The duration in milliseconds.
     */
    public toMilliseconds(): number {
        return this.milliseconds;
    }

    /**
     * Converts the Duration to seconds.
     *
     * @returns {number} - The duration in seconds.
     */
    public toSeconds(): number {
        return Math.ceil(this.milliseconds / 1000);
    }

    /**
     * Converts the Duration to seconds.
     *
     * @returns {number} - The duration in seconds.
     */
    public toExactSeconds(): number {
        return this.milliseconds / 1000;
    }

    /**
     * Converts the Duration to minutes.
     *
     * @returns {number} - The duration in minutes.
     */
    public toMinutes(): number {
        return Math.ceil(this.milliseconds / 60000);
    }
}
