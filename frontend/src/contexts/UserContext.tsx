import { createContext, createEffect, createSignal, JSX, onMount, useContext } from "solid-js";
import { apiService } from "~/api";
import { AuthResponse, User } from "~/types";

interface UserContextType {
    user: () => User | null;
    isLoggedIn: () => boolean;
    login: (authResponse: AuthResponse) => void;
    logout: () => Promise<void>;
    isLoading: () => boolean;
}

const UserContext = createContext<UserContextType>();

export function UserProvider(props: { children: JSX.Element }) {
    const [user, setUser] = createSignal<User | null>(null);
    const [isLoading, setIsLoading] = createSignal(true);

    onMount(async () => {
        try {
            const profile = await apiService.getProfile();
            if (profile) {
                setUser(profile);
            } else {
                setUser(null);
            }
        } catch (error) {
            console.error("Error fetching profile:", error);
            setUser(null);
        } finally {
            setIsLoading(false);
        }
    });


    const isLoggedIn = () => user() !== null;

    const login = (authResponse: AuthResponse) => {
        setUser(authResponse.user);
    };

    const logout = async () => {
        try {
            await apiService.logout();
        } catch (error) {
            console.warn("Logout API call failed:", error);
        } finally {
            setUser(null);
        }
    };

    const value: UserContextType = {
        user,
        isLoggedIn,
        login,
        logout,
        isLoading,
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