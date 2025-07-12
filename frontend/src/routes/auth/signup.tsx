import { A, useNavigate } from "@solidjs/router";
import { createSignal, Show } from "solid-js";
import { apiService, ApiServiceError } from "~/api";
import AuthLayout from "~/components/auth/AuthLayout";
import FormError from "~/components/forms/FormError";
import FormInput from "~/components/forms/FormInput";
import Button from "~/components/ui/Button";
import Icon from "~/components/ui/Icon";
import LoadingSpinner from "~/components/ui/LoadingSpinner";
import { useUser } from "~/contexts/UserContext";

export default function Signup() {
    const [formData, setFormData] = createSignal({
        name: "",
        email: "",
        password: "",
        confirmPassword: "",
    });
    const [errors, setErrors] = createSignal<Record<string, string>>({});
    const [isLoading, setIsLoading] = createSignal(false);
    const [showSuccess, setShowSuccess] = createSignal(false);

    const { login } = useUser();
    const navigate = useNavigate();

    const handleInputChange = (field: string, value: string) => {
        setFormData(prev => ({ ...prev, [field]: value }));
        // Clear error when user starts typing
        if (errors()[field]) {
            setErrors(prev => ({ ...prev, [field]: "" }));
        }
    };

    const validateForm = () => {
        const data = formData();
        const newErrors: Record<string, string> = {};

        if (!data.name.trim()) {
            newErrors.name = "Name is required";
        }

        if (!data.email.trim()) {
            newErrors.email = "Email is required";
        } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(data.email)) {
            newErrors.email = "Please enter a valid email address";
        }

        if (!data.password) {
            newErrors.password = "Password is required";
        } else if (data.password.length < 8) {
            newErrors.password = "Password must be at least 8 characters";
        }

        if (!data.confirmPassword) {
            newErrors.confirmPassword = "Please confirm your password";
        } else if (data.password !== data.confirmPassword) {
            newErrors.confirmPassword = "Passwords do not match";
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = async (e: Event) => {
        e.preventDefault();

        if (!validateForm()) {
            return;
        }

        setIsLoading(true);
        setErrors({});

        try {
            const response = await apiService.register({
                email: formData().email,
                password: formData().password,
                full_name: formData().name,
            });

            login(response);
            setShowSuccess(true);

            // Show success message for 2 seconds, then redirect
            setTimeout(() => {
                navigate("/dashboard");
            }, 2000);
        } catch (error) {
            console.error("Signup error:", error);
            if (error instanceof ApiServiceError) {
                if (error.status === 409) {
                    setErrors({ email: "An account with this email already exists." });
                } else {
                    setErrors({ general: error.message });
                }
            } else {
                setErrors({ general: "Something went wrong. Please try again." });
            }
        } finally {
            setIsLoading(false);
        }
    };

    const handleGoogleSignup = async () => {
        try {
            const response = await apiService.getGoogleAuthUrl();
            window.location.href = response.auth_url;
        } catch (error) {
            console.error("Google signup error:", error);
            setErrors({ general: "Failed to initialize Google signup. Please try again." });
        }
    };

    const formContent = (
        <Show
            when={!showSuccess()}
            fallback={
                <div class="text-center space-y-4">
                    <div class="w-16 h-16 bg-emerald-100 dark:bg-emerald-900/30 rounded-full flex items-center justify-center mx-auto">
                        <Icon name="v-icon" class="w-8 h-8 text-emerald-600 dark:text-emerald-400" />
                    </div>
                    <h3 class="text-lg font-semibold text-slate-900 dark:text-slate-100">
                        Account Created Successfully!
                    </h3>
                    <p class="text-sm text-slate-600 dark:text-slate-400">
                        Please check your email to verify your account. You'll be redirected to the dashboard shortly.
                    </p>
                </div>
            }
        >
            <form class="space-y-6" onSubmit={handleSubmit}>
                {/* General Error */}
                <Show when={errors().general}>
                    <FormError message={errors().general} />
                </Show>

                {/* Name Field */}
                <FormInput
                    id="name"
                    label="Full Name"
                    type="text"
                    value={formData().name}
                    onInput={(value) => handleInputChange("name", value)}
                    placeholder="Enter your full name"
                    error={errors().name}
                    required
                />

                {/* Email Field */}
                <FormInput
                    id="email"
                    label="Email Address"
                    type="email"
                    value={formData().email}
                    onInput={(value) => handleInputChange("email", value)}
                    placeholder="Enter your email"
                    error={errors().email}
                    required
                />

                {/* Password Field */}
                <FormInput
                    id="password"
                    label="Password"
                    type="password"
                    value={formData().password}
                    onInput={(value) => handleInputChange("password", value)}
                    placeholder="Create a password"
                    error={errors().password}
                    required
                />

                {/* Confirm Password Field */}
                <FormInput
                    id="confirmPassword"
                    label="Confirm Password"
                    type="password"
                    value={formData().confirmPassword}
                    onInput={(value) => handleInputChange("confirmPassword", value)}
                    placeholder="Confirm your password"
                    error={errors().confirmPassword}
                    required
                />

                {/* Terms and Privacy */}
                <div class="text-sm text-slate-600 dark:text-slate-400">
                    By creating an account, you agree to our{" "}
                    <a href="/terms" class="text-emerald-600 hover:text-emerald-700 dark:text-emerald-400 dark:hover:text-emerald-300">
                        Terms of Service
                    </a>{" "}
                    and{" "}
                    <a href="/privacy" class="text-emerald-600 hover:text-emerald-700 dark:text-emerald-400 dark:hover:text-emerald-300">
                        Privacy Policy
                    </a>
                </div>

                {/* Submit Button */}
                <Button
                    type="submit"
                    variant="primary"
                    size="lg"
                    class="w-full justify-center"
                >
                    <Show when={isLoading()}>
                        <LoadingSpinner text="Creating account..." />
                    </Show>
                    <Show when={!isLoading()}>
                        Create Account
                    </Show>
                </Button>

                {/* Divider */}
                <div class="relative">
                    <div class="absolute inset-0 flex items-center">
                        <div class="w-full border-t border-slate-300 dark:border-slate-600" />
                    </div>
                    <div class="relative flex justify-center text-sm">
                        <span class="px-2 bg-white dark:bg-slate-800 text-slate-500 dark:text-slate-400">
                            Or continue with
                        </span>
                    </div>
                </div>

                {/* Google Signup Button */}
                <div class="flex justify-center">
                    <Button
                        variant="secondary"
                        size="md"
                        class="w-full max-w-sm justify-center"
                        icon="google-icon"
                        onClick={handleGoogleSignup}
                    >
                        Continue with Google
                    </Button>
                </div>

                {/* Login Link */}
                <div class="text-center">
                    <p class="text-sm text-slate-600 dark:text-slate-400">
                        Already have an account?{" "}
                        <A
                            href="/auth/login"
                            class="text-emerald-600 hover:text-emerald-700 dark:text-emerald-400 dark:hover:text-emerald-300 font-medium"
                        >
                            Sign in
                        </A>
                    </p>
                </div>
            </form>
        </Show>
    );

    const bottomContent = (
        <>
            <p class="text-sm text-slate-600 dark:text-slate-400 mb-4">
                Start synchronizing your services today
            </p>
            <div class="flex justify-center space-x-6">
                <div class="flex items-center text-xs text-slate-500 dark:text-slate-400">
                    <Icon name="v-icon" class="w-4 h-4 mr-1 text-emerald-500" />
                    Calendar Sync
                </div>
                <div class="flex items-center text-xs text-slate-500 dark:text-slate-400">
                    <Icon name="v-icon" class="w-4 h-4 mr-1 text-emerald-500" />
                    Music Streaming
                </div>
                <div class="flex items-center text-xs text-slate-500 dark:text-slate-400">
                    <Icon name="v-icon" class="w-4 h-4 mr-1 text-emerald-500" />
                    Real-time Sync
                </div>
            </div>
        </>
    );

    return (
        <AuthLayout
            title="Create your account"
            subtitle="Join Syncer to synchronize your services seamlessly"
            bottomContent={bottomContent}
        >
            {formContent}
        </AuthLayout>
    );
}
