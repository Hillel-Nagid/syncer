
interface FormInputProps {
    id: string;
    label: string;
    type?: string;
    value: string;
    onInput: (value: string) => void;
    placeholder?: string;
    error?: string;
    required?: boolean;
}

export default function FormInput(props: FormInputProps) {
    return (
        <div>
            <label
                for={props.id}
                class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2"
            >
                {props.label}
                {props.required && <span class="text-red-500 ml-1">*</span>}
            </label>
            <input
                id={props.id}
                type={props.type || "text"}
                value={props.value}
                onInput={(e) => props.onInput(e.currentTarget.value)}
                class={`w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 transition-colors bg-white dark:bg-slate-700 text-slate-900 dark:text-slate-100 ${props.error
                    ? "border-red-300 dark:border-red-600"
                    : "border-slate-300 dark:border-slate-600"
                    }`}
                placeholder={props.placeholder}
            />
            {props.error && (
                <p class="mt-1 text-sm text-red-600 dark:text-red-400">{props.error}</p>
            )}
        </div>
    );
} 