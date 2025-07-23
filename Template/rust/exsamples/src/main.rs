// mod public_test;
// mod dynamic_code_test;
// mod thread_manager_test;
// mod test_route_ws_client;
mod test_http;

use public::sleep_ms;

fn main(){

	// public_test::test_public();
	// let _ = dynamic_code_test::test_dynamic_code();
	// thread_manager_test::test_thread();

	let rt = tokio::runtime::Runtime::new().unwrap();
    // rt.block_on(async {
	// 	test_route_ws_client::test_route_ws_client().await;
	// });
	
	rt.block_on(async {
		let _ = test_http::test_http_manager().await;
	});

	sleep_ms(10000);
}