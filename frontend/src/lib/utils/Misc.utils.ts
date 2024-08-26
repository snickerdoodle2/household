import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';
import { cubicOut } from 'svelte/easing';
import type { TransitionConfig } from 'svelte/transition';
import type { Result } from '@/types/Result.types';

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs));
}

type FlyAndScaleParams = {
    y?: number;
    x?: number;
    start?: number;
    duration?: number;
};

export const flyAndScale = (
    node: Element,
    params: FlyAndScaleParams = { y: -8, x: 0, start: 0.95, duration: 150 }
): TransitionConfig => {
    const style = getComputedStyle(node);
    const transform = style.transform === 'none' ? '' : style.transform;

    const scaleConversion = (
        valueA: number,
        scaleA: [number, number],
        scaleB: [number, number]
    ) => {
        const [minA, maxA] = scaleA;
        const [minB, maxB] = scaleB;

        const percentage = (valueA - minA) / (maxA - minA);
        const valueB = percentage * (maxB - minB) + minB;

        return valueB;
    };

    const styleToString = (style: Record<string, number | string | undefined>): string => {
        return Object.keys(style).reduce((str, key) => {
            if (style[key] === undefined) return str;
            // biome-ignore lint: code generated by shadcn-ui
            return str + `${key}:${style[key]};`;
        }, '');
    };

    return {
        duration: params.duration ?? 200,
        delay: 0,
        css: (t) => {
            const y = scaleConversion(t, [0, 1], [params.y ?? 5, 0]);
            const x = scaleConversion(t, [0, 1], [params.x ?? 0, 0]);
            const scale = scaleConversion(t, [0, 1], [params.start ?? 0.95, 1]);

            return styleToString({
                transform: `${transform} translate3d(${x}px, ${y}px, 0) scale(${scale})`,
                opacity: t,
            });
        },
        easing: cubicOut,
    };
};

export function validateName(name: string): Result<string, string> {
    // Validate the new category
    const isValidCategory = /^[A-Za-z0-9_-]+$/.test(name);
    const isValidLength = name.length >= 3 && name.length <= 15;
    const hasNoWhitespaces = !/\s/.test(name);

    if (!hasNoWhitespaces) {
        return {
            isError: true,
            error: 'Category name must not contain any whitespace characters.',
        };
    }

    if (!isValidCategory) {
        return {
            isError: true,
            error: 'Category name must be a single word containing only letters, numbers, underscores, and hyphens.',
        };
    }

    if (!isValidLength) {
        return {
            isError: true,
            error: 'Category name must be between 3 and 15 characters long.',
        };
    }

    return {
        isError: false,
        data: name,
    };
}
