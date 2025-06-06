import { A } from "@solidjs/router";
import { For } from "solid-js";

interface MobileMenuProps {
    isMobileMenuOpen: () => boolean;
    closeMenu: () => void;
    navLinks: { href: string, text: string }[];
}

export default function MobileMenu(props: MobileMenuProps) {
    return (
        < div class={`fixed inset-0 top-16 z-40 md:hidden transition-all duration-300 ease-in-out ${props.isMobileMenuOpen() ? "opacity-100 visible" : "opacity-0 invisible"}`}>
            {/* Overlay */}
            < div
                class="fixed inset-0 bg-black/30 backdrop-blur-sm transition-opacity duration-300"
                onClick={props.closeMenu}
            />

            {/* Mobile Nav Panel */}
            <div class={`relative bg-white dark:bg-slate-900 w-full p-6 border-t border-slate-200 dark:border-slate-800 shadow-lg transition-transform duration-300 ease-in-out ${props.isMobileMenuOpen() ? 'translate-y-0' : '-translate-y-full'}`}>
                <div class="flex flex-col space-y-4">
                    <For each={props.navLinks}>
                        {(link) => (
                            <A href={link.href} class="text-lg font-medium text-slate-800 dark:text-slate-200 hover:text-emerald-500" onClick={props.closeMenu}>
                                {link.text}
                            </A>
                        )}
                    </For>

                    <div class="pt-4 border-t border-slate-200 dark:border-slate-700">
                        <a
                            href="/get-started"
                            onClick={props.closeMenu}
                            class="w-full inline-flex items-center justify-center px-4 py-2 text-sm font-medium text-white bg-gradient-to-r from-emerald-600 to-emerald-500 rounded-lg"
                        >
                            Get Started
                        </a>
                    </div>
                </div>
            </div >
        </div >
    )
}