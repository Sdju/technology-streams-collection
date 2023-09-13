use std::cmp::Reverse;
use std::collections::{HashMap, HashSet};
use itertools::{Itertools, Either};
use std::fmt::Write;
use std::fs;
use serde::{Deserialize, Serialize};

#[derive(Debug)]
#[derive(Serialize)]
#[derive(Deserialize)]
pub struct VoteItem {
    pub name: String,
    pub aliases: HashSet<String>,
    pub count: i32,
    pub confirmed: bool,
}

impl VoteItem {
    pub fn new(str: &str) -> Self {
        Self {
            name: str.to_string(),
            aliases: HashSet::new(),
            count: 1,
            confirmed: false,
        }
    }
}


#[derive(Default)]
#[derive(Debug)]
#[derive(Serialize)]
#[derive(Deserialize)]
pub struct Vote {
    pub items: HashMap<String, VoteItem>,
    pub aliases: HashMap<String, String>,
    pub ban: HashSet<String>,
}

const DB_FILE: &str = "vote_db.json";

impl Vote {
    pub fn new() -> Self {
        Self {
            items: HashMap::new(),
            aliases: HashMap::new(),
            ban: HashSet::new(),
        }
    }

    pub fn vote_up(&mut self, vote_str: &str) {
        if self.ban.get(vote_str).is_some() {
            return;
        }

        if let Some(alias) = self.aliases.get(vote_str) {
            let item = self.items.get_mut(alias).expect("Alias for nonexistent item");
            item.count += 1;
            return;
        }

        self.save();

        self.items
            .entry(vote_str.to_string())
            .and_modify(|vote_item| vote_item.count += 1)
            .or_insert(VoteItem::new(vote_str));
    }

    pub fn rename(&mut self, from_str: &str, to_str: &str) {
        let item = self.items.remove(from_str).expect("Rename for nonexistent item");
        if let Some(old_item) = self.items.get_mut(to_str) {
            old_item.count += item.count;
        } else {
            self.items.insert(to_str.to_string(), item);
        }

        self.save();
    }

    pub fn confirm(&mut self, vote_str: &str) {
        self.items
            .entry(vote_str.to_string())
            .and_modify(|vote_item| {
                vote_item.confirmed = true;
            });
        self.save();
    }

    pub fn confirm_as_alias(&mut self, vote_str: &str, alias_to_str: &str) {
        let vote_item = self.items.remove(vote_str).expect("Confirm for nonexistent element");
        let count = vote_item.count;
        let mut alias_to_item = self.items.get_mut(alias_to_str).unwrap();

        alias_to_item.aliases.insert(vote_str.to_string());
        alias_to_item.count += count;
        self.aliases.insert(vote_str.to_string(), alias_to_str.to_string());
        self.save();
    }

    pub fn ban(&mut self, vote: &str) {
        self.ban.insert(vote.to_string());
        self.items.remove(vote);
    }

    pub fn get_votes_repr(&self) -> String {
        let mut items_pairs = self.items
            .iter()
            .collect::<Vec<_>>();
        items_pairs.sort_by_key(|(_, item)| Reverse(item.count));

        let (confirmed, unconfirmed): (String, String) = items_pairs
            .into_iter()
            .partition_map(|(name, item)| {
                let line = format!("{} - {}\n", name, item.count);
                match item.confirmed {
                    true => Either::Left(line),
                    false => Either::Right(line)
                }
            });

        let mut result = String::new();

        if !confirmed.is_empty() {
            writeln!(&mut result, "{confirmed}").expect("TODO: panic message");
        }

        if !unconfirmed.is_empty() {
            writeln!(&mut result, "Unconfirmed:\n{unconfirmed}").expect("TODO: panic message");
        }

        result
    }

    pub fn save(&self) {
        let data = serde_json::to_string(self).expect("WTF");
        fs::write(DB_FILE, data).expect("Can't write to file");
    }

    pub fn load(&mut self) {
        if let Ok(data) = fs::read_to_string(DB_FILE) {
            *self = serde_json::from_str(&data).expect("File data has been corrupted")
        }
    }
}
