export default function Footer() {
    return (
        <footer class="bg-slate-50 dark:bg-slate-900 text-slate-600 dark:text-slate-300">
            <div class="px-8 py-12 ">
                {/* Bottom Section */}
                <div class="border-t border-slate-200 dark:border-slate-700 flex flex-col md:flex-row justify-between items-center pt-8">
                    <p class="text-sm text-slate-500 dark:text-slate-400">
                        © 2025 Syncer. All rights reserved.
                    </p>
                    <div class="flex space-x-6 mt-4 md:mt-0">
                        <a href="#" class="text-sm text-slate-500 dark:text-slate-400 hover:text-emerald-600 dark:hover:text-emerald-400 transition-colors">
                            Privacy
                        </a>
                        <a href="#" class="text-sm text-slate-500 dark:text-slate-400 hover:text-emerald-600 dark:hover:text-emerald-400 transition-colors">
                            Terms
                        </a>
                        <a href="#" class="text-sm text-slate-500 dark:text-slate-400 hover:text-emerald-600 dark:hover:text-emerald-400 transition-colors">
                            Cookies
                        </a>
                    </div>
                </div>
            </div>
        </footer>
    );
} 