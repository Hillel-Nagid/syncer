import { createSignal, onCleanup, Show } from "solid-js";
import { useUser } from "~/contexts/UserContext";
import Icon from "../ui/Icon";

export default function UserProfile() {
    const { user, logout } = useUser();
    const [isDropdownOpen, setIsDropdownOpen] = createSignal(false);
    let dropdownRef: HTMLDivElement | undefined;

    const handleLogout = () => {
        logout();
        setIsDropdownOpen(false);
    };

    const handleClickOutside = (e: MouseEvent) => {
        const target = e.target as Node;
        if (isDropdownOpen() && dropdownRef && !dropdownRef.contains(target)) {
            setIsDropdownOpen(false);
        }
    };

    const toggleDropdown = (e: MouseEvent) => {
        e.stopPropagation();
        setIsDropdownOpen(!isDropdownOpen());
    };

    document.addEventListener('click', handleClickOutside);

    onCleanup(() => {
        document.removeEventListener('click', handleClickOutside);
    });


    return (
        <div class="relative" ref={dropdownRef}>
            <button
                type="button"
                class="flex items-center space-x-2 p-2 rounded-lg text-slate-600 dark:text-slate-400 hover:text-emerald-600 hover:dark:text-emerald-400 hover:bg-slate-100 dark:hover:bg-slate-800 focus:outline-none focus:ring-2 focus:ring-emerald-500 transition-colors"
                onClick={toggleDropdown}
            >
                <Show
                    when={user()?.profilePicture}
                    fallback={
                        <div class="w-8 h-8 bg-emerald-100 dark:bg-emerald-900/30 rounded-full flex items-center justify-center">
                            <Icon name="user-icon" class="w-5 h-5 text-emerald-600 dark:text-emerald-400" alt="User Icon" />
                        </div>
                    }
                >
                    <img
                        src={user()?.profilePicture}
                        alt={`${user()?.username}'s profile`}
                        class="w-8 h-8 rounded-full object-cover"
                    />
                </Show>
                <span class="text-sm font-medium hidden sm:block">{user()?.username}</span>
            </button>

            <Show when={isDropdownOpen()}>
                <div class="absolute right-0 mt-2 w-48 bg-white dark:bg-slate-800 rounded-lg shadow-lg border border-slate-200 dark:border-slate-700 z-50">
                    <div class="px-4 py-3 border-b border-slate-200 dark:border-slate-700">
                        <p class="text-sm font-medium text-slate-900 dark:text-slate-100">{user()?.username}</p>
                        <p class="text-xs text-slate-500 dark:text-slate-400">{user()?.email}</p>
                    </div>
                    <div class="py-1">
                        <button
                            type="button"
                            onClick={handleLogout}
                            class="w-full flex items-center px-4 py-2 text-sm text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 hover:text-emerald-600 hover:dark:text-emerald-400 transition-colors"
                        >
                            <Icon name="logout-icon" class="w-4 h-4 mr-3" alt="Logout Icon" />
                            Sign out
                        </button>
                    </div>
                </div>
            </Show>
        </div>
    );
} 