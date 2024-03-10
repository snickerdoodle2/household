use axum::routing::get;

async fn hello() -> String {
    "Hello, World!".to_string()
}

#[tokio::main]
async fn main() -> std::io::Result<()> {
    let app = axum::Router::new().route("/api/hello", get(hello));

    let listener = tokio::net::TcpListener::bind("127.0.0.1:8080").await?;

    axum::serve(listener, app).await?;

    Ok(())
}
