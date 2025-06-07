import { JSX } from "solid-js";

interface HeroProps {
    title: string;
    subtitle: string;
    class?: string;
    titleClass?: string;
    subtitleClass?: string;
}

export default function Hero(props: HeroProps): JSX.Element {
    const defaultTitleClass = "text-5xl md:text-6xl font-bold text-emerald-600 dark:text-emerald-400 mb-6";
    const defaultSubtitleClass = "text-xl md:text-2xl text-slate-700 dark:text-slate-300 mb-8 max-w-3xl mx-auto";

    return (
        <div class={`text-center ${props.class || ""}`}>
            <h1 class={props.titleClass || defaultTitleClass}>
                {props.title}
            </h1>
            <p class={props.subtitleClass || defaultSubtitleClass}>
                {props.subtitle}
            </p>
        </div>
    );
} 