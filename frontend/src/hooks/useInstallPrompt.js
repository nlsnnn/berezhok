import { useEffect, useState } from "react";

export function useInstallPrompt() {
  const [prompt, setPrompt] = useState(null);
  const [isInstalled, setIsInstalled] = useState(
    () => window.matchMedia("(display-mode: standalone)").matches
  );

  useEffect(() => {
    const handleBeforeInstall = (e) => {
      e.preventDefault();
      setPrompt(e);
    };

    const handleInstalled = () => {
      setIsInstalled(true);
      setPrompt(null);
    };

    window.addEventListener("beforeinstallprompt", handleBeforeInstall);
    window.addEventListener("appinstalled", handleInstalled);

    return () => {
      window.removeEventListener("beforeinstallprompt", handleBeforeInstall);
      window.removeEventListener("appinstalled", handleInstalled);
    };
  }, []);

  const install = async () => {
    if (!prompt) return;
    prompt.prompt();
    const { outcome } = await prompt.userChoice;
    if (outcome === "accepted") setPrompt(null);
  };

  const dismiss = () => setPrompt(null);

  return { canInstall: !!prompt && !isInstalled, install, dismiss };
}
