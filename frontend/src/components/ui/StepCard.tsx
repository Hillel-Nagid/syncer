import { JSX } from "solid-js";

interface StepCardProps {
    step: number;
    title: string;
    description: string;
    class?: string;
}

export default function StepCard(props: StepCardProps): JSX.Element {
    return (
        <div class={`text-center ${props.class || ""}`}>
            <div class="w-16 h-16 bg-emerald-100 dark:bg-emerald-900 rounded-full flex items-center justify-center mx-auto mb-4">
                <span class="text-2xl font-bold text-emerald-600 dark:text-emerald-400">{props.step}</span>
            </div>
            <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">{props.title}</h3>
            <p class="text-gray-600 dark:text-gray-400">
                {props.description}
            </p>
        </div>
    );
} 