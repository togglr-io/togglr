use std::convert::Infallible;
use uuid::Uuid;

use crate::pg::account::PgAccountService;
use crate::pg::user::PgUserService;
use togglr::account::{CreateAccount, UpdateAccount};
use togglr::user::{CreateUser, UpdateUser};

use warp::Filter;

pub struct Services {
    pub account_service: PgAccountService,
    pub user_service: PgUserService,
}

fn with_service<T: Clone + Send>(
    service: T,
) -> impl Filter<Extract = (T,), Error = std::convert::Infallible> + Clone {
    warp::any().map(move || service.clone())
}

pub fn router(services: Services) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    account_handler(services.account_service)
        .or(user_handler(services.user_service))
}

pub fn account_handler(
    account_service: PgAccountService,
) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    let account = warp::path("account");
    let get_accounts = account
        .and(warp::get())
        .and(with_service(account_service.clone()))
        .and_then(list_accounts);

    let get_account = account
        .and(warp::get())
        .and(warp::path::param())
        .and(with_service(account_service.clone()))
        .and_then(get_account);

    let post_account = account
        .and(warp::post())
        .and(warp::body::content_length_limit(1024 * 1024))
        .and(warp::body::json())
        .and(with_service(account_service.clone()))
        .and_then(create_account);

    let put_account = account
        .and(warp::put())
        .and(warp::body::content_length_limit(1024 * 1024))
        .and(warp::body::json())
        .and(with_service(account_service.clone()))
        .and_then(update_account);

    get_account
        .or(get_accounts)
        .or(post_account)
        .or(put_account)
}

pub fn user_handler(
    user_service: PgUserService
) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    let user = warp::path("user");
    let get_users = user
        .and(warp::get())
        .and(with_service(user_service.clone()))
        .and_then(list_users);

    let get_user = user
        .and(warp::get())
        .and(warp::path::param())
        .and(with_service(user_service.clone()))
        .and_then(get_user);

    let post_user = user
        .and(warp::post())
        .and(warp::body::content_length_limit(1024 * 1024))
        .and(warp::body::json())
        .and(with_service(user_service.clone()))
        .and_then(create_user);

    let put_user = user
        .and(warp::put())
        .and(warp::body::content_length_limit(1024 * 1024))
        .and(warp::body::json())
        .and(with_service(user_service.clone()))
        .and_then(update_user);

    get_user
        .or(get_users)
        .or(post_user)
        .or(put_user)
}

async fn list_accounts(account_service: PgAccountService) -> Result<impl warp::Reply, Infallible> {
    let accounts = account_service.list().await.unwrap();

    Ok(warp::reply::json(&accounts))
}

async fn get_account(
    id: String,
    account_service: PgAccountService,
) -> Result<impl warp::Reply, Infallible> {
    let uid = match Uuid::parse_str(&id) {
        Ok(i) => i,
        Err(_) => return Ok(warp::reply::json(&String::from(""))),
    };

    let account = account_service.get(uid).await.unwrap();
    Ok(warp::reply::json(&account))
}

async fn create_account(
    req: CreateAccount,
    account_service: PgAccountService,
) -> Result<impl warp::Reply, Infallible> {
    let id = account_service.create(req).await.unwrap();
    Ok(warp::reply::json(&id))
}

async fn update_account(req: UpdateAccount, account_service: PgAccountService) -> Result<impl warp::Reply, Infallible> {
    account_service.update(req).await.unwrap();
    Ok(warp::reply())
}

async fn list_users(user_service: PgUserService) -> Result<impl warp::Reply, Infallible> {
    let users = user_service.list().await.unwrap();

    Ok(warp::reply::json(&users))
}

async fn get_user(
    id: String,
    user_service: PgUserService,
) -> Result<impl warp::Reply, Infallible> {
    let uid = match Uuid::parse_str(&id) {
        Ok(i) => i,
        Err(_) => return Ok(warp::reply::json(&String::from(""))),
    };

    let user = user_service.get(uid).await.unwrap();
    Ok(warp::reply::json(&user))
}

async fn create_user(
    req: CreateUser,
    user_service: PgUserService,
) -> Result<impl warp::Reply, Infallible> {
    let id = user_service.create(req).await.unwrap();
    Ok(warp::reply::json(&id))
}

async fn update_user(req: UpdateUser, user_service: PgUserService) -> Result<impl warp::Reply, Infallible> {
    user_service.update(req).await.unwrap();
    Ok(warp::reply())
}
