"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

type SyncButtonProps = {
  apiBase: string;
};

export function SyncButton({ apiBase }: SyncButtonProps) {
  const router = useRouter();
  const [status, setStatus] = useState<string | null>(null);

  const handleClick = async () => {
    setStatus("Syncing...");

    try {
      const response = await fetch(`${apiBase}/api/v1/me/sync`, {
        method: "POST",
      });

      if (!response.ok) {
        const body = await response.json().catch(() => null);
        throw new Error(body?.error ?? `Sync failed with ${response.status}`);
      }

      setStatus("Sync complete");
      router.refresh();
    } catch (error) {
      setStatus(error instanceof Error ? error.message : "Sync failed");
    }
  };

  return (
    <div>
      <button className="btn primary" type="button" onClick={handleClick}>
        Sync My Matches
      </button>
      {status ? <p className="sync-status">{status}</p> : null}
    </div>
  );
}
