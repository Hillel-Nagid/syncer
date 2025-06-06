import { JSX } from "solid-js";

interface ButtonProps {
    children: JSX.Element;
    variant?: "primary" | "secondary";
    size?: "sm" | "md" | "lg";
    onClick?: () => void;
    class?: string;
    type?: "button" | "submit" | "reset";
}

export default function Button(props: ButtonProps) {
    const variant = () => props.variant || "primary";
    const size = () => props.size || "md";

    const baseClasses = "font-semibold transition-colors rounded-lg cursor-pointer";

    const variantClasses = () => {
        switch (variant()) {
            case "primary":
                return "bg-emerald-600 hover:bg-emerald-700 text-white";
            case "secondary":
                return "border border-emerald-600 text-emerald-600 hover:bg-emerald-50 dark:hover:bg-emerald-900/20";
            default:
                return "bg-emerald-600 hover:bg-emerald-700 text-white";
        }
    };

    const sizeClasses = () => {
        switch (size()) {
            case "sm":
                return "px-4 py-2 text-sm";
            case "md":
                return "px-8 py-3";
            case "lg":
                return "px-8 py-3 text-lg";
            default:
                return "px-8 py-3";
        }
    };

    const classes = () =>
        `${baseClasses} ${variantClasses()} ${sizeClasses()} ${props.class || ""}`;

    return (
        <button
            type={props.type || "button"}
            class={classes()}
            onClick={props.onClick}
        >
            {props.children}
        </button>
    );
} 