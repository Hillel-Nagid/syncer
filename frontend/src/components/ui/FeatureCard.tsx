import { JSX } from "solid-js";
import type { IconName } from "~/types";
import Icon from "./Icon";

interface FeatureCardProps {
    icon: IconName;
    title: string;
    description: string;
    class?: string;
    iconContainerClass?: string;
    variant?: "default" | "colored";
    hover?: boolean;
}

export default function FeatureCard(props: FeatureCardProps): JSX.Element {
    const variant = () => props.variant || "default";
    const enableHover = () => props.hover !== false; // Enable hover by default

    const baseClasses = "rounded-2xl p-8 transition-all duration-300 ease-in-out";

    const variantClasses = () => {
        const hoverClasses = enableHover() ? "hover:shadow-2xl hover:-translate-y-1" : "";

        switch (variant()) {
            case "default":
                return `bg-white dark:bg-slate-800 shadow-xl border border-slate-200 dark:border-slate-700 ${hoverClasses}`;
            case "colored":
                return `bg-emerald-600 dark:bg-emerald-700 text-white shadow-xl ${hoverClasses}`;
            default:
                return `bg-white dark:bg-slate-800 shadow-xl border border-slate-200 dark:border-slate-700 ${hoverClasses}`;
        }
    };

    const iconContainerClasses = () => {
        const base = "w-12 h-12 rounded-lg flex items-center justify-center mb-4 transition-all duration-300 ease-in-out";
        const hoverClasses = enableHover() ? "group-hover:scale-110" : "";

        if (props.iconContainerClass) {
            return `${base} ${hoverClasses} ${props.iconContainerClass}`;
        }

        switch (variant()) {
            case "colored":
                return `${base} ${hoverClasses} bg-emerald-400 dark:bg-emerald-500`;
            default:
                return `${base} ${hoverClasses} bg-emerald-100 dark:bg-emerald-900/30`;
        }
    };

    const iconClasses = () => {
        const baseIconClasses = "transition-all duration-300 ease-in-out";
        const hoverClasses = enableHover() ? "group-hover:scale-110" : "";

        switch (variant()) {
            case "colored":
                return `w-6 h-6 text-white ${baseIconClasses} ${hoverClasses}`;
            default:
                return `w-6 h-6 text-emerald-600 dark:text-emerald-400 ${baseIconClasses} ${hoverClasses}`;
        }
    };

    const titleClasses = () => {
        const baseClasses = "transition-colors duration-300 ease-in-out";

        switch (variant()) {
            case "colored":
                return `text-xl font-semibold mb-3 text-white ${baseClasses}`;
            default:
                return `text-xl font-semibold mb-3 text-slate-900 dark:text-slate-100 ${baseClasses}`;
        }
    };

    const descriptionClasses = () => {
        const baseClasses = "transition-colors duration-300 ease-in-out";

        switch (variant()) {
            case "colored":
                return `text-emerald-100 dark:text-emerald-200 ${baseClasses}`;
            default:
                return `text-slate-600 dark:text-slate-400 ${baseClasses}`;
        }
    };

    const cardClasses = () => {
        const groupClass = enableHover() ? "group" : "";
        return `${baseClasses} ${variantClasses()} ${groupClass} ${props.class || ""}`;
    };

    return (
        <div class={cardClasses()}>
            <div class={iconContainerClasses()}>
                <Icon name={props.icon} class={iconClasses()} />
            </div>
            <h3 class={titleClasses()}>{props.title}</h3>
            <p class={descriptionClasses()}>
                {props.description}
            </p>
        </div>
    );
} 