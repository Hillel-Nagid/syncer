import { JSX } from "solid-js";

interface IconProps {
    class?: string;
}

export default function MenuIcon(props: IconProps): JSX.Element {
    return (
        <svg class={props.class} fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
        </svg>
    );
} 