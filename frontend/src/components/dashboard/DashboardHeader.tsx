import Button from "~/components/ui/Button";
import Icon from "~/components/ui/Icon";

interface DashboardHeaderProps {
    title: string;
    description: string;
    addButtonText: string;
    onBack: () => void;
    onAdd: () => void;
}

export default function DashboardHeader(props: DashboardHeaderProps) {
    return (
        <div class="flex items-center justify-between mb-8">
            <div class="flex items-center">
                <Button
                    variant="secondary"
                    size="sm"
                    onClick={props.onBack}
                    class="mr-4"
                >
                    <Icon name="arrow-left" class="w-4 h-4 mr-2" />
                    Back
                </Button>
                <div>
                    <h1 class="text-3xl md:text-4xl font-bold text-gray-900 dark:text-white">
                        {props.title}
                    </h1>
                    <p class="text-lg text-gray-600 dark:text-gray-300 mt-2">
                        {props.description}
                    </p>
                </div>
            </div>
            <Button
                variant="primary"
                size="md"
                onClick={props.onAdd}
            >
                <Icon name="plus-icon" class="w-5 h-5 mr-2" />
                {props.addButtonText}
            </Button>
        </div>
    );
} 