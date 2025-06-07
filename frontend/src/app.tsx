import { Router } from "@solidjs/router";
import { FileRoutes } from "@solidjs/start/router";
import { Suspense } from "solid-js";
import Footer from "~/components/navigation/Footer";
import Nav from "~/components/navigation/Nav";
import { ThemeProvider } from "~/contexts/ThemeContext";
import { UserProvider } from "~/contexts/UserContext";
import "./app.css";

export default function App() {
  return (
    <ThemeProvider>
      <UserProvider>
        <Router
          root={props => (
            <main class="min-h-screen bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100">
              <Nav />
              <Suspense>{props.children}</Suspense>
              <Footer />
            </main>
          )}
        >
          <FileRoutes />
        </Router>
      </UserProvider>
    </ThemeProvider>
  );
}
