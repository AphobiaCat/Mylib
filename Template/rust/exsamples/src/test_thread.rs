use thread_manager::recv_msg;
use thread_manager::send_msg;
use thread_manager::try_recv_msg;
use thread_manager::ThreadManager;

use public::{sleep_ms, DBG_LOG};

fn task(id: usize) {
    if let Some(msg) = try_recv_msg::<i32>(id) {
        DBG_LOG!("thread:", id, " get msg (early): ", msg);
    } else {
        DBG_LOG!("thread:", id, " can't get msg (early)");
    }

    let msg: i32 = recv_msg::<i32>(id).expect("try get msg error");

    DBG_LOG!("thread:", id, " received msg: ", msg);
}

pub fn test_thread() {


    let mut manager = ThreadManager::new();

    manager.spawn_task(|| task(0));
    manager.spawn_task(|| task(1));


    sleep_ms(1000);

    DBG_LOG!("send msg");
    send_msg::<i32>(0, 1);
    send_msg::<i32>(1, 2);

    DBG_LOG!("wait recv msg");

    sleep_ms(1000);
    manager.join();
    DBG_LOG!("finish test");
}
