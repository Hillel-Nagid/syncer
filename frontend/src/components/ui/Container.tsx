import { JSX } from "solid-js";

interface ContainerProps {
    children: JSX.Element;
    class?: string;
    maxWidth?: "sm" | "md" | "lg" | "xl" | "2xl" | "4xl";
}

export default function Container(props: ContainerProps): JSX.Element {
    const maxWidth = () => props.maxWidth || "4xl";

    const maxWidthClasses = () => {
        switch (maxWidth()) {
            case "sm":
                return "max-w-sm";
            case "md":
                return "max-w-md";
            case "lg":
                return "max-w-lg";
            case "xl":
                return "max-w-xl";
            case "2xl":
                return "max-w-2xl";
            case "4xl":
                return "max-w-4xl";
            default:
                return "max-w-4xl";
        }
    };

    const classes = () => `${maxWidthClasses()} mx-auto ${props.class || ""}`;

    return (
        <div class={classes()}>
            {props.children}
        </div>
    );
} 