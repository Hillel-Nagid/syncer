import { JSX } from "solid-js";

interface CardProps {
    children: JSX.Element;
    variant?: "default" | "colored" | "numbered";
    class?: string;
    backgroundColor?: string;
    padding?: "sm" | "md" | "lg" | "xl" | "none";
    hover?: boolean;
}

export default function Card(props: CardProps) {
    const variant = () => props.variant || "default";
    const padding = () => props.padding || "lg";
    const enableHover = () => props.hover !== false; // Enable hover by default

    const baseClasses = "rounded-2xl transition-all duration-300 ease-in-out";

    const paddingClasses = () => {
        switch (padding()) {
            case "none":
                return "";
            case "sm":
                return "p-4";
            case "md":
                return "p-6";
            case "lg":
                return "p-8";
            case "xl":
                return "p-12";
            default:
                return "p-8";
        }
    };

    const variantClasses = () => {
        const hoverClasses = enableHover() ? "hover:shadow-2xl hover:-translate-y-1" : "";

        switch (variant()) {
            case "default":
                return `bg-white dark:bg-slate-800 shadow-xl border border-slate-200 dark:border-slate-700 ${hoverClasses}`;
            case "colored":
                if (props.backgroundColor) {
                    // When custom background color is provided, ensure it works with dark theme
                    return `${props.backgroundColor} text-white dark:text-white shadow-xl ${hoverClasses}`;
                }
                return `bg-emerald-600 dark:bg-emerald-700 text-white shadow-xl ${hoverClasses}`;
            case "numbered":
                return `bg-white dark:bg-slate-800 shadow-xl border border-slate-200 dark:border-slate-700 ${hoverClasses}`;
            default:
                return `bg-white dark:bg-slate-800 shadow-xl border border-slate-200 dark:border-slate-700 ${hoverClasses}`;
        }
    };

    const classes = () => `${baseClasses} ${paddingClasses()} ${variantClasses()} ${props.class || ""}`;

    return (
        <div class={classes()}>
            {props.children}
        </div>
    );
} 