import { A } from "@solidjs/router";
import { Component, createSignal, For } from "solid-js";
import Button from "~/components/ui/Button";
import Card from "~/components/ui/Card";
import Container from "~/components/ui/Container";
import Hero from "~/components/ui/Hero";
import Icon from "~/components/ui/Icon";
import Page from "~/components/ui/Page";

// Reusable FAQ Item Component
interface FAQItemProps {
    question: string;
    answer: string;
    isOpen: boolean;
    onToggle: () => void;
}

const FAQItem: Component<FAQItemProps> = (props) => {
    return (
        <div class="border-b border-slate-200 dark:border-slate-700 last:border-b-0">
            <button
                type="button"
                class="w-full py-6 text-left flex items-center justify-between gap-4 hover:bg-slate-50 dark:hover:bg-slate-700/50 transition-colors duration-200 rounded-lg px-2"
                onClick={props.onToggle}
            >
                <span class="text-lg font-semibold text-slate-800 dark:text-slate-200 pr-4">
                    {props.question}
                </span>
                <div class={`transform transition-transform duration-200 flex-shrink-0 ${props.isOpen ? 'rotate-45' : 'rotate-0'}`}>
                    <Icon name="plus-icon" class="w-5 h-5 text-emerald-600 dark:text-emerald-400" />
                </div>
            </button>
            <div class={`overflow-hidden transition-all duration-300 ease-in-out ${props.isOpen ? 'max-h-96 pb-6' : 'max-h-0'}`}>
                <div class="px-2 text-slate-600 dark:text-slate-300 leading-relaxed">
                    {props.answer}
                </div>
            </div>
        </div>
    );
};

// FAQ Category Component
interface FAQCategoryProps {
    title: string;
    faqs: Array<{ question: string; answer: string; }>;
}

const FAQCategory: Component<FAQCategoryProps> = (props) => {
    const [openItems, setOpenItems] = createSignal<Set<number>>(new Set());

    const toggleItem = (index: number) => {
        setOpenItems(prev => {
            const newSet = new Set(prev);
            if (newSet.has(index)) {
                newSet.delete(index);
            } else {
                newSet.add(index);
            }
            return newSet;
        });
    };

    return (
        <Card class="mb-8">
            <h2 class="text-2xl font-bold text-emerald-800 dark:text-emerald-200 mb-6">
                {props.title}
            </h2>
            <div class="space-y-0">
                <For each={props.faqs}>
                    {(faq, index) => (
                        <FAQItem
                            question={faq.question}
                            answer={faq.answer}
                            isOpen={openItems().has(index())}
                            onToggle={() => toggleItem(index())}
                        />
                    )}
                </For>
            </div>
        </Card>
    );
};

const FAQ: Component = () => {
    const generalFAQs = [
        {
            question: "What is Syncer and how does it work?",
            answer: "Syncer is a powerful synchronization engine that connects your calendars, music streaming services, and other digital platforms. It uses secure APIs to sync data between services in real-time, ensuring all your information stays consistent across platforms without manual intervention."
        },
        {
            question: "Which services does Syncer support?",
            answer: "Syncer supports major calendar platforms including Google Calendar, Outlook, Apple Calendar, and CalDAV-compatible services. For music streaming, we support Spotify, Apple Music, YouTube Music, and Amazon Music. We're continuously adding support for more services based on user demand."
        },
        {
            question: "Is my data secure with Syncer?",
            answer: "Absolutely. We use industry-standard encryption for all data transmission and storage. We only access the minimum data necessary for synchronization and never store sensitive personal information. All connections use OAuth 2.0 or similar secure authentication protocols."
        },
        {
            question: "How much does Syncer cost?",
            answer: "Syncer offers a free tier with basic synchronization for up to 2 services. Our Pro plan at $9.99/month includes unlimited services, advanced sync rules, and priority support. Enterprise plans are available for teams and organizations with custom pricing."
        }
    ];

    const calendarFAQs = [
        {
            question: "Can I sync calendars across different platforms?",
            answer: "Yes! Syncer can synchronize events between Google Calendar, Outlook, Apple Calendar, and other CalDAV-compatible services. You can choose which calendars to sync and set up rules for how events should be handled across platforms."
        },
        {
            question: "Will Syncer duplicate my calendar events?",
            answer: "No, Syncer uses intelligent deduplication to prevent duplicate events. Our system identifies matching events across platforms and maintains a single synchronized version, even when the same event exists in multiple calendars."
        },
        {
            question: "Can I set up bidirectional calendar sync?",
            answer: "Yes, you can configure bidirectional synchronization so that events created or modified in any connected calendar will be reflected across all your other calendars. You can also set up unidirectional sync if you prefer."
        },
        {
            question: "What happens to recurring events during sync?",
            answer: "Syncer fully supports recurring events and their patterns. When syncing recurring events, we preserve the recurrence rules and handle exceptions (like moved or cancelled instances) correctly across all platforms."
        }
    ];

    const musicFAQs = [
        {
            question: "Can I sync my playlists between Spotify and Apple Music?",
            answer: "Yes! Syncer can synchronize playlists between supported music streaming services. We match songs across platforms using metadata like artist, title, and album information to find the best matches on each service."
        },
        {
            question: "What if a song isn't available on all platforms?",
            answer: "When a song isn't available on a target platform, Syncer will mark it as unavailable and continue syncing the rest of the playlist. You'll receive a report of any songs that couldn't be matched or aren't available."
        },
        {
            question: "Can I sync my listening history and favorites?",
            answer: "Yes, Syncer can synchronize your liked songs, favorites, and listening history across platforms where the APIs permit. This helps maintain your music preferences consistently across all your streaming services."
        },
        {
            question: "How often does music sync happen?",
            answer: "Music synchronization can be set to run automatically at intervals you choose (hourly, daily, weekly) or triggered manually. Real-time sync is available for premium users where supported by the streaming service APIs."
        }
    ];

    const technicalFAQs = [
        {
            question: "Do I need to keep Syncer running all the time?",
            answer: "No, Syncer runs in the cloud. Once you've connected your services and configured your sync settings, synchronization happens automatically on our servers. You only need to access the dashboard when you want to modify settings or check sync status."
        },
        {
            question: "What happens if a service is temporarily unavailable?",
            answer: "Syncer includes robust error handling and retry logic. If a service is temporarily unavailable, we'll queue the sync operations and retry them automatically once the service is back online. You'll be notified of any extended outages."
        },
        {
            question: "Can I set up custom sync rules?",
            answer: "Yes! Pro users can create custom sync rules, including filtering by keywords, date ranges, or specific calendars/playlists. You can also set up conditional rules like only syncing work events during business hours."
        },
        {
            question: "How do I disconnect a service from Syncer?",
            answer: "You can easily disconnect any service from your Syncer dashboard. Go to Settings > Connected Services, and click 'Disconnect' next to the service you want to remove. This will revoke Syncer's access to that service immediately."
        }
    ];

    return (
        <Page class="py-12 px-4 sm:px-6 lg:px-8">
            <Container>
                <Hero
                    title="Frequently Asked Questions"
                    subtitle="Find answers to common questions about Syncer's synchronization capabilities, security, and supported services."
                />

                <div class="mt-12">
                    <FAQCategory title="General Questions" faqs={generalFAQs} />
                    <FAQCategory title="Calendar Synchronization" faqs={calendarFAQs} />
                    <FAQCategory title="Music Streaming" faqs={musicFAQs} />
                    <FAQCategory title="Technical & Setup" faqs={technicalFAQs} />
                </div>

                {/* Contact Section */}
                <Card class="mt-12 text-center">
                    <h2 class="text-2xl font-bold text-emerald-800 dark:text-emerald-200 mb-4">
                        Still have questions?
                    </h2>
                    <p class="text-slate-600 dark:text-slate-300 mb-6">
                        Our support team is here to help you get the most out of Syncer.
                    </p>
                    <div class="flex flex-col sm:flex-row gap-4 justify-center">
                        <A href="/contact">
                            <Button variant="primary" size="lg">
                                Contact us
                            </Button>
                        </A>
                    </div>
                </Card>
            </Container>
        </Page>
    );
};

export default FAQ;
