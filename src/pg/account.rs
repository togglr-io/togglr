use crate::pg::error::Result;
use togglr::account::{Account, CreateAccount, UpdateAccount};

use sqlx::postgres::PgPool;
use uuid::Uuid;

#[derive(Clone)]
pub struct PgAccountService {
    pool: PgPool,
}

impl PgAccountService {
    pub fn new(pool: PgPool) -> PgAccountService {
        PgAccountService { pool }
    }

    pub async fn get(&self, id: Uuid) -> Result<Account> {
        sqlx::query_as!(
            Account,
            "
            SELECT
                id,
                name,
                created_at,
                updated_at
            FROM
                accounts
            WHERE
                id = $1
            ",
            id
        )
        .fetch_one(&self.pool)
        .await
    }

    pub async fn list(&self) -> Result<Vec<Account>> {
        sqlx::query_as!(
            Account,
            "
            SELECT
                id,
                name,
                created_at,
                updated_at
            FROM
                accounts
            "
        )
        .fetch_all(&self.pool)
        .await
    }

    pub async fn create(&self, req: CreateAccount) -> Result<Uuid> {
        let id = match req.id {
            Some(i) => i,
            None => Uuid::new_v4(),
        };

        let _ = sqlx::query("INSERT INTO accounts (id, name) values ($1, $2)")
            .bind(id)
            .bind(req.name)
            .execute(&self.pool)
            .await?;

        Ok(id)
    }

    pub async fn update(&self, req: UpdateAccount) -> Result<()> {
        sqlx::query(
            "
            UPDATE
                accounts
            SET
                name = $1
            WHERE
                id = $2
                ",
        )
        .bind(req.name)
        .bind(req.id)
        .execute(&self.pool)
        .await?;
        Ok(())
    }

    pub async fn delete(&self, id: Uuid) -> Result<()> {
        sqlx::query("DELETE FROM accounts WHERE id = $1")
            .bind(id)
            .execute(&self.pool)
            .await?;

        Ok(())
    }
}
