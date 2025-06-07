import { A } from "@solidjs/router";
import { createSignal } from "solid-js";
import AuthLayout from "~/components/auth/AuthLayout";
import Button from "~/components/ui/Button";
import FormError from "~/components/forms/FormError";
import FormInput from "~/components/forms/FormInput";
import Icon from "~/components/ui/Icon";
import LoadingSpinner from "~/components/ui/LoadingSpinner";

export default function Login() {
    const [formData, setFormData] = createSignal({
        email: "",
        password: "",
    });
    const [errors, setErrors] = createSignal<Record<string, string>>({});
    const [isLoading, setIsLoading] = createSignal(false);
    const [rememberMe, setRememberMe] = createSignal(false);

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

        if (!data.email.trim()) {
            newErrors.email = "Email is required";
        } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(data.email)) {
            newErrors.email = "Please enter a valid email address";
        }

        if (!data.password) {
            newErrors.password = "Password is required";
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

        try {
            // TODO: Implement actual login API call
            await new Promise(resolve => setTimeout(resolve, 1000)); // Simulate API call
            console.log("Login data:", { ...formData(), rememberMe: rememberMe() });
            // Redirect to dashboard on success
        } catch (error) {
            console.error("Login error:", error);
            setErrors({ general: "Invalid email or password. Please try again." });
        } finally {
            setIsLoading(false);
        }
    };

    const formContent = (
        <form class="space-y-6" onSubmit={handleSubmit}>
            {/* General Error */}
            {errors().general && <FormError message={errors().general} />}

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
                placeholder="Enter your password"
                error={errors().password}
                required
            />

            {/* Remember Me and Forgot Password */}
            <div class="flex items-center justify-between">
                <div class="flex items-center">
                    <input
                        id="remember-me"
                        type="checkbox"
                        checked={rememberMe()}
                        onChange={(e) => setRememberMe(e.currentTarget.checked)}
                        class="h-4 w-4 text-emerald-600 focus:ring-emerald-500 border-slate-300 dark:border-slate-600 rounded bg-white dark:bg-slate-700"
                    />
                    <label for="remember-me" class="ml-2 block text-sm text-slate-700 dark:text-slate-300">
                        Remember me
                    </label>
                </div>

                <div class="text-sm">
                    <a
                        href="/auth/forgot-password"
                        class="text-emerald-600 hover:text-emerald-700 dark:text-emerald-400 dark:hover:text-emerald-300 font-medium"
                    >
                        Forgot your password?
                    </a>
                </div>
            </div>

            {/* Submit Button */}
            <Button
                type="submit"
                variant="primary"
                size="lg"
                class="w-full justify-center"
            >
                {isLoading() ? <LoadingSpinner text="Signing in..." /> : "Sign In"}
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

            {/* Social Login Button */}
            <div class="flex justify-center">
                <Button
                    variant="secondary"
                    size="md"
                    class="w-full max-w-sm justify-center"
                    icon="google-icon"
                    onClick={() => console.log("Google login")}
                >
                    Continue with Google
                </Button>
            </div>

            {/* Sign Up Link */}
            <div class="text-center">
                <p class="text-sm text-slate-600 dark:text-slate-400">
                    Don't have an account?{" "}
                    <A
                        href="/auth/signup"
                        class="text-emerald-600 hover:text-emerald-700 dark:text-emerald-400 dark:hover:text-emerald-300 font-medium"
                    >
                        Sign up
                    </A>
                </p>
            </div>
        </form>
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
            title="Welcome back"
            subtitle="Sign in to your Syncer account"
            bottomContent={bottomContent}
        >
            {formContent}
        </AuthLayout>
    );
}
