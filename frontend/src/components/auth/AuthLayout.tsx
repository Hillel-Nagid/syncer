import { JSX } from "solid-js";
import SyncIcon from "../SyncIcon";

interface AuthLayoutProps {
    title: string;
    subtitle: string;
    children: JSX.Element;
    bottomContent?: JSX.Element;
}

export default function AuthLayout(props: AuthLayoutProps) {
    return (
        <div class="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8 bg-gradient-to-br from-emerald-50 to-teal-50 dark:from-slate-900 dark:to-slate-800">
            <div class="max-w-md w-full space-y-8">
                {/* Header */}
                <div class="text-center">
                    <div class="flex justify-center mb-6">
                        <div class="w-16 h-16 bg-emerald-600 rounded-2xl flex items-center justify-center">
                            <SyncIcon size="md" class="text-white" />
                        </div>
                    </div>
                    <h2 class="text-3xl font-bold text-slate-900 dark:text-slate-100 mb-2">
                        {props.title}
                    </h2>
                    <p class="text-slate-600 dark:text-slate-400 text-sm">
                        {props.subtitle}
                    </p>
                </div>

                {/* Main Content */}
                <div class="bg-white dark:bg-slate-800 shadow-xl rounded-2xl px-8 py-8 border border-slate-200 dark:border-slate-700">
                    {props.children}
                </div>

                {/* Bottom Content */}
                {props.bottomContent && (
                    <div class="text-center">
                        {props.bottomContent}
                    </div>
                )}
            </div>
        </div>
    );
} 