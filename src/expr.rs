#[derive(Debug, PartialEq, Eq)]
enum Kind {
    Ident,
    Keyword,
    StringLit,
    IntLit,
    FloatLit,
    Op,
    GroupStart,
    GroupEnd,
    Whitespace,
    Unknown,
}

#[derive(Debug)]
pub struct Token {
    kind: Kind,
    val: String,
}

impl Token {
    fn new(kind: Kind, val: String) -> Token {
        Token { kind, val }
    }
}

fn is_allowed_ident(ch: char, is_start: bool) -> bool {
    if is_start {
        ch.is_alphabetic() || ch == '_'
    } else {
        ch.is_alphanumeric() || ch == '_'
    }
}

fn is_terminal(ch: char) -> bool {
    ch.is_whitespace() || is_op(ch) || ch == '(' || ch == ')'
}

fn is_op(ch: char) -> bool {
    match ch {
        '=' => true,
        '!' => true,
        '&' => true,
        '|' => true,
        '>' => true,
        '<' => true,
        _ => false,
    }
}

fn get_kind(ch: char) -> Kind {
    if is_allowed_ident(ch, true) {
        return Kind::Ident;
    }

    if ch == '"' {
        return Kind::StringLit;
    }

    if ch.is_numeric() {
        return Kind::IntLit;
    }

    if is_op(ch) {
        return Kind::Op;
    }

    if ch == '(' {
        return Kind::GroupStart;
    }

    if ch == ')' {
        return Kind::GroupEnd;
    }

    if ch.is_whitespace() {
        return Kind::Whitespace;
    }

    Kind::Unknown
}

fn is_keyword(keyword: &str) -> bool {
    match keyword {
        "true" => true,
        "false" => true,
        _ => false,
    }
}

fn is_number(ch: char) -> bool {
    ch.is_numeric() || ch == '_'
}

pub fn lex(input: &str) -> Vec<Token> {
    let mut tokens: Vec<Token> = Vec::new();
    let mut kind = Kind::Unknown;
    let mut buf = String::with_capacity(512);
    let mut last: char = ' ';

    for ch in input.chars() {
        match kind {
            Kind::Ident => {
                if is_terminal(ch) {
                    if is_keyword(&buf) {
                        kind = Kind::Keyword;
                    }
                    tokens.push(Token::new(kind, String::from(&buf)));
                    buf.clear();
                    buf.push(ch);
                    kind = get_kind(ch);
                } else {
                    buf.push(ch);
                }
            }
            Kind::Op => {
                if buf.len() == 2 || !is_op(ch) {
                    tokens.push(Token::new(kind, String::from(&buf)));
                    buf.clear();
                    kind = get_kind(ch);
                }
                buf.push(ch);
            }
            Kind::StringLit => {
                if ch == '"' && last != '\\' {
                    last = ch; // last only matters in the context of a string literal
                    tokens.push(Token::new(kind, String::from(&buf)));
                    buf.clear();
                    kind = Kind::Unknown;
                } else {
                    buf.push(ch);
                }
            }
            Kind::IntLit => {
                if is_number(ch) {
                    buf.push(ch);
                } else if ch == '.' {
                    kind = Kind::FloatLit;
                    buf.push(ch);
                } else {
                    tokens.push(Token::new(kind, String::from(&buf)));
                    buf.clear();
                    buf.push(ch);
                    kind = get_kind(ch);
                }
            }
            Kind::FloatLit => {
                if is_number(ch) {
                    buf.push(ch);
                } else {
                    tokens.push(Token::new(kind, String::from(&buf)));
                    buf.clear();
                    buf.push(ch);
                    kind = get_kind(ch);
                }
            }
            Kind::GroupStart => {
                tokens.push(Token::new(kind, String::from("(")));
                buf.clear();
                buf.push(ch);
                kind = get_kind(ch);
            }
            Kind::GroupEnd => {
                tokens.push(Token::new(kind, String::from(")")));
                buf.clear();
                buf.push(ch);
                kind = get_kind(ch);
            }
            Kind::Whitespace => {
                buf.clear();
                buf.push(ch);
                kind = get_kind(ch);
            }
            Kind::Keyword => {
                unreachable!();
            }
            Kind::Unknown => {
                kind = get_kind(ch);
                if !matches!(kind, Kind::StringLit) {
                    buf.push(ch);
                }
            }
        };
    }

    // make sure to tokenize any leftover input
    if buf.len() > 0 {
        tokens.push(Token::new(kind, String::from(&buf)));
    }

    tokens
}

#[cfg(test)]
mod tests {
    use super::{lex, Kind};

    #[test]
    fn lex_simple_expr() {
        let input = "userType == \"admin\"";
        let tokens = lex(&input);

        assert!(tokens.len() == 3);
        assert!(matches!(tokens[0].kind, Kind::Ident));
        assert!(matches!(tokens[1].kind, Kind::Op));
        assert!(matches!(tokens[2].kind, Kind::StringLit));
    }

    #[test]
    fn lex_numbers() {
        let input = "age > 18 && bill < 1_000.25";
        let tokens = lex(&input);
        let expected_kinds = vec![
            Kind::Ident,
            Kind::Op,
            Kind::IntLit,
            Kind::Op,
            Kind::Ident,
            Kind::Op,
            Kind::FloatLit,
        ];

        assert!(tokens.len() == expected_kinds.len());

        for (idx, expected) in expected_kinds.into_iter().enumerate() {
            assert_eq!(tokens[idx].kind, expected);
        }
    }

    #[test]
    fn lex_group_expr() {
        let input = "userType == \"admin\" && (flag == true || otherFlag == false)";
        let tokens = lex(&input);
        println!("Tokens: {:?}", tokens);
        let expected_kinds = vec![
            Kind::Ident,      // userType
            Kind::Op,         // ==
            Kind::StringLit,  // admin
            Kind::Op,         // &&
            Kind::GroupStart, // (
            Kind::Ident,      // flag
            Kind::Op,         // ==
            Kind::Keyword,    // true
            Kind::Op,         // ||
            Kind::Ident,      // otherFlag
            Kind::Op,         // ==
            Kind::Keyword,    // false
            Kind::GroupEnd,   // )
        ];

        for (idx, expected) in expected_kinds.into_iter().enumerate() {
            assert_eq!(tokens[idx].kind, expected);
        }
    }
}
