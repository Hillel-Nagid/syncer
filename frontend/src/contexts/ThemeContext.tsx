import { createContext, createSignal, JSX, onMount, useContext } from "solid-js";

type Theme = "light" | "dark";

interface ThemeContextType {
    theme: () => Theme;
    toggleTheme: () => void;
}

const ThemeContext = createContext<ThemeContextType>();

export function ThemeProvider(props: { children: JSX.Element }) {
    const [theme, setTheme] = createSignal<Theme>("dark");

    onMount(() => {
        const savedTheme = localStorage.getItem("theme") as Theme;
        const systemPrefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;

        const initialTheme = savedTheme || (systemPrefersDark ? "dark" : "light");
        setTheme(initialTheme);
        applyTheme(initialTheme);
    });

    const applyTheme = (newTheme: Theme) => {
        const html = document.documentElement;
        if (newTheme === "dark") {
            html.classList.add("dark");
        } else {
            html.classList.remove("dark");
        }
        localStorage.setItem("theme", newTheme);
    };

    const toggleTheme = () => {
        const newTheme = theme() === "dark" ? "light" : "dark";
        setTheme(newTheme);
        applyTheme(newTheme);
    };

    const value: ThemeContextType = {
        theme,
        toggleTheme,
    };

    return (
        <ThemeContext.Provider value={value}>
            {props.children}
        </ThemeContext.Provider>
    );
}

export function useTheme() {
    const context = useContext(ThemeContext);
    if (!context) {
        throw new Error("useTheme must be used within a ThemeProvider");
    }
    return context;
} 