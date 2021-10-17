use std::str::FromStr;

pub struct Env {
    prefix: Option<String>,
}

impl Env {
    pub fn new() -> Env {
        Env { prefix: None }
    }

    pub fn with_prefix(prefix: &str) -> Env {
        Env {
            prefix: Some(String::from(prefix)),
        }
    }

    pub fn get<T: FromStr>(&self, key: &str, default: Option<T>) -> Result<T, String> {
        let key = match &self.prefix {
            Some(prefix) => format!("{}_{}", prefix, key),
            None => key.to_string(),
        };

        get(&key, default)
    }
}

/// Fetches values from the environment with an optional default in the case that the key does not
/// exist. Any type that implements the FromStr can be returned.
/// ```
/// use crate::env;
///
/// let port = env::get::<u16>("SOME_ENV_PORT", Some(8080));
/// assert_eq!(port, 8080);
/// ```
pub fn get<T: FromStr>(key: &str, default: Option<T>) -> Result<T, String> {
    let parse_result = match std::env::var(key) {
        Ok(val) => T::from_str(&val),
        Err(_) => match default {
            Some(def) => return Ok(def),
            None => return Err(String::from("key not found")),
        },
    };

    match parse_result {
        Ok(val) => Ok(val),
        Err(_) => Err(String::from("parsing failure")),
    }
}

#[cfg(test)]
mod tests {
    #[test]
    fn get_existing() {
        std::env::set_var("TOGGLR_PORT", "42");
        let port = super::get::<u16>("TOGGLR_PORT", Some(24));
        match port {
            Ok(val) => assert_eq!(val, 42),
            _ => unreachable!(),
        }
    }

    #[test]
    fn get_default() {
        std::env::remove_var("TOGGLR_PORT");
        let port = super::get::<u16>("TOGGLR_PORT", Some(24));
        match port {
            Ok(val) => assert_eq!(val, 24),
            _ => unreachable!(),
        }
    }
}
