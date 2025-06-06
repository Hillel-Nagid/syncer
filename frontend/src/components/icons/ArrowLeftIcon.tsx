import { JSX } from "solid-js";

interface ArrowLeftIconProps {
    class?: string;
}

export default function ArrowLeftIcon(props: ArrowLeftIconProps): JSX.Element {
    return (
        <svg
            class={props.class}
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
        >
            <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M10 19l-7-7m0 0l7-7m-7 7h18"
            />
        </svg>
    );
} 