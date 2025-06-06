interface SyncIconProps {
    size?: "sm" | "md" | "lg" | "xl";
    class?: string;
}

export default function SyncIcon(props: SyncIconProps) {
    const sizeClasses = () => {
        switch (props.size || "md") {
            case "sm":
                return "w-4 h-4";
            case "md":
                return "w-8 h-8";
            case "lg":
                return "w-12 h-12";
            case "xl":
                return "w-16 h-16";
            default:
                return "w-8 h-8";
        }
    };

    return (
        <svg
            class={`${sizeClasses()} ${props.class || ""}`}
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
        >
            <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
            />
        </svg>
    );
} 