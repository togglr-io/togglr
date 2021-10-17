use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "lowercase")]
pub enum Identity {
    Github,
    Basic,
    Google,
    Unknown,
}

impl From<&str> for Identity {
    fn from(input: &str) -> Self {
        match input.to_lowercase().as_str() {
            "github" => Identity::Github,
            "basic" => Identity::Basic,
            "google" => Identity::Google,
            _ => Identity::Unknown,
        }
    }
}

impl From<String> for Identity {
    fn from(input: String) -> Self {
        Identity::from(input.as_str())
    }
}

impl From<Identity> for String {
    fn from(input: Identity) -> Self {
        match input {
            Identity::Github => String::from("github"),
            Identity::Basic => String::from("basic"),
            Identity::Google => String::from("google"),
            Identity::Unknown => String::from("unknown"),
        }
    }
}

#[derive(Debug, Serialize, Deserialize)]
pub struct User {
    pub id: Uuid,
    pub name: String,
    pub email: String,
    pub identity: Identity,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CreateUser {
    pub id: Option<Uuid>,
    pub name: String,
    pub email: String,
    pub identity: Identity,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct UpdateUser {
    pub id: Uuid,
    pub name: String,
    pub email: String,
}

#[cfg(test)]
mod tests {
    use super::Identity;

    #[test]
    fn identity_strings() {
        let github_str = String::from(Identity::Github);
        let google_str = String::from(Identity::Google);
        let basic_str = String::from(Identity::Basic);
        let unknown_str = "random";

        assert!(matches!(Identity::from(github_str), Identity::Github));
        assert!(matches!(Identity::from(google_str), Identity::Google));
        assert!(matches!(Identity::from(basic_str), Identity::Basic));
        assert!(matches!(Identity::from(unknown_str), Identity::Unknown));
    }
}
