import { createEffect, JSX, onCleanup, Show } from "solid-js";
import Icon from "../ui/Icon";

interface ModalProps {
    isOpen: () => boolean;
    onClose: () => void;
    title?: string;
    footer?: JSX.Element;
    size?: "sm" | "md" | "lg" | "xl";
    children: JSX.Element;
}

export default function Modal(props: ModalProps) {
    const size = () => props.size || "md";

    // Handle escape key
    const handleKeyDown = (e: KeyboardEvent) => {
        if (e.key === "Escape" && props.isOpen()) {
            props.onClose();
        }
    };

    // Prevent body scroll when modal is open
    createEffect(() => {
        if (props.isOpen()) {
            document.body.classList.add("overflow-hidden");
            document.addEventListener("keydown", handleKeyDown);
        } else {
            document.body.classList.remove("overflow-hidden");
            document.removeEventListener("keydown", handleKeyDown);
        }
    });

    onCleanup(() => {
        document.body.classList.remove("overflow-hidden");
        document.removeEventListener("keydown", handleKeyDown);
    });

    const sizeClasses = () => {
        switch (size()) {
            case "sm":
                return "max-w-md";
            case "md":
                return "max-w-lg";
            case "lg":
                return "max-w-2xl";
            case "xl":
                return "max-w-4xl";
            default:
                return "max-w-lg";
        }
    };

    // Custom scrollbar styles
    const scrollbarStyles = `
        .modal-scrollbar::-webkit-scrollbar {
            width: 8px;
        }
        
        .modal-scrollbar::-webkit-scrollbar-track {
            background: transparent;
            border-radius: 4px;
        }
        
        .modal-scrollbar::-webkit-scrollbar-thumb {
            background: rgba(16, 185, 129, 0.3);
            border-radius: 4px;
            transition: all 0.2s ease;
        }
        
        .modal-scrollbar::-webkit-scrollbar-thumb:hover {
            background: rgba(16, 185, 129, 0.5);
        }
        
        .modal-scrollbar::-webkit-scrollbar-thumb:active {
            background: rgba(16, 185, 129, 0.7);
        }
        
        .dark .modal-scrollbar::-webkit-scrollbar-thumb {
            background: rgba(52, 211, 153, 0.4);
        }
        
        .dark .modal-scrollbar::-webkit-scrollbar-thumb:hover {
            background: rgba(52, 211, 153, 0.6);
        }
        
        .dark .modal-scrollbar::-webkit-scrollbar-thumb:active {
            background: rgba(52, 211, 153, 0.8);
        }
        
        /* Firefox scrollbar */
        .modal-scrollbar {
            scrollbar-width: thin;
            scrollbar-color: rgba(16, 185, 129, 0.3) transparent;
        }
        
        .dark .modal-scrollbar {
            scrollbar-color: rgba(52, 211, 153, 0.4) transparent;
        }
    `;

    return (
        <Show when={props.isOpen()}>
            {/* Inject custom scrollbar styles */}
            <style>{scrollbarStyles}</style>

            {/* Backdrop */}
            <div
                class={`fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm transition-all duration-300 ${props.isOpen() ? "opacity-100" : "opacity-0"
                    }`}
                onClick={props.onClose}
            >
                {/* Modal Container */}
                <div
                    class={`relative w-full ${sizeClasses()} max-h-[90vh] bg-white dark:bg-gray-800 rounded-xl shadow-2xl border border-gray-200 dark:border-gray-700 transform transition-all duration-300 flex flex-col ${props.isOpen() ? "scale-100 opacity-100" : "scale-95 opacity-0"
                        }`}
                    onClick={(e) => e.stopPropagation()}
                >
                    {/* Header - Always Visible */}
                    <Show when={props.title}>
                        <div class="flex items-center justify-between p-6 border-b border-gray-200 dark:border-gray-700 flex-shrink-0 bg-white dark:bg-gray-800 rounded-t-xl">
                            <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
                                {props.title}
                            </h2>
                            <button
                                type="button"
                                onClick={props.onClose}
                                class="p-2 text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
                            >
                                <span class="sr-only">Close modal</span>
                                <Icon name="plus-icon" class="w-5 h-5 transform rotate-45" />
                            </button>
                        </div>
                    </Show>

                    {/* Scrollable Content */}
                    <div class="flex-1 overflow-y-auto p-6 modal-scrollbar">
                        {props.children}
                    </div>

                    {/* Footer - Always Visible */}
                    <Show when={props.footer}>
                        <div class="border-t border-gray-200 dark:border-gray-700 p-6 flex-shrink-0 bg-white dark:bg-gray-800 rounded-b-xl">
                            {props.footer}
                        </div>
                    </Show>
                </div>
            </div>
        </Show>
    );
} 