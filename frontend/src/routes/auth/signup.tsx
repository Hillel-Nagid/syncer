import { A } from "@solidjs/router";
import { createSignal } from "solid-js";
import AuthLayout from "~/components/AuthLayout";
import Button from "~/components/Button";
import FormError from "~/components/FormError";
import FormInput from "~/components/FormInput";
import Icon from "~/components/Icon";
import LoadingSpinner from "~/components/LoadingSpinner";

export default function Signup() {
    const [formData, setFormData] = createSignal({
        name: "",
        email: "",
        password: "",
        confirmPassword: "",
    });
    const [errors, setErrors] = createSignal<Record<string, string>>({});
    const [isLoading, setIsLoading] = createSignal(false);

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

        try {
            // TODO: Implement actual signup API call
            await new Promise(resolve => setTimeout(resolve, 1000)); // Simulate API call
            console.log("Signup data:", formData());
            // Redirect to dashboard or login on success
        } catch (error) {
            console.error("Signup error:", error);
            setErrors({ general: "Something went wrong. Please try again." });
        } finally {
            setIsLoading(false);
        }
    };

    const formContent = (
        <form class="space-y-6" onSubmit={handleSubmit}>
            {/* General Error */}
            {errors().general && <FormError message={errors().general} />}

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
                {isLoading() ? <LoadingSpinner text="Creating account..." /> : "Create Account"}
            </Button>

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
