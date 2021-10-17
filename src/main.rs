use sqlx::postgres::PgPoolOptions;

mod env;
mod pg;
mod http;

use env::Env;
use http::Services;

use pg::account::PgAccountService;
use pg::user::PgUserService;

#[tokio::main]
async fn main() -> Result<(), sqlx::Error> {
    let env = Env::with_prefix("TOGGLE");

    // load database configuration from env
    let host: String = env.get("DB_HOST", Some(String::from("localhost"))).unwrap();
    let user: String = env.get("DB_USER", Some(String::from("toggle"))).unwrap();
    let password: String = env
        .get("DB_PASSWORD", Some(String::from("toggle")))
        .unwrap();
    let db_name: String = env.get("DB_NAME", Some(String::from("toggle"))).unwrap();
    let port: u16 = env.get("DB_PORT", Some(5432)).unwrap();
    let max_connections: u32 = env.get("DB_PORT", Some(10)).unwrap();

    let dsn = format!(
        "postgres://{}:{}@{}:{}/{}",
        user, password, host, port, db_name
    );

    // create connection pool
    let pool = PgPoolOptions::new()
        .max_connections(max_connections)
        .connect(&dsn)
        .await?;

    // create services utilizing connection pool
    let account_service = PgAccountService::new(pool.clone());
    let user_service = PgUserService::new(pool.clone());
    let services = Services {
        account_service: PgAccountService::new(pool.clone()),
        user_service: PgUserService::new(pool.clone()),
    };


    warp::serve(http::router(services))
        .run(([127, 0, 0, 1], 3030))
        .await;
    Ok(())
}



