use anyhow::Context;
use sqlx::{PgPool, Row};
use std::env;
use std::time::Duration;
use tracing::{error, info};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt()
        .with_env_filter("info")
        .with_target(false)
        .init();

    let database_url = env::var("DATABASE_URL")
        .unwrap_or_else(|_| "postgres://postgres:postgres@localhost:5432/cursed_apple_stats?sslmode=disable".to_string());

    let interval_seconds = env::var("ANALYTICS_INTERVAL_SECONDS")
        .ok()
        .and_then(|v| v.parse::<u64>().ok())
        .unwrap_or(300);

    let pool = PgPool::connect(&database_url)
        .await
        .context("failed to connect to postgres")?;

    info!("worker started; interval={}s", interval_seconds);

    let mut ticker = tokio::time::interval(Duration::from_secs(interval_seconds));

    loop {
        ticker.tick().await;

        if let Err(err) = recompute_basic_insights(&pool).await {
            error!(error = %err, "analytics iteration failed");
        }
    }
}

async fn recompute_basic_insights(pool: &PgPool) -> anyhow::Result<()> {
    let rows = sqlx::query(
        r#"
        SELECT account_id, COUNT(*)::BIGINT AS games
        FROM player_matches
        GROUP BY account_id
        "#,
    )
    .fetch_all(pool)
    .await
    .context("failed to load player match counts")?;

    for row in rows {
        let account_id: i64 = row.get("account_id");
        let games: i64 = row.get("games");

        let text = if games >= 50 {
            format!("Grinder mode: {} tracked matches and counting.", games)
        } else {
            format!("Warmup arc: {} tracked matches so far.", games)
        };

        sqlx::query(
            r#"
            INSERT INTO generated_insights (account_id, insight_type, insight_text, score)
            VALUES ($1, 'volume', $2, $3)
            "#,
        )
        .bind(account_id)
        .bind(text)
        .bind(games as f64)
        .execute(pool)
        .await
        .context("failed to insert generated insight")?;
    }

    info!("analytics iteration complete");
    Ok(())
}
