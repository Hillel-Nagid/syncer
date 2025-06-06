import { JSX } from "solid-js";

interface IconProps {
    class?: string;
}

export default function TidalIcon(props: IconProps): JSX.Element {
    return (
        <svg class={props.class} xmlns="http://www.w3.org/2000/svg" width="800px" height="800px" viewBox="0 0 48 48"><defs></defs><path class="a" d="M10.3333,12.3333l6.8334,6.8334L24,12.3333l6.8333,6.8334,6.8334-6.8334L44.5,19.1667,37.6667,26l-6.8334-6.8333L24,26l6.8333,6.8333L24,39.6667l-6.8333-6.8334L24,26l-6.8333-6.8333L10.3333,26,3.5,19.1667Z" /></svg>
    );
} 