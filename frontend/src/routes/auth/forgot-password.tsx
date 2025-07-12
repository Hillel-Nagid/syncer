import { A } from "@solidjs/router";
import { createSignal, Show } from "solid-js";
import { apiService, ApiServiceError } from "~/api";
import AuthLayout from "~/components/auth/AuthLayout";
import FormError from "~/components/forms/FormError";
import FormInput from "~/components/forms/FormInput";
import Button from "~/components/ui/Button";
import Icon from "~/components/ui/Icon";
import LoadingSpinner from "~/components/ui/LoadingSpinner";

export default function ForgotPassword() {
    const [email, setEmail] = createSignal("");
    const [errors, setErrors] = createSignal<Record<string, string>>({});
    const [isLoading, setIsLoading] = createSignal(false);
    const [showSuccess, setShowSuccess] = createSignal(false);

    const handleSubmit = async (e: Event) => {
        e.preventDefault();

        if (!email().trim()) {
            setErrors({ email: "Email is required" });
            return;
        }

        if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email())) {
            setErrors({ email: "Please enter a valid email address" });
            return;
        }

        setIsLoading(true);
        setErrors({});

        try {
            await apiService.forgotPassword({ email: email() });
            setShowSuccess(true);
        } catch (error) {
            console.error("Password reset error:", error);
            if (error instanceof ApiServiceError) {
                setErrors({ general: error.message });
            } else {
                setErrors({ general: "Something went wrong. Please try again." });
            }
        } finally {
            setIsLoading(false);
        }
    };

    const formContent = (
        <Show
            when={!showSuccess()}
            fallback={
                <div class="text-center space-y-4">
                    <div class="w-16 h-16 bg-emerald-100 dark:bg-emerald-900/30 rounded-full flex items-center justify-center mx-auto">
                        <Icon name="mail-icon" class="w-8 h-8 text-emerald-600 dark:text-emerald-400" />
                    </div>
                    <h3 class="text-lg font-semibold text-slate-900 dark:text-slate-100">
                        Check your email
                    </h3>
                    <p class="text-sm text-slate-600 dark:text-slate-400">
                        We've sent password reset instructions to {email()}
                    </p>
                    <div class="pt-4">
                        <A
                            href="/auth/login"
                            class="text-emerald-600 hover:text-emerald-700 dark:text-emerald-400 dark:hover:text-emerald-300 font-medium"
                        >
                            Back to login
                        </A>
                    </div>
                </div>
            }
        >
            <form class="space-y-6" onSubmit={handleSubmit}>
                {/* General Error */}
                <Show when={errors().general}>
                    <FormError message={errors().general} />
                </Show>

                <div class="text-center mb-6">
                    <p class="text-sm text-slate-600 dark:text-slate-400">
                        Enter your email address and we'll send you a link to reset your password.
                    </p>
                </div>

                {/* Email Field */}
                <FormInput
                    id="email"
                    label="Email Address"
                    type="email"
                    value={email()}
                    onInput={(value) => {
                        setEmail(value);
                        if (errors().email) {
                            setErrors(prev => ({ ...prev, email: "" }));
                        }
                    }}
                    placeholder="Enter your email"
                    error={errors().email}
                    required
                />

                {/* Submit Button */}
                <Button
                    type="submit"
                    variant="primary"
                    size="lg"
                    class="w-full justify-center"
                >
                    <Show when={isLoading()}>
                        <LoadingSpinner text="Sending..." />
                    </Show>
                    <Show when={!isLoading()}>
                        Send Reset Link
                    </Show>
                </Button>

                {/* Back to Login Link */}
                <div class="text-center">
                    <A
                        href="/auth/login"
                        class="text-sm text-emerald-600 hover:text-emerald-700 dark:text-emerald-400 dark:hover:text-emerald-300 font-medium"
                    >
                        Back to login
                    </A>
                </div>
            </form>
        </Show>
    );

    return (
        <AuthLayout
            title="Forgot your password?"
            subtitle="No worries, we'll send you reset instructions"
        >
            {formContent}
        </AuthLayout>
    );
} 