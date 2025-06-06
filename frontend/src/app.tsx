import { Router } from "@solidjs/router";
import { FileRoutes } from "@solidjs/start/router";
import { Suspense } from "solid-js";
import Nav from "~/components/Nav";
import { ThemeProvider } from "~/contexts/ThemeContext";
import "./app.css";

export default function App() {
  return (
    <ThemeProvider>
      <Router
        root={props => (
          <main class="min-h-screen bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100">
            <Nav />
            <Suspense>{props.children}</Suspense>
          </main>
        )}
      >
        <FileRoutes />
      </Router>
    </ThemeProvider>
  );
}
