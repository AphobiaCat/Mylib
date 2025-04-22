use std::error::Error;

use public::DBG_LOG;
use dynamic_code::DynamicCode;

pub fn test_dynamic_code() -> Result<(), Box<dyn Error>> {

    let script = r#"
        pub fn add(a, b) {
            a + b
        }

        pub fn sum_matrix(matrix) {
            let sum = 0;
            for row in matrix {
                for col in row {
                    sum += col;
                }
            }
            sum
        }
    "#;

    let mut dynamic_code = DynamicCode::new(script)?;

    let result:i32 = dynamic_code.use_func("add", (100244, 2048))?;

    DBG_LOG!("result is ", result);

    let result:i64 = dynamic_code.use_func_dyn::<i64>("sum_matrix", "[[[1, 2, 3], [4, 5, 6]]]")?;
    
    DBG_LOG!("result is ", result);

    Ok(())   
}
