use leptos::*;
use leptos::leptos_dom::logging::console_log;
use serde::{Deserialize, Serialize};
use thiserror::Error;

#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize)]
pub struct Cat {
    url: String,
}

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

type CatCount = usize;

async fn fetch_cats(count: CatCount) -> error::Result<Vec<String>> {
    if count > 0 {
        // make the request
        let res = reqwasm::http::Request::get(&format!(
            "https://api.thecatapi.com/v1/images/search?limit={count}",
        ))
            .send()
            .await?
            // convert it to JSON
            .json::<Vec<Cat>>()
            .await?
            // extract the URL field for each cat
            .into_iter()
            .take(count)
            .map(|cat| cat.url)
            .collect::<Vec<_>>();
        Ok(res)
    } else {
        Err(CatError::NonZeroCats.into())
    }
}

async fn fetch_votes(_: CatCount) -> error::Result<VotesJsonRepr> {
    // make the request
    let res = reqwasm::http::Request::get(&format!(
        "http://127.0.0.1:8000/votes",
    ))
        .send()
        .await?
        .json::<VotesJsonRepr>()
        .await?;
    Ok(res)
}

#[component]
pub fn SimpleCounter(initial_value: i32) -> impl IntoView {
    let (cat_count, set_cat_count) = create_signal::<CatCount>(10);

    // we use local_resource here because
    // 1) our error type isn't serializable/deserializable
    // 2) we're not doing server-side rendering in this example anyway
    //    (during SSR, create_resource will begin loading on the server and resolve on the client)
    let votes = create_local_resource(cat_count, fetch_votes);

    let confirmed_votes_view = move || votes.and_then(|data: &VotesJsonRepr| {
        data.confirmed.iter()
            .map(|s| view! { <p>{&s.0} -  {s.1}</p> })
            .collect_view()
    });

    let unconfirmed_votes_view = move || votes.and_then(|data: &VotesJsonRepr| {
        data.unconfirmed.iter()
            .map(|s| view! { <p>{&s.0} -  {s.1}</p> })
            .collect_view()
    });

    let (value, set_value) = create_signal(initial_value);
    //
    // let clear = move |_| set_value(0);
    // let decrement = move |_| set_value.update(|value| *value -= 1);
    // let increment = move |_| set_value.update(|value| *value += 1);

    // create user interfaces with the declarative `view!` macro
    view! {
        <div>
            // <button on:click=clear>Clear</button>
            // <button on:click=decrement>-1</button>
            // text nodes can be quoted or unquoted
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