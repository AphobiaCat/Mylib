use public::{DBG_LOG, parse_json, build_json};

pub fn test_public(){
    DBG_LOG!("hello world :", 123);
}


#[derive(Debug, Serialize, Deserialize)]
struct Person {
    name: String,
    age: u32,
}

fn test_json() {
    let json = r#"{ "name": "Dunty", "age": 27 }"#;

    let person: Person = parse_json(json).unwrap();
    DBG_LOG!("Parsed: ", person);

    let json_str = build_json(&person).unwrap();
    DBG_LOG!("Serialized: ", json_str);
}


define_global!(TEST_VAL, u32, 0_u32);

fn example() {
    let mut val = TEST_VAL.lock().unwrap();
    *val += 1;
    println!("val: {}", *val);
}

#[derive(Debug)]
struct TestS{
    a   : i32,
    b   : f32,
}

impl TestS{
    fn new(a_val: i32, b_val: f32)->Self{
        Self{
            a: a_val,
            b: b_val,
        }
    }
}

define_global_async!(TEST_VAL_2, TestS, TestS::new(1_i32, 2_f32));


fn example2() {
    let mut val = TEST_VAL_2.lock().unwrap();
    val.a += 1_i32;
    val.b += 1_f32;
    
    println!("val: {:?}", *val);
}

fn test_global(){
    example();
    example();
    example();

    let rt = tokio::runtime::Runtime::new().unwrap();
    rt.block_on(example2());
    rt.block_on(example2());
    rt.block_on(example2());
}