use std::error::Error;

use public::DBG_LOG;

fn main() -> Result<(), Box<dyn Error>> {

    DBG_LOG!("start");

    Ok(())  
}
