import { useNavigate, useSearchParams } from "@solidjs/router";
import { createSignal, onMount, Show } from "solid-js";
import { apiService, ApiServiceError } from "~/api";
import AuthLayout from "~/components/auth/AuthLayout";
import Button from "~/components/ui/Button";
import Icon from "~/components/ui/Icon";
import LoadingSpinner from "~/components/ui/LoadingSpinner";

export default function VerifyEmail() {
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();
    const [isLoading, setIsLoading] = createSignal(true);
    const [error, setError] = createSignal<string | null>(null);
    const [success, setSuccess] = createSignal(false);

    onMount(async () => {
        const tokenParam = searchParams.token;
        const token = Array.isArray(tokenParam) ? tokenParam[0] : tokenParam;

        if (!token) {
            setError("Invalid verification link");
            setIsLoading(false);
            return;
        }

        try {
            await apiService.verifyEmail({ token });
            setSuccess(true);
        } catch (error) {
            console.error("Email verification error:", error);
            if (error instanceof ApiServiceError) {
                setError(error.message);
            } else {
                setError("Something went wrong. Please try again.");
            }
        } finally {
            setIsLoading(false);
        }
    });

    const handleContinue = () => {
        navigate("/dashboard");
    };

    const content = (
        <div class="text-center space-y-6">
            <Show when={isLoading()}>
                <div class="space-y-4">
                    <LoadingSpinner />
                    <p class="text-slate-600 dark:text-slate-400">
                        Verifying your email...
                    </p>
                </div>
            </Show>

            <Show when={success()}>
                <div class="space-y-4">
                    <div class="w-16 h-16 bg-emerald-100 dark:bg-emerald-900/30 rounded-full flex items-center justify-center mx-auto">
                        <Icon name="v-icon" class="w-8 h-8 text-emerald-600 dark:text-emerald-400" />
                    </div>
                    <h3 class="text-lg font-semibold text-slate-900 dark:text-slate-100">
                        Email Verified Successfully!
                    </h3>
                    <p class="text-sm text-slate-600 dark:text-slate-400">
                        Your email has been verified. You can now access all features of Syncer.
                    </p>
                    <Button
                        variant="primary"
                        size="lg"
                        onClick={handleContinue}
                        class="w-full justify-center"
                    >
                        Continue to Dashboard
                    </Button>
                </div>
            </Show>

            <Show when={error()}>
                <div class="space-y-4">
                    <div class="w-16 h-16 bg-red-100 dark:bg-red-900/30 rounded-full flex items-center justify-center mx-auto">
                        <Icon name="mail-icon" class="w-8 h-8 text-red-600 dark:text-red-400" />
                    </div>
                    <h3 class="text-lg font-semibold text-slate-900 dark:text-slate-100">
                        Verification Failed
                    </h3>
                    <p class="text-sm text-slate-600 dark:text-slate-400">
                        {error()}
                    </p>
                    <Button
                        variant="secondary"
                        size="lg"
                        onClick={() => navigate("/auth/login")}
                        class="w-full justify-center"
                    >
                        Back to Login
                    </Button>
                </div>
            </Show>
        </div>
    );

    return (
        <AuthLayout
            title="Email Verification"
            subtitle="We're verifying your email address"
        >
            {content}
        </AuthLayout>
    );
} 