use std::error::Error;

use public::DBG_LOG;
use dynamic_code::{DynamicCode, register_rust_function_i64, register_rust_function_matrix};

pub fn negative(a: i64) -> i64{
    -a
}

pub fn add_matrices(a: Vec<Vec<f32>>, b: Vec<Vec<f32>>) -> Vec<Vec<f32>> {
    a.iter()
        .zip(b.iter())
        .map(|(row_a, row_b)| {
            row_a.iter()
                .zip(row_b.iter())
                .map(|(x, y)| x + y)
                .collect()
        })
        .collect()
}


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

        pub fn negative(a){
            rust_negative(a)
        }

        pub fn add_matrix(a, b) {
            let result = rust_add_matrices(a, b);
            result
        }
    "#;

    register_rust_function_i64("rust_negative", negative);
    register_rust_function_matrix("rust_add_matrices", add_matrices);

    let mut dynamic_code = DynamicCode::new(script)?;

    let result:i32 = dynamic_code.use_func("add", (100244, 2048))?;

    DBG_LOG!("result is ", result);

    let result:i64 = dynamic_code.use_func_dyn::<i64>("sum_matrix", "[[[1, 2, 3], [4, 5, 6]]]")?;
    
    DBG_LOG!("result is ", result);

    let result:i64 = dynamic_code.use_func_dyn::<i64>("negative", "[1]")?;
    
    DBG_LOG!("result is ", result);

    let result:Vec<Vec<f32>> = dynamic_code.use_func_dyn::<Vec<Vec<f32>>>("add_matrix", "[[[1.0, 2.0, 3.0]], [[2.0, 3.0, 4.0]]]")?;
    
    DBG_LOG!("result is ", result);

    Ok(())   
}
