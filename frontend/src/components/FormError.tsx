interface FormErrorProps {
    message: string;
}

export default function FormError(props: FormErrorProps) {
    return (
        <div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-600 dark:text-red-400 px-4 py-3 rounded-lg text-sm">
            {props.message}
        </div>
    );
} 