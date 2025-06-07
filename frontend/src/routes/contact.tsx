import { createSignal } from "solid-js";
import FormInput from "~/components/forms/FormInput";
import FormTextarea from "~/components/forms/FormTextarea";
import Button from "~/components/ui/Button";
import Card from "~/components/ui/Card";
import Container from "~/components/ui/Container";
import Icon from "~/components/ui/Icon";
import Page from "~/components/ui/Page";

interface ContactFormData {
    name: string;
    email: string;
    subject: string;
    message: string;
}

interface ContactFormErrors {
    name?: string;
    email?: string;
    subject?: string;
    message?: string;
    general?: string;
}

export default function Contact() {
    const [formData, setFormData] = createSignal<ContactFormData>({
        name: "",
        email: "",
        subject: "",
        message: "",
    });

    const [errors, setErrors] = createSignal<ContactFormErrors>({});
    const [isSubmitting, setIsSubmitting] = createSignal(false);
    const [isSubmitted, setIsSubmitted] = createSignal(false);

    const handleInputChange = (field: keyof ContactFormData, value: string) => {
        setFormData(prev => ({ ...prev, [field]: value }));
        // Clear error when user starts typing
        if (errors()[field]) {
            setErrors(prev => ({ ...prev, [field]: undefined }));
        }
    };

    const validateForm = (data: ContactFormData): ContactFormErrors => {
        const newErrors: ContactFormErrors = {};

        if (!data.name.trim()) {
            newErrors.name = "Name is required";
        }

        if (!data.email.trim()) {
            newErrors.email = "Email is required";
        } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(data.email)) {
            newErrors.email = "Please enter a valid email address";
        }

        if (!data.subject.trim()) {
            newErrors.subject = "Subject is required";
        }

        if (!data.message.trim()) {
            newErrors.message = "Message is required";
        } else if (data.message.trim().length < 10) {
            newErrors.message = "Message must be at least 10 characters long";
        }

        return newErrors;
    };

    const handleSubmit = async (e: Event) => {
        e.preventDefault();
        const data = formData();
        const validationErrors = validateForm(data);

        if (Object.keys(validationErrors).length > 0) {
            setErrors(validationErrors);
            return;
        }

        setIsSubmitting(true);
        setErrors({});

        try {
            // Here you would typically send the data to your backend
            // For now, we'll simulate a successful submission
            await new Promise(resolve => setTimeout(resolve, 1000));

            setIsSubmitted(true);
            setFormData({
                name: "",
                email: "",
                subject: "",
                message: "",
            });
        } catch (error) {
            setErrors({ general: "Something went wrong. Please try again." });
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <Page>
            <Container maxWidth="4xl" class="px-6 py-12">
                <div class="text-center mb-12">
                    <h1 class="text-4xl font-bold text-slate-900 dark:text-slate-100 mb-4">
                        Get in Touch
                    </h1>
                    <p class="text-xl text-slate-600 dark:text-slate-300 max-w-2xl mx-auto">
                        Have questions about Syncer? Need help with integration?
                        We'd love to hear from you. Send us a message and we'll respond as soon as possible.
                    </p>
                </div>

                <div class="flex flex-col gap-8">


                    <div class="lg:col-span-2">
                        <Card>
                            <div class="p-8">
                                {isSubmitted() ? (
                                    <div class="text-center py-8">
                                        <div class="w-16 h-16 bg-emerald-100 dark:bg-emerald-900/20 rounded-full flex items-center justify-center mx-auto mb-4">
                                            <Icon name="v-icon" class="w-8 h-8 text-emerald-600 dark:text-emerald-400" />
                                        </div>
                                        <h3 class="text-2xl font-semibold text-slate-900 dark:text-slate-100 mb-2">
                                            Message Sent!
                                        </h3>
                                        <p class="text-slate-600 dark:text-slate-400 mb-6">
                                            Thank you for contacting us. We'll get back to you within 24 hours.
                                        </p>
                                        <Button onClick={() => setIsSubmitted(false)}>
                                            Send Another Message
                                        </Button>
                                    </div>
                                ) : (
                                    <>
                                        <h2 class="text-2xl font-semibold text-slate-900 dark:text-slate-100 mb-6">
                                            Send us a Message
                                        </h2>

                                        {errors().general && (
                                            <div class="mb-6 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
                                                <p class="text-red-700 dark:text-red-400">{errors().general}</p>
                                            </div>
                                        )}

                                        <form onSubmit={handleSubmit} class="space-y-6">
                                            <div class="grid md:grid-cols-2 gap-6">
                                                <FormInput
                                                    id="name"
                                                    label="Full Name"
                                                    type="text"
                                                    value={formData().name}
                                                    onInput={(value) => handleInputChange("name", value)}
                                                    placeholder="Your full name"
                                                    error={errors().name}
                                                    required
                                                />
                                                <FormInput
                                                    id="email"
                                                    label="Email Address"
                                                    type="email"
                                                    value={formData().email}
                                                    onInput={(value) => handleInputChange("email", value)}
                                                    placeholder="your.email@example.com"
                                                    error={errors().email}
                                                    required
                                                />
                                            </div>

                                            <FormInput
                                                id="subject"
                                                label="Subject"
                                                type="text"
                                                value={formData().subject}
                                                onInput={(value) => handleInputChange("subject", value)}
                                                placeholder="What is your message about?"
                                                error={errors().subject}
                                                required
                                            />

                                            <FormTextarea
                                                id="message"
                                                label="Message"
                                                value={formData().message}
                                                onInput={(value) => handleInputChange("message", value)}
                                                placeholder="Tell us more about your question or feedback..."
                                                error={errors().message}
                                                required
                                                rows={6}
                                            />

                                            <div class="flex justify-center">
                                                <Button
                                                    type="submit"
                                                    icon="mail-icon"
                                                    size="lg"
                                                    class={isSubmitting() ? "opacity-50 cursor-not-allowed" : ""}
                                                >
                                                    {isSubmitting() ? "Sending..." : "Send Message"}
                                                </Button>
                                            </div>
                                        </form>
                                    </>
                                )}
                            </div>
                        </Card>
                    </div>
                </div>
            </Container>
        </Page>
    );
}
