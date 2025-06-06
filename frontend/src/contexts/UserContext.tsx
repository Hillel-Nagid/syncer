import { createContext, createSignal, JSX, onMount, useContext } from "solid-js";
import { User } from "~/types";

interface UserContextType {
    user: () => User | null;
    isLoggedIn: () => boolean;
    login: (userData: User) => void;
    logout: () => void;
}

const UserContext = createContext<UserContextType>();

export function UserProvider(props: { children: JSX.Element }) {
    const [user, setUser] = createSignal<User | null>(null);

    onMount(() => {
        // Check localStorage for existing user session
        const savedUser = localStorage.getItem("user");
        if (savedUser) {
            try {
                setUser(JSON.parse(savedUser));
            } catch (error) {
                console.error("Failed to parse saved user data:", error);
                localStorage.removeItem("user");
            }
        }
    });

    const isLoggedIn = () => user() !== null;

    const login = (userData: User) => {
        setUser(userData);
        localStorage.setItem("user", JSON.stringify(userData));
    };

    const logout = () => {
        setUser(null);
        localStorage.removeItem("user");
    };

    const value: UserContextType = {
        user,
        isLoggedIn,
        login,
        logout,
    };

    return (
        <UserContext.Provider value={value}>
            {props.children}
        </UserContext.Provider>
    );
}

export function useUser() {
    const context = useContext(UserContext);
    if (!context) {
        throw new Error("useUser must be used within a UserProvider");
    }
    return context;
} 