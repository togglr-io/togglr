use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::postgres::PgPool;
use uuid::Uuid;

use crate::pg::error::Result;
use togglr::user::{CreateUser, Identity, UpdateUser, User};

#[derive(Clone)]
pub struct PgUserService {
    pool: PgPool,
}

#[derive(Debug, Serialize, Deserialize)]
struct PgUser {
    id: Uuid,
    name: String,
    email: String,
    identity_type: String,
    created_at: DateTime<Utc>,
    updated_at: DateTime<Utc>,
}

impl From<User> for PgUser {
    fn from(user: User) -> Self {
        PgUser {
            id: user.id,
            name: user.name,
            email: user.email,
            identity_type: String::from(user.identity),
            created_at: user.created_at,
            updated_at: user.updated_at,
        }
    }
}

impl Into<User> for PgUser {
    fn into(self) -> User {
        User {
            id: self.id,
            name: self.name,
            email: self.email,
            identity: Identity::from(self.identity_type),
            created_at: self.created_at,
            updated_at: self.updated_at,
        }
    }
}

impl PgUserService {
    pub fn new(pool: PgPool) -> PgUserService {
        PgUserService { pool }
    }

    pub async fn get(&self, id: Uuid) -> Result<User> {
        let pg_user = sqlx::query_as!(
            PgUser,
            "
            SELECT
                id,
                name,
                email,
                identity_type,
                created_at,
                updated_at
            FROM
                users
            WHERE
                id = $1
            ",
            id
        )
        .fetch_one(&self.pool)
        .await?;

        Ok(pg_user.into())
    }

    pub async fn list(&self) -> Result<Vec<User>> {
        let pg_users = sqlx::query_as!(
            PgUser,
            "
            SELECT
                id,
                name,
                email,
                identity_type,
                created_at,
                updated_at
            FROM
                users
            "
        )
        .fetch_all(&self.pool)
        .await?;

        // map from a Vec<PgUser> to a Vec<User>
        Ok(pg_users.into_iter().map(|user| user.into()).collect())
    }

    pub async fn create(&self, req: CreateUser) -> Result<Uuid> {
        let id = match req.id {
            Some(i) => i,
            None => Uuid::new_v4(),
        };

        let identity_type = String::from(req.identity);
        let _ = sqlx::query(
            "
            INSERT INTO users (id, name, email, identity_type)
            VALUES ($1, $2, $3, $4)
            ",
        )
        .bind(id)
        .bind(req.name)
        .bind(req.email)
        .bind(identity_type)
        .execute(&self.pool)
        .await?;
        Ok(id)
    }

    pub async fn update(&self, req: UpdateUser) -> Result<()> {
        let _ = sqlx::query(
            "
            UPDATE
                users
            SET
                name = $1,
                email = $2
            WHERE
                id = $3
            ",
        )
        .bind(req.name)
        .bind(req.email)
        .bind(req.id)
        .execute(&self.pool)
        .await?;

        Ok(())
    }
}
