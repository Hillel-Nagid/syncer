import { JSX } from "solid-js";

interface IconProps {
    class?: string;
}

export default function DeezerIcon(props: IconProps): JSX.Element {
    return (
        <svg class={props.class} xmlns="http://www.w3.org/2000/svg" viewBox="0 0 250 163.17">
            <defs>
                <linearGradient id="e" x1="90.71" x2="90.5" y1="62.45" y2="49.29" gradientTransform="matrix(2.4611 0 0 2.4611 0 -77.943)" gradientUnits="userSpaceOnUse">
                    <stop stop-color="#358c7b" offset="0" />
                    <stop stop-color="#33a65e" offset=".5256" />
                </linearGradient>
                <linearGradient id="a" x1="79.29" x2="101.91" y1="79.49" y2="67.97" gradientTransform="matrix(2.4611 0 0 2.4611 0 -77.943)" gradientUnits="userSpaceOnUse">
                    <stop stop-color="#222b90" offset="0" />
                    <stop stop-color="#367b99" offset="1" />
                </linearGradient>
                <linearGradient id="b" x2="21.95" y1="91.55" y2="91.55" gradientTransform="matrix(2.4611 0 0 2.4611 0 -77.943)" gradientUnits="userSpaceOnUse">
                    <stop stop-color="#f90" offset="0" />
                    <stop stop-color="#ff8000" offset="1" />
                </linearGradient>
                <linearGradient id="h" x1="26.55" x2="48.49" y1="91.55" y2="91.55" gradientTransform="matrix(2.4611 0 0 2.4611 0 -77.943)" gradientUnits="userSpaceOnUse">
                    <stop stop-color="#ff8000" offset="0" />
                    <stop stop-color="#cc1953" offset="1" />
                </linearGradient>
                <linearGradient id="g" x1="53.09" x2="75.03" y1="91.55" y2="91.55" gradientTransform="matrix(2.4611 0 0 2.4611 0 -77.943)" gradientUnits="userSpaceOnUse">
                    <stop stop-color="#cc1953" offset="0" />
                    <stop stop-color="#241284" offset="1" />
                </linearGradient>
                <linearGradient id="i" x1="79.6" x2="101.55" y1="91.55" y2="91.55" gradientTransform="matrix(2.4611 0 0 2.4611 0 -77.943)" gradientUnits="userSpaceOnUse">
                    <stop stop-color="#222b90" offset="0" />
                    <stop stop-color="#3559a6" offset="1" />
                </linearGradient>
                <linearGradient id="c" x1="52.22" x2="75.9" y1="77.19" y2="70.27" gradientTransform="matrix(2.4611 0 0 2.4611 0 -77.943)" gradientUnits="userSpaceOnUse">
                    <stop stop-color="#cc1953" offset="0" />
                    <stop stop-color="#241284" offset="1" />
                </linearGradient>
                <linearGradient id="f" x1="25.76" x2="49.27" y1="69.45" y2="78.01" gradientTransform="matrix(2.4611 0 0 2.4611 0 -77.943)" gradientUnits="userSpaceOnUse">
                    <stop stop-color="#fc0" offset=".0026698" />
                    <stop stop-color="#ce1938" offset=".9999" />
                </linearGradient>
                <linearGradient id="d" x1="28.39" x2="46.65" y1="47.49" y2="64.33" gradientTransform="matrix(2.4611 0 0 2.4611 0 -77.943)" gradientUnits="userSpaceOnUse">
                    <stop stop-color="#ffd100" offset=".0026698" />
                    <stop stop-color="#fd5a22" offset="1" />
                </linearGradient>
            </defs>
            <path d="m250 0v31.625h-54.046v-31.625h54.046z" fill="#40ab5d" stroke-width="2.4611" />
            <path d="m250 43.832v31.625h-54.046v-31.625h54.046z" fill="url(#e)" stroke-width="2.4611" />
            <path d="m250 87.69v31.625h-54.046v-31.625h54.046z" fill="url(#a)" stroke-width="2.4611" />
            <path d="m54.046 131.55v31.625h-54.046v-31.625h54.046z" fill="url(#b)" stroke-width="2.4611" />
            <path d="m119.34 131.55v31.625h-54.071v-31.625h54.071z" fill="url(#h)" stroke-width="2.4611" />
            <path d="m184.73 131.55v31.625h-54.071v-31.625h54.071z" fill="url(#g)" stroke-width="2.4611" />
            <path d="m250 131.55v31.625h-54.046v-31.625h54.046z" fill="url(#i)" stroke-width="2.4611" />
            <path d="m184.73 87.69v31.625h-54.071v-31.625h54.071z" fill="url(#c)" stroke-width="2.4611" />
            <path d="m119.34 87.69v31.625h-54.071v-31.625h54.071z" fill="url(#f)" stroke-width="2.4611" />
            <path d="m119.34 43.832v31.625h-54.071v-31.625h54.071z" fill="url(#d)" stroke-width="2.4611" />
        </svg>
    );
} 