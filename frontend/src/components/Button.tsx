import { JSX } from "solid-js";
import type { IconName } from "~/types";
import Icon from "./Icon";

interface ButtonProps {
    children: JSX.Element;
    variant?: "primary" | "secondary";
    size?: "sm" | "md" | "lg" | "xl";
    onClick?: () => void;
    class?: string;
    type?: "button" | "submit" | "reset";
    icon?: IconName;
    iconPosition?: "left" | "right";
}

export default function Button(props: ButtonProps) {
    const variant = () => props.variant || "primary";
    const size = () => props.size || "md";
    const iconPosition = () => props.iconPosition || "left";

    const baseClasses = "font-semibold transition-colors rounded-lg cursor-pointer flex items-center";

    const variantClasses = () => {
        switch (variant()) {
            case "primary":
                return "bg-emerald-600 hover:bg-emerald-700 text-white";
            case "secondary":
                return "border border-emerald-600 text-emerald-600 hover:bg-emerald-600 hover:text-white dark:hover:bg-emerald-900/50 dark:hover:text-white";
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
            case "xl":
                return "px-12 py-4 text-xl";
            default:
                return "px-8 py-3";
        }
    };

    const iconSize = () => {
        switch (size()) {
            case "sm":
                return "w-4 h-4";
            case "md":
                return "w-5 h-5";
            case "lg":
                return "w-5 h-5";
            case "xl":
                return "w-6 h-6";
            default:
                return "w-5 h-5";
        }
    };

    const iconSpacing = () => {
        const spacing = iconPosition() === "left" ? "mr-2" : "ml-2";
        return props.icon ? spacing : "";
    };

    const classes = () =>
        `${baseClasses} ${variantClasses()} ${sizeClasses()} ${props.class || ""}`;

    return (
        <button
            type={props.type || "button"}
            class={classes()}
            onClick={props.onClick}
        >
            {props.icon && iconPosition() === "left" && (
                <Icon name={props.icon} class={`${iconSize()} ${iconSpacing()}`} />
            )}
            {props.children}
            {props.icon && iconPosition() === "right" && (
                <Icon name={props.icon} class={`${iconSize()} ${iconSpacing()}`} />
            )}
        </button>
    );
} 