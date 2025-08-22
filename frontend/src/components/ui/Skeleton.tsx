interface SkeletonProps {
    variant?: 'text' | 'avatar' | 'card' | 'button' | 'custom';
    size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl';
    width?: string;
    height?: string;
    class?: string;
    rounded?: 'none' | 'sm' | 'md' | 'lg' | 'xl' | '2xl' | 'full';
    animate?: boolean;
}

export default function Skeleton(props: SkeletonProps) {
    const variant = () => props.variant || 'text';
    const size = () => props.size || 'md';
    const rounded = () => props.rounded || 'md';
    const animate = () => props.animate !== false; // Enable animation by default

    const baseClasses = 'bg-slate-200 dark:bg-slate-700';

    const animationClasses = () =>
        animate()
            ? 'animate-pulse bg-gradient-to-r from-slate-200 via-slate-300 to-slate-200 dark:from-slate-700 dark:via-slate-600 dark:to-slate-700 bg-[length:200%_100%] animate-shimmer'
            : '';

    const roundedClasses = () => {
        switch (rounded()) {
            case 'none':
                return 'rounded-none';
            case 'sm':
                return 'rounded-sm';
            case 'md':
                return 'rounded-md';
            case 'lg':
                return 'rounded-lg';
            case 'xl':
                return 'rounded-xl';
            case '2xl':
                return 'rounded-2xl';
            case 'full':
                return 'rounded-full';
            default:
                return 'rounded-md';
        }
    };

    const variantClasses = () => {
        switch (variant()) {
            case 'text':
                return getTextClasses();
            case 'avatar':
                return getAvatarClasses();
            case 'card':
                return getCardClasses();
            case 'button':
                return getButtonClasses();
            case 'custom':
                return getCustomClasses();
            default:
                return getTextClasses();
        }
    };

    const getTextClasses = () => {
        const heights = {
            xs: 'h-3',
            sm: 'h-4',
            md: 'h-5',
            lg: 'h-6',
            xl: 'h-7',
            '2xl': 'h-8',
        };

        const widths = {
            xs: 'w-16',
            sm: 'w-24',
            md: 'w-32',
            lg: 'w-48',
            xl: 'w-64',
            '2xl': 'w-80',
        };

        return `${heights[size()]} ${widths[size()]}`;
    };

    const getAvatarClasses = () => {
        const sizes = {
            xs: 'w-6 h-6',
            sm: 'w-8 h-8',
            md: 'w-10 h-10',
            lg: 'w-12 h-12',
            xl: 'w-16 h-16',
            '2xl': 'w-20 h-20',
        };

        return `${sizes[size()]} rounded-full`;
    };

    const getCardClasses = () => {
        const sizes = {
            xs: 'w-32 h-24',
            sm: 'w-48 h-32',
            md: 'w-64 h-40',
            lg: 'w-80 h-48',
            xl: 'w-96 h-56',
            '2xl': 'w-full h-64',
        };

        return `${sizes[size()]} rounded-2xl`;
    };

    const getButtonClasses = () => {
        const sizes = {
            xs: 'w-16 h-6',
            sm: 'w-20 h-8',
            md: 'w-24 h-10',
            lg: 'w-32 h-12',
            xl: 'w-40 h-14',
            '2xl': 'w-48 h-16',
        };

        return `${sizes[size()]} rounded-lg`;
    };

    const getCustomClasses = () => {
        const width = props.width || 'w-full';
        const height = props.height || 'h-4';
        return `${width} ${height}`;
    };

    const classes = () => {
        const customRounded = variant() === 'avatar' ? '' : roundedClasses();
        return `${baseClasses} ${animationClasses()} ${variantClasses()} ${customRounded} ${
            props.class || ''
        }`;
    };

    return <div class={classes()} />;
}

// Skeleton composition components for common use cases
export function SkeletonText(
    props: Omit<SkeletonProps, 'variant'> & { lines?: number }
) {
    const lines = () => props.lines || 1;

    if (lines() === 1) {
        return <Skeleton {...props} variant='text' />;
    }

    return (
        <div class='space-y-2'>
            {Array.from({ length: lines() }, (_, i) => (
                <Skeleton
                    {...props}
                    variant='text'
                    width={i === lines() - 1 ? 'w-3/4' : 'w-full'}
                />
            ))}
        </div>
    );
}

export function SkeletonCard(props: Omit<SkeletonProps, 'variant'>) {
    return (
        <div class='space-y-4'>
            <Skeleton {...props} variant='card' />
            <div class='space-y-2'>
                <SkeletonText size='lg' />
                <SkeletonText lines={2} size='sm' />
            </div>
        </div>
    );
}

export function SkeletonProfile(props: Omit<SkeletonProps, 'variant'>) {
    return (
        <div class='flex items-center space-x-2 p-2 rounded-lg'>
            <Skeleton {...props} variant='avatar' size='sm' />
            <SkeletonText size='sm' width='w-20' class='hidden sm:block' />
        </div>
    );
}
