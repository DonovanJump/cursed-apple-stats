const cards = [
  {
    title: "Most Played Hero",
    value: "Haze",
    sub: "42 matches",
  },
  {
    title: "Most Purchased Item",
    value: "Unstoppable",
    sub: "31 purchases",
  },
  {
    title: "Cursed Insight",
    value: "Greed Score 8.3",
    sub: "Wins spike after defensive 2nd item",
  },
];

export default function HomePage() {
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
          <button className="btn primary">Sign In With Steam</button>
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
    </main>
  );
}
