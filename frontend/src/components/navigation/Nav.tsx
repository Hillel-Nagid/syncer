import { A, useLocation } from '@solidjs/router';
import {
    createEffect,
    createSignal,
    For,
    Match,
    onCleanup,
    Switch,
} from 'solid-js';
import { useUser } from '~/contexts/UserContext';
import UserProfile from '../auth/UserProfile';
import Button from '../ui/Button';
import Icon from '../ui/Icon';
import { SkeletonProfile } from '../ui/Skeleton';
import MobileMenu from './MobileMenu';
import ThemeToggle from './ThemeToggle';

const navLinks = [
    { href: '/', text: 'Home' },
    { href: '/about', text: 'About' },
    { href: '/dashboard', text: 'Dashboard' },
    { href: '/faq', text: 'FAQ' },
    { href: '/contact', text: 'Contact' },
];

export default function Nav() {
    const location = useLocation();
    const { isLoggedIn, isLoading } = useUser();
    const [isMobileMenuOpen, setIsMobileMenuOpen] = createSignal(false);

    const active = (path: string) =>
        path == location.pathname
            ? 'text-emerald-600 dark:text-emerald-400 font-medium'
            : 'text-slate-600 dark:text-slate-400 hover:text-emerald-600 hover:dark:text-emerald-400 transition-colors';

    createEffect(() => {
        if (isMobileMenuOpen()) {
            document.body.classList.add('overflow-hidden', 'md:overflow-auto');
        } else {
            document.body.classList.remove('overflow-hidden', 'md:overflow-auto');
        }
        onCleanup(() =>
            document.body.classList.remove('overflow-hidden', 'md:overflow-auto')
        );
    });

    const closeMenu = () => setIsMobileMenuOpen(false);

    return (
        <>
            <nav class='sticky top-0 z-50 w-full border-b border-slate-200/60 dark:border-slate-700/60 bg-white/80 dark:bg-slate-900/80 backdrop-blur-xl supports-[backdrop-filter]:bg-white/70 supports-[backdrop-filter]:dark:bg-slate-900/70'>
                <div class='max-w-6xl mx-auto px-4 sm:px-6 lg:px-8'>
                    <div class='flex h-16 items-center justify-between'>
                        {/* Logo/Brand */}
                        <div class='flex-shrink-0'>
                            <a
                                href='/'
                                class='flex items-center space-x-3 text-slate-900 dark:text-slate-100 hover:text-emerald-600 hover:dark:text-emerald-400 transition-colors group'
                                onClick={closeMenu}
                            >
                                <div class='w-8 h-8 bg-gradient-to-br from-emerald-500 to-emerald-600 rounded-lg flex items-center justify-center shadow-lg group-hover:shadow-emerald-500/25 transition-shadow'>
                                    <Icon
                                        name='sync-icon'
                                        class='w-5 h-5 text-white'
                                        alt='Sync Icon'
                                    />
                                </div>
                                <span class='text-xl font-bold bg-gradient-to-r from-emerald-600 to-emerald-500 bg-clip-text text-transparent'>
                                    Syncer
                                </span>
                            </a>
                        </div>

                        {/* Desktop Navigation Links */}
                        <div class='hidden md:block'>
                            <div class='ml-10 flex items-baseline space-x-8'>
                                <For each={navLinks}>
                                    {(link) => (
                                        <A
                                            href={link.href}
                                            class={`px-3 py-2 text-sm font-medium transition-all duration-200 rounded-md ${active(
                                                link.href
                                            )}`}
                                        >
                                            {link.text}
                                        </A>
                                    )}
                                </For>
                            </div>
                        </div>

                        <div class='hidden md:flex items-center space-x-4'>
                            <Switch
                                fallback={
                                    <A href='/auth/signup'>
                                        <Button variant='primary' size='sm'>
                                            Sign Up
                                        </Button>
                                    </A>
                                }
                            >
                                <Match when={isLoggedIn()}>
                                    <UserProfile />
                                </Match>
                                <Match when={!isLoggedIn() && isLoading()}>
                                    <SkeletonProfile />
                                </Match>
                            </Switch>
                            <ThemeToggle />
                        </div>

                        <div class='md:hidden flex items-center space-x-2'>
                            <ThemeToggle />
                            <button
                                type='button'
                                class='inline-flex items-center justify-center p-2 rounded-md text-slate-600 dark:text-slate-400 hover:text-emerald-600 hover:dark:text-emerald-400 hover:bg-slate-100 dark:hover:bg-slate-800 focus:outline-none focus:ring-2 focus:ring-emerald-500 transition-colors'
                                attr:aria-expanded={isMobileMenuOpen()}
                                onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen())}
                            >
                                <span class='sr-only'>Open main menu</span>
                                <Icon name='menu-icon' class='w-6 h-6' alt='Menu Icon' />
                            </button>
                        </div>
                    </div>
                </div>
            </nav>
            <MobileMenu
                isMobileMenuOpen={isMobileMenuOpen}
                closeMenu={closeMenu}
                navLinks={navLinks}
            />
        </>
    );
}
