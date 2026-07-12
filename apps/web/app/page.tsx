import { SyncButton } from "./sync-button";

type MeResponse = {
  steam_id64: string;
  account_id: number;
  display_name: string;
  last_synced_at: string | null;
};

type RecentMatch = {
  match_id: number;
  start_time: string;
  duration_seconds: number;
  hero_id: number;
  hero_level?: number;
  kills?: number;
  deaths?: number;
  assists?: number;
  net_worth?: number;
  won?: boolean;
};

async function fetchJson<T>(url: string): Promise<T | null> {
  try {
    const response = await fetch(url, { cache: "no-store" });
    if (!response.ok) {
      return null;
    }
    return (await response.json()) as T;
  } catch {
    return null;
  }
}

export default async function HomePage() {
  const apiBase = process.env.API_BASE_URL ?? "http://localhost:8080";
  const [me, matchesPayload] = await Promise.all([
    fetchJson<MeResponse>(`${apiBase}/api/v1/me`),
    fetchJson<{ matches: RecentMatch[] }>(`${apiBase}/api/v1/me/matches`),
  ]);

  const matches = matchesPayload?.matches ?? [];
  const heroCount = new Map<number, number>();
  for (const match of matches) {
    heroCount.set(match.hero_id, (heroCount.get(match.hero_id) ?? 0) + 1);
  }
  const topHero = [...heroCount.entries()].sort((a, b) => b[1] - a[1])[0];
  const topHeroLabel = topHero ? `Hero ${topHero[0]}` : "No matches yet";
  const topHeroCount = topHero ? `${topHero[1]} matches` : "Run sync first";

  const cards = [
    {
      title: "Player",
      value: me?.display_name ?? "Not synced",
      sub: me
        ? `${me.steam_id64} • account ${me.account_id}`
        : "Start the API sync to populate this",
    },
    {
      title: "Recent Matches",
      value: String(matches.length),
      sub: me?.last_synced_at
        ? `Last synced ${new Date(me.last_synced_at).toLocaleString()}`
        : "No sync timestamp yet",
    },
    {
      title: "Most Played Hero",
      value: topHeroLabel,
      sub: topHeroCount,
    },
  ];

  return (
    <main className="page">
      <section className="hero">
        <p className="eyebrow">Deadlock Analytics</p>
        <h1>Cursed Apple Stats</h1>
        <p className="lead">
          Track match history, hero habits, item trends, and weird friend-group
          stats powered by Go APIs, Rust analytics, and Postgres snapshots.
        </p>
        <div className="actions">
          <SyncButton apiBase={apiBase} />
          <button className="btn ghost">Generate Analysis</button>
        </div>
      </section>

      <section className="grid" aria-label="summary cards">
        {cards.map((card) => (
          <article key={card.title} className="card">
            <p className="card-title">{card.title}</p>
            <p className="card-value">{card.value}</p>
            <p className="card-sub">{card.sub}</p>
          </article>
        ))}
      </section>

      <section className="matches-section">
        <div className="matches-header">
          <h2>Recent Matches</h2>
          <p>
            {matches.length
              ? "Latest synced match history from Deadlock API"
              : "Run sync to load match history"}
          </p>
        </div>
        <div className="matches-list">
          {matches.map((match) => (
            <article key={match.match_id} className="match-row">
              <div>
                <p className="match-hero">Hero {match.hero_id}</p>
                <p className="match-meta">
                  {new Date(match.start_time).toLocaleString()} •{" "}
                  {Math.round(match.duration_seconds / 60)} min
                </p>
              </div>
              <div className="match-stats">
                <span>
                  {match.kills ?? 0}/{match.deaths ?? 0}/{match.assists ?? 0}
                </span>
                <span>{match.net_worth ?? 0} nw</span>
                <span>
                  {match.won === undefined
                    ? "pending"
                    : match.won
                      ? "win"
                      : "loss"}
                </span>
              </div>
            </article>
          ))}
        </div>
      </section>
    </main>
  );
}
