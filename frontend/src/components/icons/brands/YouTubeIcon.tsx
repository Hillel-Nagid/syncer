import { JSX } from "solid-js";

interface IconProps {
    class?: string;
}

export default function YouTubeIcon(props: IconProps): JSX.Element {
    return (
        <svg x="0px" y="0px" viewBox="0 0 176 176" enable-background="new 0 0 176 176" class={props.class}>
            <g id="XMLID_167_">
                <circle id="XMLID_791_" fill="#FF0000" cx="88" cy="88" r="88" />
                <path id="XMLID_42_" fill="#FFFFFF" d="M88,46c23.1,0,42,18.8,42,42s-18.8,42-42,42s-42-18.8-42-42S64.9,46,88,46 M88,42   c-25.4,0-46,20.6-46,46s20.6,46,46,46s46-20.6,46-46S113.4,42,88,42L88,42z" />
                <polygon id="XMLID_274_" fill="#FFFFFF" points="72,111 111,87 72,65  " />
            </g>
        </svg>
    );
} 