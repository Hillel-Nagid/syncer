import { useNavigate } from "@solidjs/router";
import { createSignal } from "solid-js";
import CalendarDashboard from "~/components/dashboard/CalendarDashboard";
import Page from "~/components/ui/Page";

export default function CalendarRoute() {
    const navigate = useNavigate();
    const [isTransitioning, setIsTransitioning] = createSignal(false);

    const handleBack = () => {
        setIsTransitioning(true);
        setTimeout(() => {
            navigate("/dashboard");
        }, 200);
    };

    return (
        <Page>
            <div class={`transition-opacity duration-300 ease-in-out ${isTransitioning() ? 'opacity-0' : 'opacity-100'}`}>
                <CalendarDashboard onBack={handleBack} />
            </div>
        </Page>
    );
} 