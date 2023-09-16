use leptos::*;
use leptos::leptos_dom::logging::console_log;
use serde::{Deserialize, Serialize};
use thiserror::Error;
use reqwasm::http::Request;

#[derive(Clone, Default, Debug, Serialize, Deserialize)]
struct VotesJsonRepr {
    confirmed: Vec<(String, i32)>,
    unconfirmed: Vec<(String, i32)>,
}

#[derive(Error, Clone, Debug)]
pub enum CatError {
    #[error("Please request more than zero cats.")]
    NonZeroCats,
}

async fn fetch_votes() -> error::Result<VotesJsonRepr> {
    // make the request
    let res = Request::get(&format!(
        "http://127.0.0.1:8000/votes",
    ))
        .send()
        .await?
        .json::<VotesJsonRepr>()
        .await?;
    Ok(res)
}

#[component]
pub fn ConfirmedItem(name: String, votes: i32, set_votes: WriteSignal<VotesJsonRepr>) -> impl IntoView {
    let name_copy = name.to_owned();
    let vote_up = create_action(move |_| {
        let name_copy = name_copy.to_owned();
        async move {
            let test = Request::get(&format!(
                "http://127.0.0.1:8000/vote?item={name_copy}",
            ))
                .send()
                .await
                .unwrap()
                .json::<VotesJsonRepr>()
                .await
                .unwrap();
            set_votes(test);
        }
    });

    view! {
        <div class="confirmed-item">
            <span>{&name} -  {votes}</span>
            <button
                class="confirmed-item__increment"
                on:click=move |_| vote_up.dispatch(())
            >+</button>
        </div>
    }
}

#[component]
pub fn SimpleCounter(initial_value: i32) -> impl IntoView {
    // let (confirmed, set_confirmed) = create_signal(Vec::<(String, i32)>::new());
    // let (unconfirmed, set_unconfirmed) = create_signal(Vec::<(String, i32)>::new());
    let (votes, set_votes) = create_signal(VotesJsonRepr::default());

    let (loading, set_loading) = create_signal(false);
    let votes_res = create_local_resource(move || loading.get(), move |_| {
        fetch_votes()
    });

    // Я СДАЛСЯ, Я НЕ ПОНЯЛ КАК НОРМАЛЬНО РЕШИТЬ ПРОБЛЕМЫ И ЗАПИХАЛ ВСЕ ХОТЬ КАК-НИБУДЬ
    let confirmed_votes_view = move || {
        if (votes.get().confirmed.is_empty()) {
            return votes_res.and_then(|data: &VotesJsonRepr| {
                data.confirmed.iter()
                    .map(|s| view! {
                <ConfirmedItem
                    name=s.0.clone()
                    votes=s.1
                    set_votes=set_votes
                />
            })
                    .collect_view()
            });
        }
        return Some(Ok(votes.get().confirmed.iter()
            .map(|s| view! {
                <ConfirmedItem
                    name=s.0.clone()
                    votes=s.1
                    set_votes=set_votes
                />
            })
            .collect_view()));
    };

    let unconfirmed_votes_view = move || votes_res.and_then(|data: &VotesJsonRepr| {
        data.unconfirmed.iter()
            .map(|s| view! { <p>{&s.0} -  {s.1}</p> })
            .collect_view()
    });

    let (value, set_value) = create_signal(initial_value);

    view! {
        <div>
            {confirmed_votes_view }
            <span> Unconfirmed </span>
            {unconfirmed_votes_view }
        </div>
    }
}

pub fn main() {
    mount_to_body(|| view! {
        <SimpleCounter initial_value=3 />
    })
}