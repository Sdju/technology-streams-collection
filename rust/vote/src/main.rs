mod vote;

use crate::vote::Vote;
use axum::extract::Query;
use axum::{extract::State, routing::get, Router};
use serde::Deserialize;
use std::sync::{Arc, Mutex};
use tower_http::cors::{Any, CorsLayer};

#[derive(Clone)]
struct AppState {
    vote: Arc<Mutex<Vote>>,
}

async fn show(State(state): State<AppState>) -> String {
    state.vote.lock().unwrap().get_votes_repr()
}

async fn votes_json(State(state): State<AppState>) -> String {
    state.vote.lock().unwrap().get_as_json()
}

#[derive(Deserialize)]
struct VoteUpQuery {
    item: String,
}

async fn vote_up(Query(query): Query<VoteUpQuery>, State(state): State<AppState>) -> String {
    println!("+ {}", query.item);
    let mut vote = state.vote.lock().unwrap();
    vote.vote_up(&query.item);
    vote.get_as_json()
}

#[derive(Deserialize)]
struct BanQuery {
    item: String,
}

async fn ban(Query(query): Query<BanQuery>, State(state): State<AppState>) -> String {
    println!("❌ {}", query.item);
    let mut vote = state.vote.lock().unwrap();
    vote.ban(&query.item);
    vote.get_as_json()
}

#[derive(Deserialize)]
struct RenameQuery {
    item: String,
    to: String,
}

async fn rename(Query(query): Query<RenameQuery>, State(state): State<AppState>) -> String {
    println!("ℹ️ {} -> {}", query.item, query.to);
    let mut vote = state.vote.lock().unwrap();
    vote.rename(&query.item, &query.to);
    vote.get_as_json()
}

#[derive(Deserialize)]
struct ConfirmQuery {
    item: String,
    target: Option<String>,
}

async fn confirm(Query(query): Query<ConfirmQuery>, State(state): State<AppState>) -> String {
    let mut vote = state.vote.lock().unwrap();
    if let Some(target) = query.target {
        println!("✔️ {} -> {}", query.item, target);
        vote.confirm_as_alias(&query.item, &target);
    } else {
        println!("✔️ {}", query.item);
        vote.confirm(&query.item);
    }
    vote.get_as_json()
}

#[tokio::main]
async fn main() {
    let shared_state = AppState {
        vote: Arc::new(Mutex::new(Vote::new())),
    };

    {
        let mut vote = shared_state.vote.lock().unwrap();
        vote.load();
    }

    let cors = CorsLayer::new().allow_origin(Any);

    let app = Router::new()
        .route("/", get(show))
        .route("/show", get(show))
        .route("/votes", get(votes_json))
        .route("/vote", get(vote_up))
        .route("/confirm", get(confirm))
        .route("/ban", get(ban))
        .route("/rename", get(rename))
        .layer(cors)
        .with_state(shared_state);

    println!("Server started");

    axum::Server::bind(&"127.0.0.1:8000".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
