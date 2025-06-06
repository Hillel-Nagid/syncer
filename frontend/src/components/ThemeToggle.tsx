import { useTheme } from "~/contexts/ThemeContext";
import Icon from "./Icon";

export default function ThemeToggle() {
    const { theme, toggleTheme } = useTheme();

    return (
        <button
            onClick={toggleTheme}
            type="button"
            class="relative p-2 rounded-lg text-slate-600 dark:text-slate-400 hover:text-emerald-600 hover:dark:text-emerald-400 hover:bg-slate-100 dark:hover:bg-slate-800 focus:outline-none focus:ring-2 focus:ring-emerald-500 transition-colors"
            aria-label={`Switch to ${theme() === "dark" ? "light" : "dark"} mode`}
            title={`Switch to ${theme() === "dark" ? "light" : "dark"} mode`}
        >
            <div class="relative w-5 h-5">
                <div
                    class={`absolute inset-0 transition-all duration-300 ${theme() === "light"
                        ? "rotate-0 scale-100 opacity-100"
                        : "rotate-90 scale-0 opacity-0"
                        }`}
                >
                    <Icon name="sun-icon" class="w-5 h-5" alt="Sun Icon" />
                </div>

                <div
                    class={`absolute inset-0 transition-all duration-300 ${theme() === "dark"
                        ? "rotate-0 scale-100 opacity-100"
                        : "-rotate-90 scale-0 opacity-0"
                        }`}
                >
                    <Icon name="moon-icon" class="w-5 h-5" alt="Moon Icon" />
                </div>
            </div>
        </button>
    );
} 