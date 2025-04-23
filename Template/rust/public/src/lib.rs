use std::{thread, time};
use serde::{Serialize, Deserialize};
use serde_json::{self, Error as SerdeError};

#[macro_export]
macro_rules! DBG_LOG {
    ($($arg:expr),*) => {
        {
            let file = std::file!();
            let line = std::line!();

            let args = vec![$(format!("{:?}", $arg)),*];

            let args_str = args.join("");

            println!("{:<20}|{:^5}| logs: {}", file, line, args_str);
        }
    };
}

#[macro_export]
macro_rules! DBG_ERR {
    ($($arg:expr),*) => {
        {
            let file = std::file!();
            let line = std::line!();

            let args = vec![$(format!("{:?}", $arg)),*];

            let args_str = args.join("");

            println!("\x1b[31m{:<20}|{:^5}| logs: {}\x1b[0m", file, line, args_str);
        }
    };
}

#[macro_export]
macro_rules! define_global {
    ($name:ident, $ty:ty, $expr:expr) => {
        static $name: once_cell::sync::Lazy<std::sync::Mutex<$ty>> = once_cell::sync::Lazy::new(|| std::sync::Mutex::new($expr));
    };
}

#[macro_export]
macro_rules! define_global_async {
    ($name:ident, $ty:ty, $expr:expr) => {
        static $name: once_cell::sync::Lazy<tokio::sync::Mutex<$ty>> = once_cell::sync::Lazy::new(|| tokio::sync::Mutex::new($expr));
    };
}

pub fn sleep_ms(duration_ms : u64){
    let duration = time::Duration::from_millis(duration_ms);
    thread::sleep(duration);
}


pub fn parse_json<T>(json_str: &str) -> Result<T, SerdeError>
where
    T: for<'de> Deserialize<'de>,
{
    if let Ok(value) = serde_json::from_str::<T>(json_str) {
        return Ok(value);
    }

    if let Ok(unescaped) = serde_json::from_str::<String>(json_str) {
        serde_json::from_str::<T>(&unescaped)
    } else {
        serde_json::from_str::<T>(json_str)
    }
}

pub fn build_json<T>(value: &T) -> Result<String, SerdeError>
where
    T: Serialize,
{
    serde_json::to_string(value)
}

