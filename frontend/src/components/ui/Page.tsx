import { JSX } from "solid-js";

interface PageProps {
    children: JSX.Element;
    class?: string;
}

export default function Page(props: PageProps) {
    return (
        <div class={`min-h-screen bg-gradient-to-br from-emerald-50 to-emerald-100 dark:from-slate-900 dark:to-slate-800  ${props.class || ""}`}>
            {props.children}
        </div>
    );
}