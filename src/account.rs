use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Account {
    pub id: Uuid,
    pub name: String,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

impl Account {
    pub fn new(name: &str) -> Account {
        Account {
            id: Uuid::new_v4(),
            name: String::from(name),
            created_at: Utc::now(),
            updated_at: Utc::now(),
        }
    }
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct CreateAccount {
    pub id: Option<Uuid>,
    pub name: String,
}

impl CreateAccount {
    pub fn new(name: &str) -> CreateAccount {
        CreateAccount {
            id: Some(Uuid::new_v4()),
            name: String::from(name),
        }
    }
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct UpdateAccount {
    pub id: Uuid,
    pub name: String,
}

impl UpdateAccount {
    pub fn new(id: Uuid, name: &str) -> UpdateAccount {
        UpdateAccount {
            id,
            name: String::from(name),
        }
    }
}

// pub trait AccountService {
//     fn get(&self, id: Uuid) -> Result<Account, String>;
//     fn list(&self) -> Result<Vec<Account>, String>;
//     fn create(&self, account: Account) -> Result<Uuid, String>;
//     fn update(&self, req: UpdateAccountReq) -> Result<(), String>;
//     fn delete(&self, id: Uuid) -> Result<(), String>;
// }
