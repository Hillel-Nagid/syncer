import { JSX } from "solid-js";

interface IconProps {
    class?: string;
}

export default function LightningIcon(props: IconProps): JSX.Element {
    return (
        <svg class={props.class} fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
        </svg>
    );
} 