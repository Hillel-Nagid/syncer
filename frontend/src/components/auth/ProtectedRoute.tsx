import { useNavigate } from "@solidjs/router";
import { createSignal, JSX, onMount, Show } from "solid-js";
import { apiService } from "~/api";
import { useUser } from "~/contexts/UserContext";
import LoadingSpinner from "../ui/LoadingSpinner";

interface ProtectedRouteProps {
    children: JSX.Element;
    fallback?: JSX.Element;
}

export default function ProtectedRoute(props: ProtectedRouteProps) {
    const { isLoggedIn, isLoading } = useUser();
    const navigate = useNavigate();
    const [validatingSession, setValidatingSession] = createSignal(false);
    const [sessionValid, setSessionValid] = createSignal(false);

    onMount(async () => {
        if (isLoggedIn() && !isLoading()) {
            setValidatingSession(true);
            try {
                await apiService.getProfile();
                setSessionValid(true);
            } catch (error) {
                console.warn("Session validation failed:", error);
                setSessionValid(false);
                navigate("/auth/login", { replace: true });
            } finally {
                setValidatingSession(false);
            }
        }
    });

    if (isLoading() || validatingSession()) {
        return (
            <div class="flex items-center justify-center min-h-screen">
                <LoadingSpinner text="Loading..." />
            </div>
        );
    }

    if (!isLoggedIn()) {
        navigate("/auth/login", { replace: true });
        return null;
    }

    if (!sessionValid() && !validatingSession()) {
        return null;
    }

    return (
        <Show when={isLoggedIn() && sessionValid()} fallback={props.fallback}>
            {props.children}
        </Show>
    );
} 