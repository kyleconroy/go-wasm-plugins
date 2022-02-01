use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
struct HelloRequest {
    name: String,
}

fn main() {
    let req = HelloRequest {
        name: String::from("foo"),
    };

    // Serialize it to a JSON string.
    let j = serde_json::to_string(&req).unwrap();

    // Print, write to a file, or send to an HTTP server.
    println!("{}", j);
}
